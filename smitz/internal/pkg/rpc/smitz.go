package rpc

import (
	"cjavellana.me/kong/smitz/internal/pkg/ipc"

	"context"
)

type Smitz struct {
	ipc.UnimplementedSmitzServer
}

func (s *Smitz) NodeInfo(ctx context.Context, in *ipc.NodeInfoRequest) (*ipc.NodeInfoResponse, error) {
	return &ipc.NodeInfoResponse{
		Plugins: &ipc.Plugins{
			AvailableOnServer: &ipc.Plugins_AvailableOnServer{
				PluginName: []string {"One", "Two"},
			},
			EnabledInCluster: &ipc.Plugins_EnabledInCluster{
				PluginName: []string {"Three", "Four"},
			},
		},
		Configuration: &ipc.Configuration{
			Name: []string {"One", "Two"},
		},
	}, nil
}