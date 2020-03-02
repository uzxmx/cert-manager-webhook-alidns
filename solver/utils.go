package solver

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/jetstack/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/jetstack/cert-manager/pkg/issuer/acme/dns/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func extractRR(fqdn, domain string) string {
	name := util.UnFqdn(fqdn)
	if idx := strings.Index(name, "."+domain); idx != -1 {
		return name[:idx]
	}
	return name
}

func extractConfig(ch *v1alpha1.ChallengeRequest) (*AliDNSSolverConfig, error) {
	config := &AliDNSSolverConfig{}

	if ch.Config == nil {
		return config, nil
	}

	if err := json.Unmarshal(ch.Config.Raw, config); err != nil {
		return config, fmt.Errorf("failed to decode solver config %v", err)
	}

	return config, nil
}

func (s *AliDNSSolver) newAliDNSClient(ch *v1alpha1.ChallengeRequest, config *AliDNSSolverConfig) (*alidns.Client, error) {
	accessKeyID := config.AccessKeyID
	accessKeySecret := config.AccessKeySecret

	if len(accessKeySecret) == 0 && len(accessKeyID) == 0 {
		ref := config.AccessKeyRef
		if len(ref.SecretName) == 0 {
			return nil, errors.New("Referenced secret name is required")
		}
		if len(ref.AccessKeyIDKey) == 0 {
			return nil, errors.New("AccessKeyIDKey is required")
		}
		if len(ref.AccessKeySecretKey) == 0 {
			return nil, errors.New("AccessKeySecretKey is required")
		}

		secret, err := s.client.CoreV1().Secrets(ch.ResourceNamespace).Get(ref.SecretName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		data, ok := secret.Data[ref.AccessKeyIDKey]
		if !ok {
			return nil, errors.New("failed to get access key id")
		}
		accessKeyID = string(data)

		data, ok = secret.Data[ref.AccessKeySecretKey]
		if !ok {
			return nil, errors.New("failed to get access key secret")
		}
		accessKeySecret = string(data)
	}

	if len(accessKeyID) == 0 || len(accessKeySecret) == 0 {
		return nil, errors.New("accessKeyID or accessKeySecret cannot be empty")
	}

	client, err := alidns.NewClientWithAccessKey(
		config.RegionID,
		accessKeyID,
		accessKeySecret,
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}
