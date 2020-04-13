package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/blang/semver"
	"github.com/heroku/docker-registry-client/registry"
	"github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/tools/components"
	"github.com/interconnectedcloud/qdr-operator/tools/constants"
	"github.com/interconnectedcloud/qdr-operator/tools/util"
	"github.com/interconnectedcloud/qdr-operator/version"
	"net/http"
	"strconv"

	oimagev1 "github.com/openshift/api/image/v1"
	routev1 "github.com/openshift/api/route/v1"
	csvv1 "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1alpha1"
	olmversion "github.com/operator-framework/operator-lifecycle-manager/pkg/lib/version"

	"github.com/tidwall/sjson"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	//"net/http"
	"os"
	//"sort"
	"strings"
	"time"
)

var (
	rh              = "Red Hat"
	maturity        = "alpha"
	major, minor, _ = components.MajorMinorMicro("7.6.0")
	csv             = csvSetting{

		Name:         "qdr",
		DisplayName:  "Red Hat Integration - AMQ Interconnect",
		OperatorName: "qdr-operator",
		CsvDir:       "qdr-operator",
		//Registry:      constants.QuayImageRegistry,
		Registry:  constants.RedHatImageRegistry,
		Context:   "amq7",
		ImageName: "amq-interconnect-operator",
		//Tag:          version.Version,
		Tag: constants.OPeratorImagTag,
	}
)

