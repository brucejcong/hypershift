package kas

import (
	hyperv1 "github.com/openshift/hypershift/api/v1alpha1"
	"github.com/openshift/hypershift/control-plane-operator/controllers/hostedcontrolplane/manifests"
	"github.com/openshift/hypershift/support/config"
	"github.com/openshift/hypershift/support/util"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestReconcileKubeAPIServerDeployment(t *testing.T) {

	// Setup expected values that are universal

	// Setup hypershift hosted control plane.
	targetNamespace := "test"
	kubeAPIDeployment := manifests.KASDeployment(targetNamespace)
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
		params           KubeAPIServerParams
		activeKey        []byte
		backupKey        []byte
	}{
		// empty deployment config
		{
			cm: corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-kube-api-config",
					Namespace: targetNamespace,
				},
				Data: map[string]string{"config.json": "test-data"},
			},
			deploymentConfig: config.DeploymentConfig{},
			params: KubeAPIServerParams{
				CloudProvider: "test-cloud-provider",
				APIServerPort: util.APIPortWithDefault(hcp, config.DefaultAPIServerPort),
			},
		},
	}
	for _, tc := range testCases {
		expectedMinReadySeconds := kubeAPIDeployment.Spec.MinReadySeconds
		err := ReconcileKubeAPIServerDeployment(kubeAPIDeployment, hcp, ownerRef, tc.deploymentConfig, tc.params.NamedCertificates(), tc.params.CloudProvider,
			tc.params.CloudProviderConfig, tc.params.CloudProviderCreds, tc.params.Images, &tc.cm, tc.params.AuditWebhookRef, tc.activeKey, tc.backupKey, tc.params.APIServerPort)
		assert.NoError(t, err)
		assert.Equal(t, expectedMinReadySeconds, kubeAPIDeployment.Spec.MinReadySeconds)
	}
}
