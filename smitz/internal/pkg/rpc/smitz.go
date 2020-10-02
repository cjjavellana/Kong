package rpc

import (
	"cjavellana.me/kong/smitz/internal/pkg/ipc"
	"cjavellana.me/kong/smitz/internal/pkg/kong"
	"context"
)

type KongAdminProxy struct {
	ipc.UnimplementedAdminServer

	kongAdmin *kong.Admin
}

func New(adminUrl string) *KongAdminProxy {
	return &KongAdminProxy{
		kongAdmin: kong.New(adminUrl),
	}
}

func (s *KongAdminProxy) NodeInfo(ctx context.Context, in *ipc.NodeInfoRequest) (*ipc.NodeInfoResponse, error) {
	// _, _ := s.kongAdmin.NodeInfo()

	return &ipc.NodeInfoResponse{
		Plugins: &ipc.Plugins{
			EnabledInCluster: []string{},
			AvailableOnServer: map[string]bool{
				"grpc-web":       true,
				"correlation-id": true,
				"pre-function":   true,
				"cors":           true,
			},
		},
		Configuration: &ipc.Configuration{
			Plugins:                  []string{"bundled"},
			CassandraUsername:        "kong",
			CassandraReadConsistency: "ONE",
			NginxStreamDirectives: []*ipc.KeyValuePair{
				{
					Name:  "lua_shared_dict",
					Value: "stream_prometheus_metrics 5m",
				},
				{
					Name:  "ssl_prefer_server_ciphers",
					Value: "off",
				},
				{
					Name:  "ssl_protocols",
					Value: "TLSv1.2 TLSv1.3",
				},
			},
		},
	}, nil
}
