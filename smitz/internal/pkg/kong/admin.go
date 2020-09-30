package kong

import (
	"io/ioutil"
	"net/http"
)

const (
	adminUrl = "http://127.0.0.1:8001"
	status   = adminUrl + "/status"
	service  = adminUrl + "/service"
)

func GetStatus() (string, error) {
	resp, err := http.Get(status)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
