package cfg

import "flag"

type Config struct {
	// the url of the management server
	CyclopsUrl string

	// the api key that is used to authenticate
	// this client
	CyclopsApiKey string

	// Kong's Admin Url
	KongAdminUrl string
}

func ReadConfig() *Config {
	cyclopsUrl := flag.String(
		"cyclops-url",
		"http://cyclops:8080",
		"Url of the cyclops management server",
	)

	cyclopsApiKey := flag.String(
		"cyclops-api-key",
		"-",
		"The API Key to the cyclops management server",
	)

	kongAdminUrl := flag.String(
		"kong-admin-url",
		"http://127.0.0.1:8001",
		"Kong's Admin Url",
	)

	flag.Parse()

	return &Config{
		CyclopsUrl:    *cyclopsUrl,
		CyclopsApiKey: *cyclopsApiKey,
		KongAdminUrl:  *kongAdminUrl,
	}
}
