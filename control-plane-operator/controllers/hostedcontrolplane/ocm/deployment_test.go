package ocm

import (
	hyperv1 "github.com/openshift/hypershift/api/v1alpha1"
	"github.com/openshift/hypershift/control-plane-operator/controllers/hostedcontrolplane/manifests"
	"github.com/openshift/hypershift/support/config"
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
		},
	}
	for _, tc := range testCases {
		err := ReconcileDeployment(ocmDeployment, ownerRef, "ocmImage", &tc.cm, tc.deploymentConfig)
		assert.NoError(t, err)
		assert.Equal(t, tc.expectedDeployStrategy, ocmDeployment.Spec.Strategy)
		ocmDeployment.Spec.Strategy.Type = "hello"
		assert.NotEqual(t, tc.expectedDeployStrategy, ocmDeployment.Spec.Strategy)

		//
		ocmDeployment.Spec.Template.Spec.TerminationGracePeriodSeconds = pointer.Int64(60)
		err = ReconcileDeployment(ocmDeployment, ownerRef, "ocmImage", &tc.cm, tc.deploymentConfig)
		assert.NoError(t, err)
		assert.Equal(t, pointer.Int64(60), ocmDeployment.Spec.Template.Spec.TerminationGracePeriodSeconds)
		assert.Equal(t, tc.expectedDeployStrategy, ocmDeployment.Spec.Strategy)
	}

}
