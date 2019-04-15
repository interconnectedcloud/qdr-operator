package openshift

import (
	"golang.org/x/text/language"
	"golang.org/x/text/search"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var (
	openshift_detected *bool
	log                = logf.Log.WithName("openshift")
)

func IsOpenShift() bool {
	if openshift_detected == nil {
		isos := detectOpenShift()
		openshift_detected = &isos
	}
	return *openshift_detected
}

// IsOpenShift checks for the OpenShift API
func detectOpenShift() bool {

	log.Info("Detect if OpenShift is running")

	config, err := config.GetConfig()
	if err != nil {
		log.Error(err, "Error getting config: %v")
		return false
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error(err, "Error getting client set: %v")
		return false
	}

	apiSchema, err := clientset.OpenAPISchema()
	if err != nil {
		log.Error(err, "Error getting api schema: %v")
		return false
	}
	return stringSearch(apiSchema.GetInfo().Title, "openshift")
}

func stringSearch(str string, substr string) bool {
	if start, _ := search.New(language.English, search.IgnoreCase).IndexString(str, substr); start == -1 {
		return false
	}
	return true
}
