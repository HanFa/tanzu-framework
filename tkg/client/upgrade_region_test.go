// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	capi "sigs.k8s.io/cluster-api/api/v1beta1"
	clusterctlv1 "sigs.k8s.io/cluster-api/cmd/clusterctl/api/v1alpha3"
	clusterctl "sigs.k8s.io/cluster-api/cmd/clusterctl/client"
	crtclient "sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/vmware-tanzu/tanzu-framework/tkg/client"
	"github.com/vmware-tanzu/tanzu-framework/tkg/clusterclient"
	"github.com/vmware-tanzu/tanzu-framework/tkg/constants"
	"github.com/vmware-tanzu/tanzu-framework/tkg/fakes"
)

var _ = Describe("Unit tests for upgrade management cluster", func() {
	var (
		err                   error
		regionalClusterClient *fakes.ClusterClient
		providerUpgradeClient *fakes.ProvidersUpgradeClient
		tkgClient             *TkgClient
		upgradeClusterOptions UpgradeClusterOptions
		context               string
	)

	BeforeEach(func() {
		regionalClusterClient = &fakes.ClusterClient{}
		providerUpgradeClient = &fakes.ProvidersUpgradeClient{}

		tkgClient, err = CreateTKGClient("../fakes/config/config.yaml", testingDir, "../fakes/config/bom/tkg-bom-v1.3.1.yaml", 2*time.Millisecond)
		Expect(err).NotTo(HaveOccurred())

		context = "fakeContext"
		upgradeClusterOptions = UpgradeClusterOptions{
			ClusterName:       "fake-cluster-name",
			Namespace:         "fake-namespace",
			KubernetesVersion: newK8sVersion,
			IsRegionalCluster: true,
			Kubeconfig:        "../fakes/config/kubeconfig/config1.yaml",
		}
	})

	Describe("When upgrading management cluster", func() {
		BeforeEach(func() {
			newK8sVersion = "v1.18.0+vmware.1"
			currentK8sVersion = "v1.17.3+vmware.2"
			setupBomFile("../fakes/config/bom/tkg-bom-v1.3.1.yaml", testingDir)
			regionalClusterClient.IsPacificRegionalClusterReturns(false, nil)
		})

		Context("When upgrading management cluster providers", func() {
			JustBeforeEach(func() {
				err = tkgClient.DoProvidersUpgrade(regionalClusterClient, context, providerUpgradeClient, &upgradeClusterOptions)
			})
			Context("When reading upgrade information from BOM file fails", func() {
				BeforeEach(func() {
					updateDefaultBoMFileName(testingDir, "tkg-bom-v1.3.1-fake.yaml")
				})
				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("unable to read in configuration from BOM file"))
				})
			})
			Context("When getting upgrade information fails due to failure to get the current providers information", func() {
				BeforeEach(func() {
					regionalClusterClient.ListResourcesReturns(errors.New("fake ListResourceError"))
				})
				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("fake ListResourceError"))
				})
			})

			Context("When current providers versions are up to date", func() {
				BeforeEach(func() {
					regionalClusterClient.ListResourcesCalls(func(providers interface{}, options ...crtclient.ListOption) error {
						installedProviders, _ := providers.(*clusterctlv1.ProviderList)
						installedProviders.Items = []clusterctlv1.Provider{
							{ObjectMeta: metav1.ObjectMeta{Namespace: "capi-system", Name: "cluster-api"}, Type: "CoreProvider", Version: "v0.3.11", ProviderName: "cluster-api"},
							{ObjectMeta: metav1.ObjectMeta{Namespace: "capi-kubeadm-bootstrap-system", Name: "bootstrap-kubeadm"}, Type: "BootstrapProvider", Version: "v0.3.11", ProviderName: "kubeadm"},
							{ObjectMeta: metav1.ObjectMeta{Namespace: "capi-kubeadm-control-plane-system", Name: "control-plane-kubeadm"}, Type: "ControlPlaneProvider", Version: "v0.3.11", ProviderName: "kubeadm"},
							{ObjectMeta: metav1.ObjectMeta{Namespace: "capv-system", Name: "infrastructure-vsphere"}, Type: "InfrastructureProvider", Version: "v0.7.1", ProviderName: "vsphere"},
						}
						return nil
					})
				})
				It("should not apply the providers version upgrade", func() {
					Expect(providerUpgradeClient.ApplyUpgradeCallCount()).Should(Equal(1))
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("When providers version upgrade failed", func() {
				BeforeEach(func() {
					regionalClusterClient.ListResourcesCalls(func(providers interface{}, options ...crtclient.ListOption) error {
						installedProviders, _ := providers.(*clusterctlv1.ProviderList)
						installedProviders.Items = []clusterctlv1.Provider{
							{ObjectMeta: metav1.ObjectMeta{Namespace: "capi-system", Name: "cluster-api"}, Type: "CoreProvider", Version: "v0.3.3", ProviderName: "cluster-api"},
							{ObjectMeta: metav1.ObjectMeta{Namespace: "capi-kubeadm-bootstrap-system", Name: "bootstrap-kubeadm"}, Type: "BootstrapProvider", Version: "v0.3.3", ProviderName: "kubeadm"},
							{ObjectMeta: metav1.ObjectMeta{Namespace: "capi-kubeadm-control-plane-system", Name: "control-plane-kubeadm"}, Type: "ControlPlaneProvider", Version: "v0.3.3", ProviderName: "kubeadm"},
							{ObjectMeta: metav1.ObjectMeta{Namespace: "capv-system", Name: "infrastructure-vsphere"}, Type: "InfrastructureProvider", Version: "v0.6.3", ProviderName: "vsphere"},
						}
						return nil
					})
					providerUpgradeClient.ApplyUpgradeReturns(errors.New("fake-providers upgrade failed"))
				})
				It("should return error", func() {
					Expect(providerUpgradeClient.ApplyUpgradeCallCount()).Should(Equal(1))
					Expect(err.Error()).To(ContainSubstring("fake-providers upgrade failed"))
				})
			})
			Context("When providers current versions and the latest versions are same", func() {
				var applyUpgradeOptionsRecvd clusterctl.ApplyUpgradeOptions
				BeforeEach(func() {
					regionalClusterClient.ListResourcesCalls(func(providers interface{}, options ...crtclient.ListOption) error {
						installedProviders, _ := providers.(*clusterctlv1.ProviderList)
						installedProviders.Items = []clusterctlv1.Provider{
							{ObjectMeta: metav1.ObjectMeta{Namespace: "capi-system", Name: "cluster-api"}, Type: "CoreProvider", Version: "v0.3.3", ProviderName: "cluster-api"},
							{ObjectMeta: metav1.ObjectMeta{Namespace: "capi-kubeadm-bootstrap-system", Name: "bootstrap-kubeadm"}, Type: "BootstrapProvider", Version: "v0.3.3", ProviderName: "kubeadm"},
							{ObjectMeta: metav1.ObjectMeta{Namespace: "capi-kubeadm-control-plane-system", Name: "control-plane-kubeadm"}, Type: "ControlPlaneProvider", Version: "v0.3.3", ProviderName: "kubeadm"},
							{ObjectMeta: metav1.ObjectMeta{Namespace: "capv-system", Name: "infrastructure-vsphere"}, Type: "InfrastructureProvider", Version: "v0.6.3", ProviderName: "vsphere"},
						}
						return nil
					})
					providerUpgradeClient.ApplyUpgradeCalls(func(applyUpgradeOptions *clusterctl.ApplyUpgradeOptions) error {
						applyUpgradeOptionsRecvd = *applyUpgradeOptions
						return nil
					})
				})
				It("should still apply the providers version upgrade to the same versions", func() {
					Expect(applyUpgradeOptionsRecvd.CoreProvider).Should(Equal("capi-system/cluster-api:v0.3.11"))
					Expect(applyUpgradeOptionsRecvd.BootstrapProviders[0]).Should(Equal("capi-kubeadm-bootstrap-system/kubeadm:v0.3.11"))
					Expect(applyUpgradeOptionsRecvd.ControlPlaneProviders[0]).Should(Equal("capi-kubeadm-control-plane-system/kubeadm:v0.3.11"))
					Expect(applyUpgradeOptionsRecvd.InfrastructureProviders[0]).Should(Equal("capv-system/vsphere:v0.7.1"))
					Expect(providerUpgradeClient.ApplyUpgradeCallCount()).Should(Equal(1))
					Expect(err).NotTo(HaveOccurred())
				})
			})
			Context("When some(cluster-api, Infrastructure-vsphere) providers current versions are not up to date", func() {
				var applyUpgradeOptionsRecvd clusterctl.ApplyUpgradeOptions
				BeforeEach(func() {
					regionalClusterClient.ListResourcesCalls(func(providers interface{}, options ...crtclient.ListOption) error {
						installedProviders, _ := providers.(*clusterctlv1.ProviderList)
						installedProviders.Items = []clusterctlv1.Provider{
							{ObjectMeta: metav1.ObjectMeta{Namespace: "capi-system", Name: "cluster-api"}, Type: "CoreProvider", Version: "v0.3.2", ProviderName: "cluster-api"},
							{ObjectMeta: metav1.ObjectMeta{Namespace: "capi-kubeadm-bootstrap-system", Name: "bootstrap-kubeadm"}, Type: "BootstrapProvider", Version: "v0.3.11", ProviderName: "kubeadm"},
							{ObjectMeta: metav1.ObjectMeta{Namespace: "capi-kubeadm-control-plane-system", Name: "control-plane-kubeadm"}, Type: "ControlPlaneProvider", Version: "v0.3.11", ProviderName: "kubeadm"},
							{ObjectMeta: metav1.ObjectMeta{Namespace: "capv-system", Name: "infrastructure-vsphere"}, Type: "InfrastructureProvider", Version: "v0.6.2", ProviderName: "vsphere"},
						}
						return nil
					})
					providerUpgradeClient.ApplyUpgradeCalls(func(applyUpgradeOptions *clusterctl.ApplyUpgradeOptions) error {
						applyUpgradeOptionsRecvd = *applyUpgradeOptions
						return nil
					})
				})
				It("should apply the providers version upgrade only for the outdated providers to the latest versions", func() {
					Expect(applyUpgradeOptionsRecvd.CoreProvider).Should(Equal("capi-system/cluster-api:v0.3.11"))
					Expect(len(applyUpgradeOptionsRecvd.BootstrapProviders)).Should(Equal(0))
					Expect(len(applyUpgradeOptionsRecvd.ControlPlaneProviders)).Should(Equal(0))
					Expect(applyUpgradeOptionsRecvd.InfrastructureProviders[0]).Should(Equal("capv-system/vsphere:v0.7.1"))
					Expect(providerUpgradeClient.ApplyUpgradeCallCount()).Should(Equal(1))
					Expect(err).NotTo(HaveOccurred())
				})
			})
			Context("When providers current versions are out dated and providers upgraded successfully", func() {
				var applyUpgradeOptionsRecvd clusterctl.ApplyUpgradeOptions
				BeforeEach(func() {
					regionalClusterClient.ListResourcesCalls(func(providers interface{}, options ...crtclient.ListOption) error {
						installedProviders, _ := providers.(*clusterctlv1.ProviderList)
						installedProviders.Items = []clusterctlv1.Provider{
							{ObjectMeta: metav1.ObjectMeta{Namespace: "capi-system", Name: "cluster-api"}, Type: "CoreProvider", Version: "v0.3.1", ProviderName: "cluster-api"},
							{ObjectMeta: metav1.ObjectMeta{Namespace: "capi-kubeadm-bootstrap-system", Name: "bootstrap-kubeadm"}, Type: "BootstrapProvider", Version: "v0.3.1", ProviderName: "kubeadm"},
							{ObjectMeta: metav1.ObjectMeta{Namespace: "capi-kubeadm-control-plane-system", Name: "control-plane-kubeadm"}, Type: "ControlPlaneProvider", Version: "v0.3.1", ProviderName: "kubeadm"},
							{ObjectMeta: metav1.ObjectMeta{Namespace: "capv-system", Name: "infrastructure-vsphere"}, Type: "InfrastructureProvider", Version: "v0.6.1", ProviderName: "vsphere"},
						}
						return nil
					})
					providerUpgradeClient.ApplyUpgradeCalls(func(applyUpgradeOptions *clusterctl.ApplyUpgradeOptions) error {
						applyUpgradeOptionsRecvd = *applyUpgradeOptions
						return nil
					})
				})
				It("should upgrade providers to the latest versions successfully", func() {
					Expect(applyUpgradeOptionsRecvd.CoreProvider).Should(Equal("capi-system/cluster-api:v0.3.11"))
					Expect(applyUpgradeOptionsRecvd.BootstrapProviders[0]).Should(Equal("capi-kubeadm-bootstrap-system/kubeadm:v0.3.11"))
					Expect(applyUpgradeOptionsRecvd.ControlPlaneProviders[0]).Should(Equal("capi-kubeadm-control-plane-system/kubeadm:v0.3.11"))
					Expect(applyUpgradeOptionsRecvd.InfrastructureProviders[0]).Should(Equal("capv-system/vsphere:v0.7.1"))
					Expect(providerUpgradeClient.ApplyUpgradeCallCount()).Should(Equal(1))
					Expect(err).NotTo(HaveOccurred())
				})
			})
		})

		Context("When validating compatibility before management cluster upgrade", func() {
			JustBeforeEach(func() {
				err = tkgClient.ValidateManagementClusterUpgradeVersionCompatibility(&upgradeClusterOptions, regionalClusterClient)
			})
			Context("When getting management cluster TKG version fails", func() {
				BeforeEach(func() {
					regionalClusterClient.GetManagementClusterTKGVersionReturns("", errors.New("fake GetManagementClusterTKGVersion error"))
				})
				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
					errString := `unable to get tkg version of management cluster "fake-cluster-name" in namespace "fake-namespace": fake GetManagementClusterTKGVersion error`
					Expect(err.Error()).To(ContainSubstring(errString))
				})
			})
			Context("When the management cluster TKG version is invalid ", func() {
				BeforeEach(func() {
					regionalClusterClient.GetManagementClusterTKGVersionReturns("InvalidTKGSemanticVersion", nil)
				})
				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("unable to parse management cluster version InvalidTKGSemanticVersion"))
				})
			})
			Context("When TKG version user trying to upgrade has major version greater than 2", func() {
				BeforeEach(func() {
					regionalClusterClient.GetManagementClusterTKGVersionReturns("v1.5.4", nil)
					tkgClient, err = CreateTKGClient("../fakes/config/config.yaml", testingDir, "../fakes/config/bom/tkg-bom-v3.0.0.yaml", 2*time.Millisecond)
					Expect(err).NotTo(HaveOccurred())
				})
				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("upgrading to beyond major version 2 is not yet supported."))
				})
			})
			Context("When the management cluster version is less than 1.6.x during v2.1.0 upgrade", func() {
				BeforeEach(func() {
					tkgClient, err = CreateTKGClient("../fakes/config/config.yaml", testingDir, "../fakes/config/bom/tkg-bom-v2.1.0.yaml", 2*time.Millisecond)
					Expect(err).NotTo(HaveOccurred())
					regionalClusterClient.GetManagementClusterTKGVersionReturns("v1.5.4", nil)
				})
				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("management cluster version must be 1.6.x to upgrade to 2.1.x"))
				})
			})
			Context("When management cluster current TKG version greater than TKG version user trying to upgrade", func() {
				BeforeEach(func() {
					regionalClusterClient.GetManagementClusterTKGVersionReturns("v1.2.2-rc.1", nil)
				})
				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("TKG version downgrade is not supported"))
				})
			})
			Context("When upgrading the TKG version (skipping a minor version) with skipPrompt set to true", func() {
				BeforeEach(func() {
					upgradeClusterOptions.SkipPrompt = true
					regionalClusterClient.GetManagementClusterTKGVersionReturns("v1.0.2-rc.1", nil)
				})
				It("should not return an error", func() {
					Expect(err).ToNot(HaveOccurred())
				})
			})
		})
	})
	Describe("When configuring AWS_AMI_ID for aws cluster upgrade", func() {
		const (
			fakeAWSClusterName      = "fakeAWSClusterName"
			fakeAWSClusterNamespace = "fakeAWSClusterNamespace"
			fakeAWSClusterKind      = "AWSCluster"
			fakeAWSAPIVersion       = "infrastructure.cluster.x-k8s.io/v1beta2"
		)
		BeforeEach(func() {
			setupBomFile("../fakes/config/bom/tkg-bom-v1.3.1.yaml", testingDir)
			setupBomFile("../fakes/config/bom/tkr-bom-v1.18.0+vmware.1-tkg.2.yaml", testingDir)
			upgradeClusterOptions.TkrVersion = "v1.18.0+vmware.1-tkg.2"
		})
		JustBeforeEach(func() {
			err = tkgClient.ConfigureAMIID(&upgradeClusterOptions, regionalClusterClient)
		})
		Context("When failed to get the cluster", func() {
			BeforeEach(func() {
				regionalClusterClient.GetResourceCalls(func(cluster interface{}, clusterName, namespace string, postVerify clusterclient.PostVerifyrFunc, pollOptions *clusterclient.PollOptions) error {
					if clusterName == upgradeClusterOptions.ClusterName && namespace == upgradeClusterOptions.Namespace {
						return errors.New("failed to get the cluster")
					}
					return nil
				})
			})
			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to get the cluster"))
			})
		})

		Context("When failed to get the InfrastructureRef", func() {
			BeforeEach(func() {
				regionalClusterClient.GetResourceCalls(func(cluster interface{}, clusterName, namespace string, postVerify clusterclient.PostVerifyrFunc, pollOptions *clusterclient.PollOptions) error {
					if clusterName == upgradeClusterOptions.ClusterName && namespace == upgradeClusterOptions.Namespace {
						clusterObj, ok := cluster.(*capi.Cluster)
						if !ok {
							return errors.New("not a cluster")
						}
						*clusterObj = capi.Cluster{
							Spec: capi.ClusterSpec{
								InfrastructureRef: nil,
							},
						}
					}
					return nil
				})
			})
			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("doest not have a InfrastructureRef"))
			})
		})

		Context("When failed to get the aws cluster", func() {
			BeforeEach(func() {
				regionalClusterClient.GetResourceCalls(func(cluster interface{}, clusterName, namespace string, postVerify clusterclient.PostVerifyrFunc, pollOptions *clusterclient.PollOptions) error {
					if clusterName == upgradeClusterOptions.ClusterName && namespace == upgradeClusterOptions.Namespace {
						clusterObj, ok := cluster.(*capi.Cluster)
						if !ok {
							return errors.New("not a cluster")
						}
						*clusterObj = capi.Cluster{
							Spec: capi.ClusterSpec{
								InfrastructureRef: &corev1.ObjectReference{
									APIVersion: fakeAWSAPIVersion,
									Kind:       fakeAWSClusterKind,
									Name:       fakeAWSClusterName,
									Namespace:  fakeAWSClusterNamespace,
								},
							},
						}
						return nil
					} else if clusterName == fakeAWSClusterName && namespace == fakeAWSClusterNamespace {
						awsclusterObj, ok := cluster.(*unstructured.Unstructured)
						if !ok {
							return errors.New("not a unstructured aws cluster")
						}

						if awsclusterObj.GetAPIVersion() == fakeAWSAPIVersion &&
							awsclusterObj.GetKind() == fakeAWSClusterKind {
							return errors.New("failed to get aws cluster")
						}
						return nil
					}
					return errors.Errorf("invalid clusterName %s and namespace %s", clusterName, namespace)
				})
			})
			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to get aws cluster"))
			})
		})

		Context("When failed to get the aws cluster spec", func() {
			BeforeEach(func() {
				regionalClusterClient.GetResourceCalls(func(cluster interface{}, clusterName, namespace string, postVerify clusterclient.PostVerifyrFunc, pollOptions *clusterclient.PollOptions) error {
					if clusterName == upgradeClusterOptions.ClusterName && namespace == upgradeClusterOptions.Namespace {
						clusterObj, ok := cluster.(*capi.Cluster)
						if !ok {
							return errors.New("not a cluster")
						}
						*clusterObj = capi.Cluster{
							Spec: capi.ClusterSpec{
								InfrastructureRef: &corev1.ObjectReference{
									APIVersion: fakeAWSAPIVersion,
									Kind:       fakeAWSClusterKind,
									Name:       fakeAWSClusterName,
									Namespace:  fakeAWSClusterNamespace,
								},
							},
						}
						return nil
					} else if clusterName == fakeAWSClusterName && namespace == fakeAWSClusterNamespace {
						awsclusterObj, ok := cluster.(*unstructured.Unstructured)
						if !ok {
							return errors.New("not a unstructured aws cluster")
						}

						if awsclusterObj.GetAPIVersion() == fakeAWSAPIVersion &&
							awsclusterObj.GetKind() == fakeAWSClusterKind {
							awsclusterObj.Object[constants.SPEC] = nil
						}
						return nil
					}
					return errors.Errorf("invalid clusterName %s and namespace %s", clusterName, namespace)
				})
			})
			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("has empty spec"))
			})
		})

		Context("When failed to get the aws cluster region", func() {
			BeforeEach(func() {
				regionalClusterClient.GetResourceCalls(func(cluster interface{}, clusterName, namespace string, postVerify clusterclient.PostVerifyrFunc, pollOptions *clusterclient.PollOptions) error {
					if clusterName == upgradeClusterOptions.ClusterName && namespace == upgradeClusterOptions.Namespace {
						clusterObj, ok := cluster.(*capi.Cluster)
						if !ok {
							return errors.New("not a cluster")
						}
						*clusterObj = capi.Cluster{
							Spec: capi.ClusterSpec{
								InfrastructureRef: &corev1.ObjectReference{
									APIVersion: fakeAWSAPIVersion,
									Kind:       fakeAWSClusterKind,
									Name:       fakeAWSClusterName,
									Namespace:  fakeAWSClusterNamespace,
								},
							},
						}
						return nil
					} else if clusterName == fakeAWSClusterName && namespace == fakeAWSClusterNamespace {
						awsclusterObj, ok := cluster.(*unstructured.Unstructured)
						if !ok {
							return errors.New("not a unstructured aws cluster")
						}

						if awsclusterObj.GetAPIVersion() == fakeAWSAPIVersion &&
							awsclusterObj.GetKind() == fakeAWSClusterKind {
							spec := map[string]interface{}{}
							awsclusterObj.Object[constants.SPEC] = spec
						}
						return nil
					}
					return errors.Errorf("invalid clusterName %s and namespace %s", clusterName, namespace)
				})
			})
			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("has empty region"))
			})
		})

		Context("When succeeded getting the AWS_AMI_ID", func() {
			BeforeEach(func() {
				regionalClusterClient.GetResourceCalls(func(cluster interface{}, clusterName, namespace string, postVerify clusterclient.PostVerifyrFunc, pollOptions *clusterclient.PollOptions) error {
					if clusterName == upgradeClusterOptions.ClusterName && namespace == upgradeClusterOptions.Namespace {
						clusterObj, ok := cluster.(*capi.Cluster)
						if !ok {
							return errors.New("not a cluster")
						}
						*clusterObj = capi.Cluster{
							Spec: capi.ClusterSpec{
								InfrastructureRef: &corev1.ObjectReference{
									APIVersion: fakeAWSAPIVersion,
									Kind:       fakeAWSClusterKind,
									Name:       fakeAWSClusterName,
									Namespace:  fakeAWSClusterNamespace,
								},
							},
						}
						return nil
					} else if clusterName == fakeAWSClusterName && namespace == fakeAWSClusterNamespace {
						awsclusterObj, ok := cluster.(*unstructured.Unstructured)
						if !ok {
							return errors.New("not a unstructured aws cluster")
						}

						if awsclusterObj.GetAPIVersion() == fakeAWSAPIVersion &&
							awsclusterObj.GetKind() == fakeAWSClusterKind {
							spec := map[string]interface{}{}
							spec["region"] = "us-west-2"
							awsclusterObj.Object[constants.SPEC] = spec
						}
						return nil
					}
					return errors.Errorf("invalid clusterName %s and namespace %s", clusterName, namespace)
				})
			})
			It("should get the AWS_AMI_ID along with OS_NAME, OS_VERSION, and OS_ARCH", func() {
				Expect(err).ToNot(HaveOccurred())
				amiID, err := tkgClient.TKGConfigReaderWriter().Get(constants.ConfigVariableAWSAMIID)
				Expect(err).ToNot(HaveOccurred())
				Expect(amiID).To(Equal("ami-03f483756fb3350c7"))
				osName, err := tkgClient.TKGConfigReaderWriter().Get(constants.ConfigVariableOSName)
				Expect(err).ToNot(HaveOccurred())
				Expect(osName).To(Equal("ubuntu"))
				osVersion, err := tkgClient.TKGConfigReaderWriter().Get(constants.ConfigVariableOSVersion)
				Expect(err).ToNot(HaveOccurred())
				Expect(osVersion).To(Equal("20.04"))
				osArch, err := tkgClient.TKGConfigReaderWriter().Get(constants.ConfigVariableOSArch)
				Expect(err).ToNot(HaveOccurred())
				Expect(osArch).To(Equal("amd64"))
			})
		})
	})
})
