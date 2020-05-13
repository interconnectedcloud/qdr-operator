package configs

import (
	"bytes"
	"strconv"
	"text/template"

	v1alpha1 "github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/pkg/constants"
	"github.com/interconnectedcloud/qdr-operator/pkg/utils/openshift"
)

func isSslProfileDefined(m *v1alpha1.Interconnect, name string) bool {
	for _, profile := range m.Spec.SslProfiles {
		if profile.Name == name {
			return true
		}
	}
	return false
}

func isSslProfileUsed(m *v1alpha1.Interconnect, name string) bool {
	for _, listener := range m.Spec.Listeners {
		if listener.SslProfile == name {
			return true
		}
	}
	for _, listener := range m.Spec.InterRouterListeners {
		if listener.SslProfile == name {
			return true
		}
	}
	return false
}

func isDefaultSslProfileDefined(m *v1alpha1.Interconnect) bool {
	return isSslProfileDefined(m, "default")
}

func isDefaultSslProfileUsed(m *v1alpha1.Interconnect) bool {
	return isSslProfileUsed(m, "default")
}

func isInterRouterSslProfileDefined(m *v1alpha1.Interconnect) bool {
	return isSslProfileDefined(m, "inter-router")
}

func isInterRouterSslProfileUsed(m *v1alpha1.Interconnect) bool {
	return isSslProfileUsed(m, "inter-router")
}

func getExposedListeners(listeners []v1alpha1.Listener) []v1alpha1.Listener {
	exposedListeners := []v1alpha1.Listener{}
	for _, listener := range listeners {
		if listener.Expose {
			exposedListeners = append(exposedListeners, listener)
		}
	}
	return exposedListeners
}

func GetInterconnectExposedListeners(m *v1alpha1.Interconnect) []v1alpha1.Listener {
	listeners := []v1alpha1.Listener{}
	normal := getExposedListeners(m.Spec.Listeners)
	internal := getExposedListeners(m.Spec.InterRouterListeners)
	edge := getExposedListeners(m.Spec.EdgeListeners)
	listeners = append(listeners, normal...)
	listeners = append(listeners, internal...)
	listeners = append(listeners, edge...)
	return listeners
}

func GetInterconnectExposedHostnames(m *v1alpha1.Interconnect, profileName string) []string {
	var hostNames []string
	exposedListeners := GetInterconnectExposedListeners(m)
	dns := openshift.GetDnsConfig()

	for _, listener := range exposedListeners {
		if listener.SslProfile == profileName {
			target := listener.Name
			if target == "" {
				target = strconv.Itoa(int(listener.Port))
			}
			hostNames = append(hostNames, m.Name+"-"+target+"."+m.Namespace+"."+dns.Spec.BaseDomain)
		}
	}
	hostNames = append(hostNames, m.Name+"."+m.Namespace+".svc.cluster.local")

	return hostNames
}

func IsCaSecretNeeded(profile *v1alpha1.SslProfile) bool {
	// If the ca and credentials are in the same secret, don't need to mount it twice
	// If the credentials were generated and signed by the ca in the profile, then the credentials
	// secret will contain the ca public key, so the ca secret doesn't need to be mounted
	return profile.CaCert != profile.Credentials && !(profile.MutualAuth && profile.GenerateCredentials)
}

