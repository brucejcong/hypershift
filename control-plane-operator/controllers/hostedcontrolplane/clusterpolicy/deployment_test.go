package clusterpolicy

import (
	hyperv1 "github.com/openshift/hypershift/api/v1alpha1"
	"github.com/openshift/hypershift/control-plane-operator/controllers/hostedcontrolplane/manifests"
	"github.com/openshift/hypershift/support/config"
	"github.com/openshift/hypershift/support/util"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"testing"
)

func TestReconcileClusterPolicyDeployment(t *testing.T) {

	imageName := "clusterPolicyImage"
	// Setup expected values that are universal

	// Setup hypershift hosted control plane.
	targetNamespace := "test"
	clusterPolicyDeployment := manifests.ClusterPolicyControllerDeployment(targetNamespace)
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
	}{
		// empty deployment config
		{
			deploymentConfig: config.DeploymentConfig{},
		},
	}
	for _, tc := range testCases {
		var expectedTermGraceSeconds *int64 = nil
		var expectedMinReadySeconds int32 = 0
		err := ReconcileDeployment(clusterPolicyDeployment, ownerRef, imageName, tc.deploymentConfig, util.AvailabilityProberImageName, pointer.Int32(1234))
		assert.NoError(t, err)
		assert.Equal(t, expectedTermGraceSeconds, clusterPolicyDeployment.Spec.Template.Spec.TerminationGracePeriodSeconds)
		assert.Equal(t, expectedMinReadySeconds, clusterPolicyDeployment.Spec.MinReadySeconds)
	}
}
