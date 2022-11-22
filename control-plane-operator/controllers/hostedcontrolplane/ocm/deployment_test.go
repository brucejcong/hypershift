package ocm

import (
	hyperv1 "github.com/openshift/hypershift/api/v1alpha1"
	"github.com/openshift/hypershift/control-plane-operator/controllers/hostedcontrolplane/manifests"
	"github.com/openshift/hypershift/support/config"
	"github.com/openshift/hypershift/support/util"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	"testing"
)

func TestReconcileDeployments(t *testing.T) {

	// Setup expected values that are universal
	maxSurge := intstr.FromInt(1)
	maxUnavailable := intstr.FromInt(0)

	// Setup hypershift hosted control plane.
	targetNamespace := "test"
	ocmDeployment := manifests.OpenShiftControllerManagerDeployment(targetNamespace)
	hcp := &hyperv1.HostedControlPlane{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "hcp",
			Namespace: targetNamespace,
		},
	}
	hcp.Name = "name"
	hcp.Namespace = "namespace"
	ownerRef := config.OwnerRefFrom(hcp)

	testCases := []struct {
		cm                     corev1.ConfigMap
		deploymentConfig       config.DeploymentConfig
		expectedDeployStrategy appsv1.DeploymentStrategy
		expectedSelector       metav1.LabelSelector
		expectedObjMeta        metav1.ObjectMeta
		expectedSpec           corev1.PodSpec
	}{
		// empty deployment config
		{
			cm: corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ocm-config",
					Namespace: targetNamespace,
				},
				Data: map[string]string{"config.yaml": "test-data"},
			},
			deploymentConfig: config.DeploymentConfig{},
			expectedDeployStrategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxSurge:       &maxSurge,
					MaxUnavailable: &maxUnavailable,
				},
			},
			expectedSelector: metav1.LabelSelector{
				MatchLabels: openShiftControllerManagerLabels(),
			},
			expectedObjMeta: metav1.ObjectMeta{
				Labels: openShiftControllerManagerLabels(),
			},
			expectedSpec: corev1.PodSpec{
				AutomountServiceAccountToken: pointer.Bool(false),
				Containers: []corev1.Container{
					util.BuildContainer(ocmContainerMain(), buildOCMContainerMain("ocmImage")),
				},
				Volumes: []corev1.Volume{
					util.BuildVolume(ocmVolumeConfig(), buildOCMVolumeConfig),
					util.BuildVolume(ocmVolumeServingCert(), buildOCMVolumeServingCert),
					util.BuildVolume(ocmVolumeKubeconfig(), buildOCMVolumeKubeconfig),
				},
			},
		},
	}
	for _, tc := range testCases {
		err := ReconcileDeployment(ocmDeployment, ownerRef, "ocmImage", &tc.cm, tc.deploymentConfig)
		assert.NoError(t, err)
		assert.Equal(t, tc.expectedDeployStrategy, ocmDeployment.Spec.Strategy)
		assert.Equal(t, &tc.expectedSelector, ocmDeployment.Spec.Selector)

		configBytes, _ := tc.cm.Data[configKey]
		configHash := util.ComputeHash(configBytes)
		tc.expectedObjMeta.Annotations = map[string]string{
			configHashAnnotation: configHash,
		}
		assert.Equal(t, tc.expectedObjMeta, ocmDeployment.Spec.Template.ObjectMeta)
		assert.Equal(t, tc.expectedSpec, ocmDeployment.Spec.Template.Spec)

		// Check to see if other random values are changed.
		ocmDeployment.Spec.Template.Spec.TerminationGracePeriodSeconds = pointer.Int64(60)
		ocmDeployment.Spec.MinReadySeconds = int32(60)
		err = ReconcileDeployment(ocmDeployment, ownerRef, "ocmImage", &tc.cm, tc.deploymentConfig)
		assert.NoError(t, err)
		assert.Equal(t, pointer.Int64(60), ocmDeployment.Spec.Template.Spec.TerminationGracePeriodSeconds)
		assert.Equal(t, int32(60), ocmDeployment.Spec.MinReadySeconds)
	}
}