func SetInterconnectDefaults(m *v1alpha1.Interconnect, certMgrPresent bool) (bool, bool) {
	requestCert := false
	updateDefaults := false

	if m.Spec.DeploymentPlan.Size == 0 {
		m.Spec.DeploymentPlan.Size = 1
		updateDefaults = true
	}
	if m.Spec.DeploymentPlan.Role == "" {
		m.Spec.DeploymentPlan.Role = v1alpha1.RouterRoleInterior
		updateDefaults = true
	}
	if m.Spec.DeploymentPlan.Placement == "" {
		m.Spec.DeploymentPlan.Placement = v1alpha1.PlacementAny
		updateDefaults = true
	}
	if m.Spec.DeploymentPlan.LivenessPort == 0 {
		m.Spec.DeploymentPlan.LivenessPort = constants.HttpLivenessPort
		updateDefaults = true
	}

	if m.Spec.Users == "" {
		m.Spec.Users = m.Name + "-users"
	}

	if len(m.Spec.Addresses) == 0 {
		m.Spec.Addresses = append(m.Spec.Addresses, v1alpha1.Address{
			Prefix:       "closest",
			Distribution: "closest",
		})
		m.Spec.Addresses = append(m.Spec.Addresses, v1alpha1.Address{
			Prefix:       "multicast",
			Distribution: "multicast",
		})
		m.Spec.Addresses = append(m.Spec.Addresses, v1alpha1.Address{
			Prefix:       "unicast",
			Distribution: "closest",
		})
		m.Spec.Addresses = append(m.Spec.Addresses, v1alpha1.Address{
			Prefix:       "exclusive",
			Distribution: "closest",
		})
		m.Spec.Addresses = append(m.Spec.Addresses, v1alpha1.Address{
			Prefix:       "broadcast",
			Distribution: "multicast",
		})
	}

	if len(m.Spec.Listeners) == 0 {
		m.Spec.Listeners = append(m.Spec.Listeners, v1alpha1.Listener{
			Port: 5672,
		}, v1alpha1.Listener{
			Port:             8080,
			Http:             true,
			Expose:           true,
			AuthenticatePeer: true,
		})
		if certMgrPresent {
			m.Spec.Listeners = append(m.Spec.Listeners, v1alpha1.Listener{
				Port:       5671,
				SslProfile: "default",
			})
		}
		updateDefaults = true
	}
	if m.Spec.DeploymentPlan.Role == v1alpha1.RouterRoleInterior {
		if len(m.Spec.InterRouterListeners) == 0 {
			if certMgrPresent {
				m.Spec.InterRouterListeners = append(m.Spec.InterRouterListeners, v1alpha1.Listener{
					Port:             55671,
					SslProfile:       "inter-router",
					Expose:           true,
					AuthenticatePeer: true,
					SaslMechanisms:   "EXTERNAL",
				})
			} else {
				m.Spec.InterRouterListeners = append(m.Spec.InterRouterListeners, v1alpha1.Listener{
					Port: 55672,
				})
			}
			updateDefaults = true
		}
		if len(m.Spec.EdgeListeners) == 0 {
			m.Spec.EdgeListeners = append(m.Spec.EdgeListeners, v1alpha1.Listener{
				Port: 45672,
			})
			updateDefaults = true
		}
	}
	if !isDefaultSslProfileDefined(m) && isDefaultSslProfileUsed(m) {
		m.Spec.SslProfiles = append(m.Spec.SslProfiles, v1alpha1.SslProfile{
			Name:                "default",
			Credentials:         m.Name + "-default-credentials",
			GenerateCredentials: true,
		})
		updateDefaults = true
		requestCert = true
	}
	if !isInterRouterSslProfileDefined(m) && isInterRouterSslProfileUsed(m) {
		m.Spec.SslProfiles = append(m.Spec.SslProfiles, v1alpha1.SslProfile{
			Name:                "inter-router",
			Credentials:         m.Name + "-inter-router-credentials",
			CaCert:              m.Name + "-inter-router-ca",
			GenerateCredentials: true,
			GenerateCaCert:      true,
			MutualAuth:          true,
		})
		updateDefaults = true
		requestCert = true
	}
	if certMgrPresent {
		for i := range m.Spec.SslProfiles {
			if m.Spec.SslProfiles[i].Credentials == "" && m.Spec.SslProfiles[i].CaCert == "" {
				m.Spec.SslProfiles[i].Credentials = m.Name + "-" + m.Spec.SslProfiles[i].Name + "-credentials"
				m.Spec.SslProfiles[i].GenerateCredentials = true
				if m.Spec.SslProfiles[i].MutualAuth {
					m.Spec.SslProfiles[i].CaCert = m.Name + "-" + m.Spec.SslProfiles[i].Name + "-ca"
					m.Spec.SslProfiles[i].GenerateCaCert = true
				}
				updateDefaults = true
				requestCert = true
			} else if m.Spec.SslProfiles[i].Credentials == "" && m.Spec.SslProfiles[i].MutualAuth {
				m.Spec.SslProfiles[i].Credentials = m.Name + "-" + m.Spec.SslProfiles[i].Name + "-credentials"
				m.Spec.SslProfiles[i].GenerateCredentials = true
				updateDefaults = true
				requestCert = true
			} else if m.Spec.SslProfiles[i].CaCert == "" && m.Spec.SslProfiles[i].MutualAuth {
				m.Spec.SslProfiles[i].CaCert = m.Name + "-" + m.Spec.SslProfiles[i].Name + "-ca"
				m.Spec.SslProfiles[i].GenerateCaCert = true
				updateDefaults = true
				requestCert = true
			} else if m.Spec.SslProfiles[i].GenerateCredentials || m.Spec.SslProfiles[i].GenerateCaCert {
				requestCert = true
			}
		}
	}
	return requestCert && certMgrPresent, updateDefaults
}