func main() {

	imageShaMap := map[string]string{}

	operatorName := csv.Name + "-operator"
	templateStruct := &csvv1.ClusterServiceVersion{}
	templateStruct.SetGroupVersionKind(csvv1.SchemeGroupVersion.WithKind("ClusterServiceVersion"))
	csvStruct := &csvv1.ClusterServiceVersion{}
	strategySpec := &csvStrategySpec{}
	json.Unmarshal(csvStruct.Spec.InstallStrategy.StrategySpecRaw, strategySpec)

	templateStrategySpec := &csvStrategySpec{}
	deployment := components.GetDeployment(csv.OperatorName, csv.Registry, csv.Context, csv.ImageName, csv.Tag, "Always")
	templateStrategySpec.Deployments = append(templateStrategySpec.Deployments, []csvDeployments{{Name: csv.OperatorName, Spec: deployment.Spec}}...)
	role := components.GetRole(csv.OperatorName)
	templateStrategySpec.Permissions = append(templateStrategySpec.Permissions, []csvPermissions{{ServiceAccountName: deployment.Spec.Template.Spec.ServiceAccountName, Rules: role.Rules}}...)

	// Re-serialize deployments and permissions into csv strategy.
	updatedStrat, err := json.Marshal(templateStrategySpec)
	if err != nil {
		panic(err)
	}
	templateStruct.Spec.InstallStrategy.StrategySpecRaw = updatedStrat
	templateStruct.Spec.InstallStrategy.StrategyName = "deployment"
	csvVersionedName := operatorName + ".v" + version.Version
	templateStruct.Name = csvVersionedName
	templateStruct.Namespace = "placeholder"
	descrip := "QDR Operator provides the ability to deploy and manages interior and edge QDR deployments automating resource creation and administration"
	var description = ""
	description += "The Qdr Operator is a lightweight [AMQP 1.0](https://www.amqp.org/) message router for building large, highly resilient messaging networks for hybrid cloud and IoT/edge deployments. Qdr transparently learns the addresses of messaging endpoints (such as clients, servers, and message brokers) and flexibly routes messages between them.\n\n"
	description += "### Core Capabilities\n\n"
	description += "* High throughput, low latency, shortest path message forwarding based on Layer 7 address routing mechanisms.\n\n"
	description += "\n"
	description += "* `Interior` mode deployments for any arbitrary topology of geographically-distributed and interconnected Qdrs.\n\n"
	description += "* `Edge` mode deployments for extremely large scale device endpoint connectivity.\n\n"
	description += "* Automatic message traffic re-routing when the network topology changes (resiliency without restrictions).\n\n"
	description += "* Flexible addressing schemes and delivery semantics (anycast, multicast, closest, balanced).\n\n"
	description += "* Integrated management with full support for the draft AMQP management specification.\n\n"
	description += "* Full-featured security capabilities for authentication, authorization, and policy-based resource access control.\n\n"
	description += "### Operator Features\n\n"
	description += "* **Flexible deployment plans** - Configurable deployment plans are available for `interior` and `edge` mode scenarios. These plans include all dependent resources.\n\n"
	description += "* **Placement directives** - Directives are provided to control how the pods should be scheduled.\n\n"
	description += "* **Connectivity configuration defaults** - Configuration defaults are automatically generated for listeners, connectors, and SSL/TLS setup.\n\n"
	description += "* **Exposes the service** - Integrated management of OpenShift Routes for exposed listener services for client, inter-router, and edge communications.\n\n"
	description += "* **Security certificate management** - Certificates are created and managed through integration with [jetstack cert-manager](https://github.com/jetstack/cert-manager).\n\n"
	description += "### Troubleshooting\n\n"
	description += "After deploying Interconnect, check any of the following to verify that it is operating correctly:\n\n"
	description += "* The Interconnect instance\n\n"
	description += "* The Deployment (or DaemonSet) instance\n\n"
	description += "* An individual pod for the Deployment (or DaemonSet)\n\n"
	description += "* A Route created for exposed services\n\n"
	description += "In addition, use `qdstat` commands to verify connectivity.\n\n"

	repository := "https://github.com/interconnectedcloud/qdr-operator"
	examples := []string{"{\n      \"apiVersion\":\"interconnectedcloud.github.io/v1alpha1\",\n      \"kind\":\"Interconnect\",\n      \"metadata\":{\n         \"name\":\"amq-interconnect\"\n      },\n      \"spec\":{\n         \"deploymentPlan\":{\n            \"size\":2,\n            \"role\":\"interior\",\n            \"placement\":\"Any\"\n         }\n      }\n   }"}
	templateStruct.SetAnnotations(
		map[string]string{
			"createdAt":           time.Now().Format("2006-01-02 15:04:05"),
			"containerImage":      deployment.Spec.Template.Spec.Containers[0].Image,
			"description":         descrip,
			"categories":          "Messaging",
			"certified":           "false",
			"capabilities":        "Basic Install",
			"repository":          repository,
			"support":             rh,
			"tectonic-visibility": "ocs",
			"alm-examples":        "[" + strings.Join(examples, ",") + "]",
		},
	)
	templateStruct.SetLabels(
		map[string]string{
			"operator-" + csv.Name: "true",
		},
	)

	templateStruct.Spec.Keywords = []string{"qdr", "operator"}
	var opVersion olmversion.OperatorVersion
	opVersion.Version = semver.MustParse(version.Version)
	templateStruct.Spec.Version = opVersion
	//templateStruct.Spec.Replaces = operatorName + ".v" + version.PriorVersion
	templateStruct.Spec.Description = description
	templateStruct.Spec.DisplayName = csv.DisplayName
	templateStruct.Spec.Maturity = maturity
	templateStruct.Spec.Maintainers = []csvv1.Maintainer{{Name: rh, Email: "customerservice@redhat.com"}}
	templateStruct.Spec.Provider = csvv1.AppLink{Name: rh}
	templateStruct.Spec.Links = []csvv1.AppLink{
		{Name: "Product Page", URL: "https://access.redhat.com/products/red-hat-amq/"},
		{Name: "Documentation", URL: "https://access.redhat.com/documentation/en-us/red_hat_amq/" + major + "." + minor + "/html/deploying_amq_interconnect_on_openshift/index"},
	}

	templateStruct.Spec.Icon = []csvv1.Icon{
		{
			Data:      "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",
			MediaType: "image/png",
		},
	}
	tLabels := map[string]string{
		"alm-owner-" + csv.Name: operatorName,
		"operated-by":           csvVersionedName,
	}
	templateStruct.Spec.Labels = tLabels
	templateStruct.Spec.Selector = &metav1.LabelSelector{MatchLabels: tLabels}
	templateStruct.Spec.InstallModes = []csvv1.InstallMode{
		{Type: csvv1.InstallModeTypeOwnNamespace, Supported: true},
		{Type: csvv1.InstallModeTypeSingleNamespace, Supported: true},
		{Type: csvv1.InstallModeTypeMultiNamespace, Supported: false},
		{Type: csvv1.InstallModeTypeAllNamespaces, Supported: false},
	}
	templateStruct.Spec.CustomResourceDefinitions.Owned = []csvv1.CRDDescription{
		{
			Version:     v1alpha1.SchemeGroupVersion.Version,
			Kind:        "Interconnect",
			DisplayName: "Red Hat Integration - AMQ Interconnect",
			Description: "An instance of AMQ Interconnect",
			Name:        "interconnects." + v1alpha1.SchemeGroupVersion.Group,
			Resources: []csvv1.APIResourceReference{

				{
					Kind:    "Service",
					Version: corev1.SchemeGroupVersion.String(),
				},
				{
					Kind:    "Deployment",
					Version: appsv1.SchemeGroupVersion.String(),
				},

				{
					Kind:    "ServiceAccount",
					Version: corev1.SchemeGroupVersion.String(),
				},

				{
					Kind:    "interconnects",
					Version: v1alpha1.SchemeGroupVersion.String(),
				},
				{
					Kind:    "RoleBinding",
					Version: rbacv1.SchemeGroupVersion.String(),
				},
				{
					Kind:    "Role",
					Version: rbacv1.SchemeGroupVersion.String(),
				},

				{
					Kind:    "pods",
					Version: corev1.SchemeGroupVersion.String(),
				},
				{
					Kind:    "configmaps",
					Version: corev1.SchemeGroupVersion.String(),
				},

				{
					Kind:    "Route",
					Version: routev1.GroupVersion.String(),
				},

				{
					Kind:    "Secret",
					Version: corev1.SchemeGroupVersion.String(),
				},
			},
			SpecDescriptors: []csvv1.SpecDescriptor{

				{
					Description:  "The number of instances to deploy",
					DisplayName:  "Size",
					Path:         "deploymentPlan.size",
					XDescriptors: []string{"urn:alm:descriptor:com.tectonic.ui:fieldGroup:deploymentPlan", "urn:alm:descriptor:com.tectonic.ui:podCount"},
				},

				{
					Description:  "Whether router is interior or edge",
					DisplayName:  "Role",
					Path:         "deploymentPlan.role",
					XDescriptors: []string{"urn:alm:descriptor:com.tectonic.ui:fieldGroup:deploymentPlan", "urn:alm:descriptor:com.tectonic.ui:label"},
				},
				{
					Description:  "Controls placement of pods and whether to use Deployment or DaemonSet, one of `Any`, `Every`, `AntiAffinity` or `Node`",
					DisplayName:  "Placement",
					Path:         "deploymentPlan.placement",
					XDescriptors: []string{"urn:alm:descriptor:com.tectonic.ui:fieldGroup:deploymentPlan", "urn:alm:descriptor:com.tectonic.ui:label"},
				},
			},
			StatusDescriptors: []csvv1.StatusDescriptor{
				{
					Description:  "The current revision of the Interconnect cluster",
					DisplayName:  "Revision Number",
					Path:         "revNumber",
					XDescriptors: []string{"urn:alm:descriptor:io.kubernetes.podCount"},
				},
				{
					Description:  "The current pods",
					DisplayName:  "Pods",
					Path:         "pods",
					XDescriptors: []string{"urn:alm:descriptor:io.kubernetes.podCount"},
				},
				{
					Description:  "The current conditions",
					DisplayName:  "Conditions",
					Path:         "conditions",
					XDescriptors: []string{"urn:alm:descriptor:io.kubernetes.conditions"},
				},
			},
		},
	}

	opMajor, opMinor, opMicro := components.MajorMinorMicro(version.Version)
	csvFile := "deploy/olm-catalog/" + csv.CsvDir + "/" + opMajor + "." + opMinor + "." + opMicro + "/" + csvVersionedName + ".clusterserviceversion.yaml"

	imageName, _, _ := components.GetImage(deployment.Spec.Template.Spec.Containers[0].Image)
	relatedImages := []image{}

	templateStruct.Annotations["certified"] = "false"
	deployFile := "deploy/operator.yaml"
	createFile(deployFile, deployment)
	roleFile := "deploy/role.yaml"
	createFile(roleFile, role)
	relatedImages = append(relatedImages, image{Name: imageName, Image: deployment.Spec.Template.Spec.Containers[0].Image})

	imageRef := constants.ImageRef{
		TypeMeta: metav1.TypeMeta{
			APIVersion: oimagev1.SchemeGroupVersion.String(),
			Kind:       "ImageStream",
		},
		Spec: constants.ImageRefSpec{
			Tags: []constants.ImageRefTag{
				{
					// Needs to match the component name for upstream and downstream.
					Name: "amq7/amq-interconnect-openshift-container",
					From: &corev1.ObjectReference{
						// Needs to match the image that is in your CSV that you want to replace.
						Name: deployment.Spec.Template.Spec.Containers[0].Image,
						Kind: "DockerImage",
					},
				},
			},
		},
	}

	//sort.Sort(sort.Reverse(sort.StringSlice(SupportedVersions)))

	imageRef.Spec.Tags = append(imageRef.Spec.Tags, constants.ImageRefTag{
		Name: constants.Interconnect17Component,
		From: &corev1.ObjectReference{
			Name: constants.Interconnect17ImageURL,
			Kind: "DockerImage",
		},
	})

	relatedImages = append(relatedImages, getRelatedImage(constants.Interconnect17ImageURL))

	if GetBoolEnv("DIGESTS") {

		for _, tagRef := range imageRef.Spec.Tags {

			if _, ok := imageShaMap[tagRef.From.Name]; !ok {
				imageShaMap[tagRef.From.Name] = ""
				imageName, imageTag, imageContext := components.GetImage(tagRef.From.Name)
				repo := imageContext + "/" + imageName

				digests, err := RetriveFromRedHatIO(repo, imageTag)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
				if len(digests) > 1 {
					imageShaMap[tagRef.From.Name] = strings.ReplaceAll(tagRef.From.Name, ":"+imageTag, "@"+digests[len(digests)-1])
				}
			}
		}
	}
	//not sure if we required mage-references file in the future So comment out for now.

	//imageFile := "deploy/olm-catalog/" + csv.CsvDir + "/" + opMajor + "." + opMinor + "." + opMicro + "/" + "image-references"
	//createFile(imageFile, imageRef)

	var templateInterface interface{}
	if len(relatedImages) > 0 {
		templateJSON, err := json.Marshal(templateStruct)
		if err != nil {
			fmt.Println(err)
		}
		result, err := sjson.SetBytes(templateJSON, "spec.relatedImages", relatedImages)
		if err != nil {
			fmt.Println(err)

		}
		if err = json.Unmarshal(result, &templateInterface); err != nil {
			fmt.Println(err)
		}
	} else {
		templateInterface = templateStruct
	}

	// find and replace images with SHAs where necessary
	templateByte, err := json.Marshal(templateInterface)
	if err != nil {
		fmt.Println(err)
	}
	for from, to := range imageShaMap {
		if to != "" {
			templateByte = bytes.ReplaceAll(templateByte, []byte(from), []byte(to))
		}
	}
	if err = json.Unmarshal(templateByte, &templateInterface); err != nil {
		fmt.Println(err)
	}
	createFile(csvFile, &templateInterface)

	packageFile := "deploy/olm-catalog/" + csv.CsvDir + "/" + opMajor + "." + opMinor + "." + opMicro + "/interconnectedcloud.package.yaml"
	p, err := os.Create(packageFile)
	defer p.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	pwr := bufio.NewWriter(p)
	pwr.WriteString("#! package-manifest: " + csvFile + "\n")
	packagedata := packageStruct{
		PackageName: operatorName,
		Channels: []channel{
			{
				Name:       maturity,
				CurrentCSV: csvVersionedName,
			},
		},
		DefaultChannel: maturity,
	}
	util.MarshallObject(packagedata, pwr)
	pwr.Flush()
}

