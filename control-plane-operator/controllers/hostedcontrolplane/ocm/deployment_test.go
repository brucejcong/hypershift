package ocm

import (
	hyperv1 "github.com/openshift/hypershift/api/v1alpha1"
	"github.com/openshift/hypershift/control-plane-operator/controllers/hostedcontrolplane/manifests"
	"github.com/openshift/hypershift/support/config"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestReconcileOpenshiftControllerManagerDeployment(t *testing.T) {

	// Setup expected values that are universal
	imageName := "ocmImage"

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
		cm               corev1.ConfigMap
		deploymentConfig config.DeploymentConfig
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
		},
	}
	for _, tc := range testCases {
		var expectedTermGraceSeconds *int64 = nil
		var expectedMinReadySeconds int32 = 0
		err := ReconcileDeployment(ocmDeployment, ownerRef, imageName, &tc.cm, tc.deploymentConfig)
		assert.NoError(t, err)

		assert.Equal(t, expectedTermGraceSeconds, ocmDeployment.Spec.Template.Spec.TerminationGracePeriodSeconds)
		assert.Equal(t, expectedMinReadySeconds, ocmDeployment.Spec.MinReadySeconds)
	}
}
