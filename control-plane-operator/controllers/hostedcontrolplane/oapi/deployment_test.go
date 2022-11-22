package oapi

import (
	hyperv1 "github.com/openshift/hypershift/api/v1alpha1"
	"github.com/openshift/hypershift/control-plane-operator/controllers/hostedcontrolplane/manifests"
	"github.com/openshift/hypershift/support/config"
	"github.com/openshift/hypershift/support/util"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"testing"
)

func TestReconcileDeployments(t *testing.T) {

	imageName := "oapiImage"
	// Setup expected values that are universal

	// Setup hypershift hosted control plane.
	targetNamespace := "test"
	oapiDeployment := manifests.OpenShiftAPIServerDeployment(targetNamespace)
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
					Name:      "test-oapi-config",
					Namespace: targetNamespace,
				},
				Data: map[string]string{"config.yaml": "test-data"},
			},
			deploymentConfig: config.DeploymentConfig{},
		},
	}
	for _, tc := range testCases {
		err := ReconcileDeployment(oapiDeployment, ownerRef, &tc.cm, tc.deploymentConfig, imageName, "socks5ProxyImage", config.DefaultEtcdURL, util.AvailabilityProberImageName, pointer.Int32(1234))
		assert.NoError(t, err)

		// Check to see if other random values are changed.
		oapiDeployment.Spec.Template.Spec.TerminationGracePeriodSeconds = pointer.Int64(60)
		oapiDeployment.Spec.MinReadySeconds = int32(60)

		err = ReconcileDeployment(oapiDeployment, ownerRef, &tc.cm, tc.deploymentConfig, imageName, "socks5ProxyImage", config.DefaultEtcdURL, util.AvailabilityProberImageName, pointer.Int32(1234))
		assert.NoError(t, err)
		assert.Equal(t, pointer.Int64(60), oapiDeployment.Spec.Template.Spec.TerminationGracePeriodSeconds)
		assert.Equal(t, int32(60), oapiDeployment.Spec.MinReadySeconds)
	}
}
