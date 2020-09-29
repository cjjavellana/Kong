package cyclops

import (
	"bytes"
	"cjavellana.me/kong/smitz/internal/pkg/cfg"
	"cjavellana.me/kong/smitz/internal/pkg/system"
	"fmt"
	"net/http"
)

type Cyclops struct {
	cyclopsUrl string
	cyclopsApiKey string
}

func New(config *cfg.Config) *Cyclops {
	return &Cyclops {
		cyclopsUrl: config.CyclopsUrl,
		cyclopsApiKey: config.CyclopsApiKey,
	}
}

// This package defines the interface that smitz uses to communicate
// to its management server (cyclops)

// Registers this node (where smitz is running) to the management server (cyclops)
func (c *Cyclops) Register() error {
	ip, err := system.HostIP()
	if err != nil {
		return err
	}

	err = c.sendToCyclops(ip)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cyclops) sendToCyclops(hostIP string) error {
	payload := fmt.Sprintf(`{"hostIP": "%s"}`, hostIP)
	var jsonStr = []byte(payload)
	req, err := http.NewRequest("POST", c.cyclopsUrl, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Set("X-API-Key", c.cyclopsApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	_ = resp.Body.Close()

	return nil
}