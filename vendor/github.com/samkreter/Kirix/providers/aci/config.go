package aci

import (
	"fmt"
	"io"

	"github.com/BurntSushi/toml"
)

type providerConfig struct {
	ResourceGroup   string
	Region          string
	OperatingSystem string
	CPU             string
	Memory          string
	CInstances      string
}

func (p *ACIProvider) loadConfig(r io.Reader) error {
	var config providerConfig
	if _, err := toml.DecodeReader(r, &config); err != nil {
		return err
	}
	p.region = config.Region
	p.resourceGroup = config.ResourceGroup

	// Default to 20 mcpu
	p.cpu = "20"
	if config.CPU != "" {
		p.cpu = config.CPU
	}
	// Default to 100Gi
	p.memory = "100Gi"
	if config.Memory != "" {
		p.memory = config.Memory
	}
	// Default to 20 pods
	p.cinstances = "20"
	if config.CInstances != "" {
		p.cinstances = config.CInstances
	}

	// Default to Linux if the operating system was not defined in the config.
	if config.OperatingSystem == "" {
		config.OperatingSystem = "Linux"
	} else {
		if config.OperatingSystem != "Linux" && config.OperatingSystem != "Windows" {
			return fmt.Errorf("Operating system not supported")
		}
	}

	p.operatingSystem = config.OperatingSystem
	return nil
}