func ConfigForInterconnect(m *v1alpha1.Interconnect) string {
	config := `
router {
    mode: {{.DeploymentPlan.Role}}
    id: ${HOSTNAME}
}
{{range .Listeners}}
listener {
    {{- if .Name}}
    name: {{.Name}}
    {{- end}}
    {{- if .Host}}
    host: {{.Host}}
    {{- else}}
    host: 0.0.0.0
    {{- end}}
    {{- if .Port}}
    port: {{.Port}}
    {{- end}}
    {{- if .LinkCapacity}}
    linkCapacity: {{.LinkCapacity}}
    {{- end}}
    {{- if .RouteContainer}}
    role: route-container
    {{- else }}
    role: normal
    {{- end}}
    {{- if .Http}}
    http: true
    {{- end}}
    {{- if .AuthenticatePeer}}
    authenticatePeer: true
    {{- end}}
    {{- if .SaslMechanisms}}
    saslMechanisms: {{.SaslMechanisms}}
    {{- end}}
    {{- if .SslProfile}}
    sslProfile: {{.SslProfile}}
    {{- end}}
}
{{- end}}
listener {
    name: health-and-stats
    port: {{.DeploymentPlan.LivenessPort}}
    http: true
    healthz: true
    metrics: true
    websockets: false
    httpRootDir: invalid
}
{{range .InterRouterListeners}}
listener {
    {{- if .Name}}
    name: {{.Name}}
    {{- end}}
    role: inter-router
    {{- if .Host}}
    host: {{.Host}}
    {{- else}}
    host: 0.0.0.0
    {{- end}}
    {{- if .Port}}
    port: {{.Port}}
    {{- end}}
    {{- if .Cost}}
    cost: {{.Cost}}
    {{- end}}
    {{- if .LinkCapacity}}
    linkCapacity: {{.LinkCapacity}}
    {{- end}}
    {{- if .SaslMechanisms}}
    saslMechanisms: {{.SaslMechanisms}}
    {{- end}}
    {{- if .AuthenticatePeer}}
    authenticatePeer: true
    {{- end}}
    {{- if .SslProfile}}
    sslProfile: {{.SslProfile}}
    {{- end}}
}
{{- end}}
{{range .EdgeListeners}}
listener {
    {{- if .Name}}
    name: {{.Name}}
    {{- end}}
    role: edge
    {{- if .Host}}
    host: {{.Host}}
    {{- else}}
    host: 0.0.0.0
    {{- end}}
    {{- if .Port}}
    port: {{.Port}}
    {{- end}}
    {{- if .Cost}}
    cost: {{.Cost}}
    {{- end}}
    {{- if .LinkCapacity}}
    linkCapacity: {{.LinkCapacity}}
    {{- end}}
    {{- if .SaslMechanisms}}
    saslMechanisms: {{.SaslMechanisms}}
    {{- end}}
    {{- if .AuthenticatePeer}}
    authenticatePeer: true
    {{- end}}
    {{- if .SslProfile}}
    sslProfile: {{.SslProfile}}
    {{- end}}
}
{{- end}}
{{range .SslProfiles}}
sslProfile {
   name: {{.Name}}
   {{- if .Credentials}}
   certFile: /etc/qpid-dispatch-certs/{{.Name}}/{{.Credentials}}/tls.crt
   privateKeyFile: /etc/qpid-dispatch-certs/{{.Name}}/{{.Credentials}}/tls.key
   {{- end}}
   {{- if .CaCert}}
       {{- if eq .CaCert .Credentials}}
   caCertFile: /etc/qpid-dispatch-certs/{{.Name}}/{{.CaCert}}/ca.crt
       {{- else if and .GenerateCredentials .MutualAuth}}
   caCertFile: /etc/qpid-dispatch-certs/{{.Name}}/{{.Credentials}}/ca.crt
       {{- else}}
   caCertFile: /etc/qpid-dispatch-certs/{{.Name}}/{{.CaCert}}/tls.crt
       {{- end}}
   {{- end}}
}
{{- end}}
{{range .Addresses}}
address {
    {{- if .Prefix}}
    prefix: {{.Prefix}}
    {{- end}}
    {{- if .Pattern}}
    pattern: {{.Pattern}}
    {{- end}}
    {{- if .Distribution}}
    distribution: {{.Distribution}}
    {{- end}}
    {{- if .Waypoint}}
    waypoint: {{.Waypoint}}
    {{- end}}
    {{- if .IngressPhase}}
    ingressPhase: {{.IngressPhase}}
    {{- end}}
    {{- if .EgressPhase}}
    egressPhase: {{.EgressPhase}}
    {{- end}}
    {{- if .Priority}}
    priority: {{.Priority}}
    {{- end}}
    {{- if .EnableFallback}}
    enableFallback: {{.EnableFallback}}
    {{- end}}
}
{{- end}}
{{range .AutoLinks}}
autoLink {
    {{- if .Address}}
    addr: {{.Address}}
    {{- end}}
    {{- if .Direction}}
    direction: {{.Direction}}
    {{- end}}
    {{- if .ContainerId}}
    containerId: {{.ContainerId}}
    {{- end}}
    {{- if .Connection}}
    connection: {{.Connection}}
    {{- end}}
    {{- if .ExternalAddress}}
    externalAddress: {{.ExternalAddress}}
    {{- end}}
    {{- if .Phase}}
    phase: {{.Phase}}
    {{- end}}
    {{- if .Fallback}}
    fallback: {{.Fallback}}
    {{- end}}
}
{{- end}}
{{range .LinkRoutes}}
linkRoute {
    {{- if .Prefix}}
    prefix: {{.Prefix}}
    {{- end}}
    {{- if .Pattern}}
    pattern: {{.Pattern}}
    {{- end}}
    {{- if .Direction}}
    direction: {{.Direction}}
    {{- end}}
    {{- if .Connection}}
    connection: {{.Connection}}
    {{- end}}
    {{- if .ContainerId}}
    containerId: {{.ContainerId}}
    {{- end}}
    {{- if .AddExternalPrefix}}
    addExternalPrefix: {{.AddExternalPrefix}}
    {{- end}}
    {{- if .DelExternalPrefix}}
    delExternalPrefix: {{.DelExternalPrefix}}
    {{- end}}
}
{{- end}}
{{range .Connectors}}
connector {
    {{- if .Name}}
    name: {{.Name}}
    {{- end}}
    {{- if .Host}}
    host: {{.Host}}
    {{- end}}
    {{- if .Port}}
    port: {{.Port}}
    {{- end}}
    {{- if .RouteContainer}}
    role: route-container
    {{- else}}
    role: normal
    {{- end}}
    {{- if .Cost}}
    cost: {{.Cost}}
    {{- end}}
    {{- if .LinkCapacity}}
    linkCapacity: {{.LinkCapacity}}
    {{- end}}
    {{- if .SslProfile}}
    sslProfile: {{.SslProfile}}
    {{- end}}
    {{- if eq .VerifyHostname false}}
    verifyHostname: false
    {{- end}}
}
{{- end}}
{{range .InterRouterConnectors}}
connector {
    {{- if .Name}}
    name: {{.Name}}
    {{- end}}
    role: inter-router
    {{- if .Host}}
    host: {{.Host}}
    {{- end}}
    {{- if .Port}}
    port: {{.Port}}
    {{- end}}
    {{- if .Cost}}
    cost: {{.Cost}}
    {{- end}}
    {{- if .LinkCapacity}}
    linkCapacity: {{.LinkCapacity}}
    {{- end}}
    {{- if .SslProfile}}
    sslProfile: {{.SslProfile}}
    {{- end}}
    {{- if eq .VerifyHostname false}}
    verifyHostname: false
    {{- end}}
}
{{- end}}
{{range .EdgeConnectors}}
connector {
    {{- if .Name}}
    name: {{.Name}}
    {{- end}}
    role: edge
    {{- if .Host}}
    host: {{.Host}}
    {{- end}}
    {{- if .Port}}
    port: {{.Port}}
    {{- end}}
    {{- if .Cost}}
    cost: {{.Cost}}
    {{- end}}
    {{- if .LinkCapacity}}
    linkCapacity: {{.LinkCapacity}}
    {{- end}}
    {{- if .SslProfile}}
    sslProfile: {{.SslProfile}}
    {{- end}}
    {{- if eq .VerifyHostname false}}
    verifyHostname: false
    {{- end}}
}
{{- end}}`
	var buff bytes.Buffer
	qdrconfig := template.Must(template.New("qdrconfig").Parse(config))
	qdrconfig.Execute(&buff, m.Spec)
	return buff.String()
}

func ConfigForSasl(m *v1alpha1.Interconnect) string {
	config := `
pwcheck_method: auxprop
auxprop_plugin: sasldb
sasldb_path: /tmp/qdrouterd.sasldb
`
	return config
}
