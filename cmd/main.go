package main

import (
	"os"

	"github.com/jetstack/cert-manager/pkg/acme/webhook/cmd"
	"github.com/uzxmx/cert-manager-webhook-alidns/solver"
)

func main() {
	groupName := os.Getenv("GROUP_NAME")
	if groupName == "" {
		panic("GROUP_NAME must be specified")
	}
	cmd.RunWebhookServer(groupName, &solver.AliDNSSolver{})
}
