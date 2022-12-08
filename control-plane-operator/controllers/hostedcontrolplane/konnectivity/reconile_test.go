package konnectivity

import (
	hyperv1 "github.com/openshift/hypershift/api/v1alpha1"
	"github.com/openshift/hypershift/control-plane-operator/controllers/hostedcontrolplane/manifests"
	"github.com/openshift/hypershift/support/config"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestReconcileKonnectivityAgentDeployment(t *testing.T) {

	imageName := "konnectivity-agent-image"
	// Setup expected values that are universal

	// Setup hypershift hosted control plane.
	targetNamespace := "test"
	konnectivityAgentDeployment := manifests.KonnectivityAgentDeployment(targetNamespace)
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
		ips              []string
	}{
		// empty deployment config
		{
			deploymentConfig: config.DeploymentConfig{},
			ips:              []string{"1.2.3.4"},
		},
	}
	for _, tc := range testCases {
		var expectedTermGraceSeconds *int64 = nil
		var expectedMinReadySeconds int32 = 0
		err := ReconcileAgentDeployment(konnectivityAgentDeployment, ownerRef, tc.deploymentConfig, imageName, tc.ips)
		assert.NoError(t, err)
		assert.Equal(t, expectedTermGraceSeconds, konnectivityAgentDeployment.Spec.Template.Spec.TerminationGracePeriodSeconds)
		assert.Equal(t, expectedMinReadySeconds, konnectivityAgentDeployment.Spec.MinReadySeconds)
	}
}
