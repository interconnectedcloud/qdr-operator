package certificates

import (
	v1alpha1 "github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	cmv1alpha1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewSelfSignedIssuerForCR(m *v1alpha1.Qdr) *cmv1alpha1.Issuer {
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

func NewCAIssuerForCR(m *v1alpha1.Qdr, secret string) *cmv1alpha1.Issuer {
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

func NewSelfSignedCACertificateForCR(m *v1alpha1.Qdr) *cmv1alpha1.Certificate {
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

func NewCertificateForCR(m *v1alpha1.Qdr, profileName string) *cmv1alpha1.Certificate {
	cert := &cmv1alpha1.Certificate{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "certmanager.k8s.io/v1alpha1",
			Kind:       "Certificate",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-" + profileName + "-tls",
			Namespace: m.Namespace,
		},
		Spec: cmv1alpha1.CertificateSpec{
			SecretName: m.Name + "-" + profileName + "-tls",
			CommonName: m.Name + "." + m.Namespace + ".svc.cluster.local",
			DNSNames: []string{
				m.Name + "." + m.Namespace + ".svc.cluster.local",
			},
			IssuerRef: cmv1alpha1.ObjectReference{
				Name: m.Name + "-ca",
			},
		},
	}
	return cert
}

func NewCACertificateForCR(m *v1alpha1.Qdr, profileName string) *cmv1alpha1.Certificate {
	cert := &cmv1alpha1.Certificate{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "certmanager.k8s.io/v1alpha1",
			Kind:       "Certificate",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-" + profileName + "-ca",
			Namespace: m.Namespace,
		},
		Spec: cmv1alpha1.CertificateSpec{
			SecretName: m.Name + "-" + profileName + "-ca",
			CommonName: m.Name + "." + m.Namespace + ".svc.cluster.local",
			DNSNames: []string{
				m.Name + "." + m.Namespace + ".svc.cluster.local",
			},
			IsCA: true,
			IssuerRef: cmv1alpha1.ObjectReference{
				Name: m.Name + "-ca",
			},
		},
	}
	return cert
}
