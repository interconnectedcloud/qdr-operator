package certificates

import (
	v1alpha1 "github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/pkg/utils/configs"
	cmv1alpha1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"
	apiextv1b1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var (
	certmgr_detected *bool
	log              = logf.Log.WithName("certificates")
)

func DetectCertmgrIssuer() bool {
	// find certmanager issuer crd
	if certmgr_detected == nil {
		iscm := detectCertmgr()
		certmgr_detected = &iscm
	}
	return *certmgr_detected
}

func detectCertmgr() bool {
	config, err := config.GetConfig()
	if err != nil {
		log.Error(err, "Error getting config: %v")
		return false
	}

	// create a client set that includes crd schema
	extClient, err := apiextclientset.NewForConfig(config)
	if err != nil {
		log.Error(err, "Error getting ext client set: %v")
		return false
	}

	crd := &apiextv1b1.CustomResourceDefinition{}
	crd, err = extClient.ApiextensionsV1beta1().CustomResourceDefinitions().Get("issuers.certmanager.k8s.io", metav1.GetOptions{})
	if err != nil {
		log.Error(err, "Issuer crd for cert-manager not present, qdr-operator will be unable to request certificate generation")
		return false
	} else {
		log.Info("Detected certmanager issuer crd", "issuer", crd)
		return true
	}

}

func NewSelfSignedIssuerForCR(m *v1alpha1.Interconnect) *cmv1alpha1.Issuer {
	issuer := &cmv1alpha1.Issuer{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "certmanager.k8s.io/v1alpha1",
			Kind:       "Issuer",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-selfsigned",
			Namespace: m.Namespace,
		},
		Spec: cmv1alpha1.IssuerSpec{
			IssuerConfig: cmv1alpha1.IssuerConfig{
				SelfSigned: &cmv1alpha1.SelfSignedIssuer{},
			},
		},
	}
	return issuer
}

func NewCAIssuerForCR(m *v1alpha1.Interconnect, secret string) *cmv1alpha1.Issuer {
	issuer := &cmv1alpha1.Issuer{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "certmanager.k8s.io/v1alpha1",
			Kind:       "Issuer",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-ca",
			Namespace: m.Namespace,
		},
		Spec: cmv1alpha1.IssuerSpec{
			IssuerConfig: cmv1alpha1.IssuerConfig{
				CA: &cmv1alpha1.CAIssuer{
					SecretName: secret,
				},
			},
		},
	}
	return issuer
}

func NewCAIssuer(name string, namespace string, secret string) *cmv1alpha1.Issuer {
	issuer := &cmv1alpha1.Issuer{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "certmanager.k8s.io/v1alpha1",
			Kind:       "Issuer",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: cmv1alpha1.IssuerSpec{
			IssuerConfig: cmv1alpha1.IssuerConfig{
				CA: &cmv1alpha1.CAIssuer{
					SecretName: secret,
				},
			},
		},
	}
	return issuer
}

func NewSelfSignedCACertificateForCR(m *v1alpha1.Interconnect) *cmv1alpha1.Certificate {
	cert := &cmv1alpha1.Certificate{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "certmanager.k8s.io/v1alpha1",
			Kind:       "Certificate",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-selfsigned",
			Namespace: m.Namespace,
		},
		Spec: cmv1alpha1.CertificateSpec{
			SecretName: m.Name + "-selfsigned",
			CommonName: m.Name + "." + m.Namespace + ".svc.cluster.local",
			IsCA:       true,
			IssuerRef: cmv1alpha1.ObjectReference{
				Name: m.Name + "-selfsigned",
			},
		},
	}
	return cert
}

func issuerName(m *v1alpha1.Interconnect, name string) string {
	if name == "" {
		return m.Name + "-ca"
	} else {
		return name
	}

}

func NewCertificateForCR(m *v1alpha1.Interconnect, profileName string, certName string, issuer string) *cmv1alpha1.Certificate {
	hostNames := configs.GetInterconnectExposedHostnames(m, profileName)
	cert := &cmv1alpha1.Certificate{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "certmanager.k8s.io/v1alpha1",
			Kind:       "Certificate",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      certName,
			Namespace: m.Namespace,
		},
		Spec: cmv1alpha1.CertificateSpec{
			SecretName: certName,
			CommonName: m.Name,
			DNSNames:   hostNames,
			IssuerRef: cmv1alpha1.ObjectReference{
				Name: issuerName(m, issuer),
			},
		},
	}
	return cert
}

func NewCACertificateForCR(m *v1alpha1.Interconnect, name string) *cmv1alpha1.Certificate {
	cert := &cmv1alpha1.Certificate{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "certmanager.k8s.io/v1alpha1",
			Kind:       "Certificate",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: m.Namespace,
		},
		Spec: cmv1alpha1.CertificateSpec{
			SecretName: name,
			CommonName: name,
			IsCA:       true,
			IssuerRef: cmv1alpha1.ObjectReference{
				Name: m.Name + "-selfsigned",
			},
		},
	}
	return cert
}
