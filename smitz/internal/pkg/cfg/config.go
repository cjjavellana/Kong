package cfg

import "flag"

type Config struct {
	// the url of the management server
	CyclopsUrl string

	// the api key that is used to authenticate
	// this client
	CyclopsApiKey string
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

	flag.Parse()

	return &Config{
		CyclopsUrl: *cyclopsUrl,
		CyclopsApiKey: *cyclopsApiKey,
	}
}
