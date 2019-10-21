package entities

type SslProfile struct {
	EntityCommon
	Ciphers            string `json:"ciphers"`
	Protocols          string `json:"protocols"`
	CaCertFile         string `json:"caCertFile"`
	CertFile           string `json:"certFile"`
	PrivateKeyFile     string `json:"privateKeyFile"`
	PasswordFile       string `json:"passwordFile"`
	Password           string `json:"password"`
	UidFormat          string `json:"uidFormat"`
	UidNameMappingFile string `json:"uidNameMappingFile"`
}

func (SslProfile) GetEntityId() string {
	return "sslProfile"
}
