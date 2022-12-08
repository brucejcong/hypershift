package scheduler

import (
	hyperv1 "github.com/openshift/hypershift/api/v1alpha1"
	"github.com/openshift/hypershift/control-plane-operator/controllers/hostedcontrolplane/manifests"
	"github.com/openshift/hypershift/support/config"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"testing"
)

func TestReconcileSchedulerDeployment(t *testing.T) {

	// Setup expected values that are universal
	imageName := "schedulerImage"

	// Setup hypershift hosted control plane.
	targetNamespace := "test"
	schedulerDeployment := manifests.SchedulerDeployment(targetNamespace)
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
		params           KubeSchedulerParams
		featureGates     []string
		ciphers          []string
	}{
		// empty deployment config
		{
			deploymentConfig: config.DeploymentConfig{},
			params: KubeSchedulerParams{
				DisableProfiling:        false,
				AvailabilityProberImage: "availability-prober",
			},
		},
	}
	for _, tc := range testCases {
		var expectedTermGraceSeconds *int64 = nil
		var expectedMinReadySeconds int32 = 0
		err := ReconcileDeployment(schedulerDeployment, ownerRef, tc.deploymentConfig, imageName,
			tc.params.FeatureGates(), tc.params.SchedulerPolicy(), tc.params.AvailabilityProberImage,
			pointer.Int32(1234), tc.params.CipherSuites(), tc.params.MinTLSVersion(), tc.params.DisableProfiling)
		assert.NoError(t, err)
		assert.Equal(t, expectedTermGraceSeconds, schedulerDeployment.Spec.Template.Spec.TerminationGracePeriodSeconds)
		assert.Equal(t, expectedMinReadySeconds, schedulerDeployment.Spec.MinReadySeconds)
	}
}
