package main

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type env struct {
	Port uint16 `envconfig:"PORT" default:"9080"`
	ServiceCAPIEndpoint string `envconfig:"SERVICE_C_ENDPOINT" default:"service-c.namespace-c.svc.cluster.local:5000"`
}

func newEnv() (*env, error) {
	var e env
	if err := envconfig.Process("", &e); err != nil {
		return nil, fmt.Errorf("failed to process envconfig: %w", err)
	}

	return &e, nil
}
