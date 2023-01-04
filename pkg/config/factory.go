package config

import onpremConfig "github.com/fastenhealth/fastenhealth-onprem/backend/pkg/config"

func Create() (onpremConfig.Interface, error) {
	config := new(configuration)
	if err := config.Init(); err != nil {
		return nil, err
	}
	return config, nil
}
