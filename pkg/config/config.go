package config

import (
	"net/url"
)

type (
	Controller struct {
		KubeConfig    string  `json:"kubeconfig"`
		MasterAddress url.URL `json:"master-address"`
		Namespace     string  `json:"namespace"`
		Deployment    string  `json:"deployment"`
	}

	Agent struct {
		Address url.URL `json:"address"`
	}
)
