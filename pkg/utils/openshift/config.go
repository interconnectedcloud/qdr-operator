package openshift

import (
	configv1 "github.com/openshift/api/config/v1"
	configv1client "github.com/openshift/client-go/config/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var (
	logcfg = logf.Log.WithName("openshift-config")
)

func GetDnsConfig() *configv1.DNS {
	dns := &configv1.DNS{}

	config, err := config.GetConfig()
	if err != nil {
		logcfg.Error(err, "Error getting config: %v")
	}
	openshiftClient, err := configv1client.NewForConfig(config)
	if err != nil {
		logcfg.Error(err, "Error getting openshift client set: %v")
	}

	dns, err = openshiftClient.ConfigV1().DNSes().Get("cluster", metav1.GetOptions{})
	if err != nil {
		logcfg.Info("Unable to get cluster base domain, qdr-operator will be unable to include host name in requested certficates for exposed listeners")
	}
	return dns
}
