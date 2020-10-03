package rpc

import (
	"cjavellana.me/kong/smitz/internal/pkg/ipc"
	"cjavellana.me/kong/smitz/internal/pkg/kong"
	"context"
	"github.com/golang/protobuf/jsonpb"
	"strings"
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
	nodeInfoJson, _ := s.kongAdmin.NodeInfo()

	unmarshaler := jsonpb.Unmarshaler{}
	unmarshaler.AllowUnknownFields = true

	nodeInfoResponse := &ipc.NodeInfoResponse{}
	if err := unmarshaler.Unmarshal(strings.NewReader(nodeInfoJson), nodeInfoResponse); err != nil {
		return nil, err
	}

	return nodeInfoResponse, nil
}
