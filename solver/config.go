package solver

// AliDNSSolverConfig defines solver configuration.
type AliDNSSolverConfig struct {
	RegionID        string             `json:"regionId"`
	AccessKeyID     string             `json:"accessKeyId"`
	AccessKeySecret string             `json:"accessKeySecret"`
	AccessKeyRef    AliDNSAccessKeyRef `json:"accessKeyRef"`
	TTL             int                `json:"ttl"`
}

// AliDNSAccessKeyRef defines configuration when using kubernetes secret.
type AliDNSAccessKeyRef struct {
	SecretName         string `json:"name"`
	AccessKeyIDKey     string `json:"accessKeyIdKey"`
	AccessKeySecretKey string `json:"accessKeySecretKey"`
}
