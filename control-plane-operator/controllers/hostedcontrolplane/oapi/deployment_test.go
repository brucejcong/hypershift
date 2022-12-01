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

func TestReconcileOpenshiftAPIServerDeployment(t *testing.T) {

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
		expectedTermGraceSeconds := oapiDeployment.Spec.Template.Spec.TerminationGracePeriodSeconds
		expectedMinReadySeconds := oapiDeployment.Spec.MinReadySeconds
		err := ReconcileDeployment(oapiDeployment, ownerRef, &tc.cm, tc.deploymentConfig, imageName, "socks5ProxyImage", config.DefaultEtcdURL, util.AvailabilityProberImageName, pointer.Int32(1234))
		assert.NoError(t, err)
		assert.Equal(t, expectedTermGraceSeconds, oapiDeployment.Spec.Template.Spec.TerminationGracePeriodSeconds)
		assert.Equal(t, expectedMinReadySeconds, oapiDeployment.Spec.MinReadySeconds)
	}
}

func TestReconcileOpenshiftOAuthAPIServerDeployment(t *testing.T) {
	// Setup expected values that are universal

	// Setup hypershift hosted control plane.
	targetNamespace := "test"
	oauthAPIDeployment := manifests.OpenShiftOAuthAPIServerDeployment(targetNamespace)
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
		deploymentConfig config.DeploymentConfig
		params           OAuthDeploymentParams
	}{
		// empty deployment config and oauth params
		{
			deploymentConfig: config.DeploymentConfig{},
			params:           OAuthDeploymentParams{},
		},
	}
	for _, tc := range testCases {
		expectedTermGraceSeconds := oauthAPIDeployment.Spec.Template.Spec.TerminationGracePeriodSeconds
		expectedMinReadySeconds := oauthAPIDeployment.Spec.MinReadySeconds
		err := ReconcileOAuthAPIServerDeployment(oauthAPIDeployment, ownerRef, &tc.params, pointer.Int32(1234))
		assert.NoError(t, err)
		assert.Equal(t, expectedTermGraceSeconds, oauthAPIDeployment.Spec.Template.Spec.TerminationGracePeriodSeconds)
		assert.Equal(t, expectedMinReadySeconds, oauthAPIDeployment.Spec.MinReadySeconds)
	}
}
