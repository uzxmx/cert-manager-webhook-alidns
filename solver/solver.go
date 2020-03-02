package solver

import (
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/jetstack/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/jetstack/cert-manager/pkg/issuer/acme/dns/util"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// AliDNSSolver is a solver to solve DNS challenge for AliCloud.
type AliDNSSolver struct {
	client *kubernetes.Clientset
}

// Name returns the solver name.
func (s *AliDNSSolver) Name() string {
	return "alidns"
}

// Initialize initializes the solver.
func (s *AliDNSSolver) Initialize(kubeClientConfig *rest.Config, stopCh <-chan struct{}) error {
	cl, err := kubernetes.NewForConfig(kubeClientConfig)
	if err != nil {
		return err
	}

	s.client = cl
	return nil
}

// Present presents DNS chanllenge to AliCloud.
func (s *AliDNSSolver) Present(ch *v1alpha1.ChallengeRequest) error {
	log.Printf("Presenting TXT record: %v %v\n", ch.ResolvedFQDN, ch.ResolvedZone)

	config, err := extractConfig(ch)
	if err != nil {
		return err
	}

	client, err := s.newAliDNSClient(ch, config)
	if err != nil {
		return err
	}

	request := alidns.CreateAddDomainRecordRequest()
	request.DomainName = util.UnFqdn(ch.ResolvedZone)
	request.Type = "TXT"
	request.RR = extractRR(ch.ResolvedFQDN, ch.ResolvedZone)
	request.Value = ch.Key
	request.TTL = requests.NewInteger(config.TTL)

	if _, err := client.AddDomainRecord(request); err != nil {
		return err
	}

	log.Printf("Presented TXT record %v\n", ch.ResolvedFQDN)
	return nil
}

// CleanUp deletes relevant TXT record from AliCloud.
func (s *AliDNSSolver) CleanUp(ch *v1alpha1.ChallengeRequest) error {
	log.Printf("Cleaning up TXT record: %v %v\n", ch.ResolvedFQDN, ch.ResolvedZone)

	config, err := extractConfig(ch)
	if err != nil {
		return err
	}

	client, err := s.newAliDNSClient(ch, config)
	if err != nil {
		return err
	}

	request := alidns.CreateDeleteSubDomainRecordsRequest()
	request.DomainName = util.UnFqdn(ch.ResolvedZone)
	request.RR = extractRR(ch.ResolvedFQDN, ch.ResolvedZone)
	request.Type = "TXT"
	if _, err := client.DeleteSubDomainRecords(request); err != nil {
		return err
	}

	log.Printf("Cleaned up TXT record: %v %v\n", ch.ResolvedFQDN, ch.ResolvedZone)
	return nil
}
