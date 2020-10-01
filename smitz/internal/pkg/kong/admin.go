package kong

import (
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	nodeInfo = "/"
	status   = "/status"
	service  = "/service"
)

type Admin struct {
	adminUrl string
}

func New(adminUrl string) *Admin {
	// remove trailing slash as our endpoint constants (see above)
	// already has leading slashes
	if strings.HasSuffix(adminUrl, "/") {
		adminUrl = adminUrl[:len(adminUrl)-1]
	}

	return &Admin{
		adminUrl: adminUrl,
	}
}

// see https://docs.konghq.com/2.1.x/admin-api/#retrieve-node-information
func (admin *Admin) NodeInfo() (string, error) {
	return get(admin.adminUrl + nodeInfo)
}

// see https://docs.konghq.com/2.1.x/admin-api/#retrieve-node-status
func (admin *Admin) GetStatus() (string, error) {
	return get(admin.adminUrl + status)
}

func get(url string) (string, error) {
	resp, err := http.Get(url)
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
