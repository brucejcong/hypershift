package oauth

import (
	"context"
	"github.com/openshift/hypershift/api"
	hyperv1 "github.com/openshift/hypershift/api/v1alpha1"
	"github.com/openshift/hypershift/control-plane-operator/hostedclusterconfigoperator/controllers/resources/manifests"
	"github.com/openshift/hypershift/support/config"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
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
		params           OAuthServerParams
	}{
		// empty deployment config
		{
			cm: corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-oauth-config",
					Namespace: targetNamespace,
				},
				Data: map[string]string{"config.yaml": "test-data"},
			},
			deploymentConfig: config.DeploymentConfig{},
			params: OAuthServerParams{
				AvailabilityProberImage: "test-availability-image",
				Socks5ProxyImage:        "test-socks-5-proxy-image",
			},
		},
	}
	for _, tc := range testCases {
		ctx := context.Background()
		fakeClient := fake.NewClientBuilder().WithScheme(api.Scheme).Build()
		expectedTermGraceSeconds := oauthDeployment.Spec.Template.Spec.TerminationGracePeriodSeconds
		expectedMinReadySeconds := oauthDeployment.Spec.MinReadySeconds
		err := ReconcileDeployment(ctx, fakeClient, oauthDeployment, ownerRef, &tc.cm, imageName, tc.deploymentConfig, tc.params.IdentityProviders(), tc.params.OauthConfigOverrides,
			tc.params.AvailabilityProberImage, pointer.Int32(1234), tc.params.NamedCertificates(), tc.params.Socks5ProxyImage, tc.params.NoProxy)
		assert.NoError(t, err)

		assert.Equal(t, expectedTermGraceSeconds, oauthDeployment.Spec.Template.Spec.TerminationGracePeriodSeconds)
		assert.Equal(t, expectedMinReadySeconds, oauthDeployment.Spec.MinReadySeconds)
	}
}
