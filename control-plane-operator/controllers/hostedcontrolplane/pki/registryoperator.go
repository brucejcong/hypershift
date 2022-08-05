package pki

import (
	corev1 "k8s.io/api/core/v1"

	"github.com/openshift/hypershift/support/config"
)

const metricsHostname = "cluster-image-registry-operator"

func ReconcileRegistryOperatorServingCert(secret, ca *corev1.Secret, ownerRef config.OwnerRef) error {
	dnsNames := []string{
		metricsHostname,
		"localhost",
	}
	return reconcileSignedCertWithAddresses(secret, ca, ownerRef, metricsHostname, []string{"openshift"}, X509UsageClientServerAuth, dnsNames, nil)
}
