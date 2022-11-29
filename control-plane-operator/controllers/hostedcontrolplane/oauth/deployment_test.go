package oauth

import (
	hyperv1 "github.com/openshift/hypershift/api/v1alpha1"
	"github.com/openshift/hypershift/control-plane-operator/hostedclusterconfigoperator/controllers/resources/manifests"
	"github.com/openshift/hypershift/support/config"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestReconcileOauthDeployment(t *testing.T) {

	// Setup expected values that are universal
	imageName := "oauthImage"

	// Setup hypershift hosted control plane.
	targetNamespace := "test"
	oauthDeployment := manifests.OAuthDeployment(targetNamespace)
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
		expectedTermGraceSeconds := oauthDeployment.Spec.Template.Spec.TerminationGracePeriodSeconds
		expectedMinReadySeconds := oauthDeployment.Spec.MinReadySeconds
		err := ReconcileDeployment(oauthDeployment, ownerRef, imageName, &tc.cm, tc.deploymentConfig)
		assert.NoError(t, err)

		assert.Equal(t, expectedTermGraceSeconds, oauthDeployment.Spec.Template.Spec.TerminationGracePeriodSeconds)
		assert.Equal(t, expectedMinReadySeconds, oauthDeployment.Spec.MinReadySeconds)
	}
}