type csvSetting struct {
	Name         string `json:"name"`
	DisplayName  string `json:"displayName"`
	OperatorName string `json:"operatorName"`
	CsvDir       string `json:"csvDir"`
	Registry     string `json:"repository"`
	Context      string `json:"context"`
	ImageName    string `json:"imageName"`
	Tag          string `json:"tag"`
}
type csvPermissions struct {
	ServiceAccountName string              `json:"serviceAccountName"`
	Rules              []rbacv1.PolicyRule `json:"rules"`
}
type csvDeployments struct {
	Name string                `json:"name"`
	Spec appsv1.DeploymentSpec `json:"spec,omitempty"`
}
type csvStrategySpec struct {
	Permissions        []csvPermissions                `json:"permissions"`
	Deployments        []csvDeployments                `json:"deployments"`
	ClusterPermissions []StrategyDeploymentPermissions `json:"clusterPermissions,omitempty"`
}

type StrategyDeploymentPermissions struct {
	ServiceAccountName string            `json:"serviceAccountName"`
	Rules              []rbac.PolicyRule `json:"rules"`
}
type channel struct {
	Name       string `json:"name"`
	CurrentCSV string `json:"currentCSV"`
}
type packageStruct struct {
	PackageName    string    `json:"packageName"`
	Channels       []channel `json:"channels"`
	DefaultChannel string    `json:"defaultChannel"`
}
type image struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

