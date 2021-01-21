// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package paas

import (
	"github.com/spf13/pflag"
	"github.com/suse/carrier/cli/kubernetes"
	config2 "github.com/suse/carrier/cli/kubernetes/config"
	"github.com/suse/carrier/cli/paas/config"
	"github.com/suse/carrier/cli/paas/eirini"
	"github.com/suse/carrier/cli/paas/gitea"
	"github.com/suse/carrier/cli/paas/ui"
)

// Injectors from wire.go:

func BuildApp(flags *pflag.FlagSet, configOverrides func(*config.Config)) (*CarrierClient, func(), error) {
	configConfig, err := config.Load()
	if err != nil {
		return nil, nil, err
	}
	restConfig, err := config2.KubeConfig()
	if err != nil {
		return nil, nil, err
	}
	cluster, err := kubernetes.NewClusterFromClient(restConfig)
	if err != nil {
		return nil, nil, err
	}
	resolver := gitea.NewResolver(configConfig, cluster)
	client, err := gitea.NewGiteaClient(resolver)
	if err != nil {
		return nil, nil, err
	}
	uiUI := ui.NewUI()
	clientset, err := eirini.NewEiriniKubeClient(cluster)
	if err != nil {
		return nil, nil, err
	}
	carrierClient := &CarrierClient{
		giteaClient:   client,
		kubeClient:    cluster,
		ui:            uiUI,
		config:        configConfig,
		giteaResolver: resolver,
		eiriniClient:  clientset,
	}
	return carrierClient, func() {
	}, nil
}
