package openstack

import (
	hyperv1 "github.com/openshift/hypershift/api/hypershift/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

const (
	CloudConfigKey = "clouds.conf"
	CaKey          = "ca.pem"
	Provider       = "openstack"
)

// ReconcileCloudConfig reconciles as expected by Nodes Kubelet.
func ReconcileCloudConfig(secret *corev1.Secret, hcp *hyperv1.HostedControlPlane, credentialsSecret *corev1.Secret) error {
	if secret.Data == nil {
		secret.Data = map[string][]byte{}
	}
	config := string(credentialsSecret.Data[CloudConfigKey]) // TODO(emilien): Missing key handling

	// XXX(mdbooth): Don't hard-code 'openstack' cloud name. It should be in the platform config.
	config += `
[Global]
use-clouds=true
clouds-file=/etc/openstack/credentials/clouds.yaml
cloud=openstack`

	// FIXME(emilien): This is specific to CCM, we might want to have 2 versions.
	// FIXME(emilien): Is it really a good idea to have it here?
	// FIXME(emilien): How do we make it configurable?
	if hcp.Spec.Platform.OpenStack.CACertSecret != nil {
		config += "\nca-file = /etc/pki/ca-trust/extracted/pem/ca.pem\n"
	}

	config += "\n[LoadBalancer]\nmax-shared-lb = 1\nmanage-security-groups = true\n"

	secret.Data[CloudConfigKey] = []byte(config)
	return nil
}

// ReconcileTrustedCA reconciles as expected by Nodes Kubelet.
func ReconcileTrustedCA(cm *corev1.ConfigMap, hcp *hyperv1.HostedControlPlane, caConfigMap *corev1.Secret) error {
	if cm.Data == nil {
		cm.Data = map[string]string{}
	}
	cm.Data[CaKey] = string(caConfigMap.Data[CaKey]) // TODO(emilien): Missing key handling
	return nil
}