func getRelatedImage(imageURL string) image {
	imageName, _, _ := components.GetImage(imageURL)
	return image{
		Name:  imageName,
		Image: imageURL,
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func createFile(filepath string, obj interface{}) {
	f, err := os.Create(filepath)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	writer := bufio.NewWriter(f)
	util.MarshallObject(obj, writer)
	writer.Flush()
}

func GetBoolEnv(key string) bool {
	val := GetEnv(key, "false")
	ret, err := strconv.ParseBool(val)
	if err != nil {
		return false
	}
	return ret
}

func GetEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}
func RetriveFromRedHatIO(image string, imageTag string) ([]string, error) {

	url := "https://" + constants.RedHatImageRegistry

	username := "" // anonymous
	password := "" // anonymous

	if userToken := strings.Split(os.Getenv("REDHATIO_TOKEN"), ":"); len(userToken) > 1 {
		username = userToken[0]
		password = userToken[1]
	}
	hub, err := registry.New(url, username, password)
	digests := []string{}
	if err != nil {
		fmt.Println(err)
	} else {
		tags, err := hub.Tags(image)
		if err != nil {
			fmt.Println(err)
		}
		// do not follow redirects - this is critical so we can get the registry digest from Location in redirect response
		hub.Client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
		if _, exists := find(tags, imageTag); exists {
			req, err := http.NewRequest("GET", url+"/v2/"+image+"/manifests/"+imageTag, nil)
			if err != nil {
				fmt.Println(err)
			}
			req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")
			resp, err := hub.Client.Do(req)
			if err != nil {
				fmt.Println(err)
			}
			if resp != nil {
				defer resp.Body.Close()
			}
			if resp.StatusCode == 302 || resp.StatusCode == 301 {
				digestURL, err := resp.Location()
				if err != nil {
					fmt.Println(err)

				}

				if digestURL != nil {
					if url := strings.Split(digestURL.EscapedPath(), "/"); len(url) > 1 {
						digests = url

						return url, err
					}
				}
			}
		}

	}
	return digests, err
}
func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
