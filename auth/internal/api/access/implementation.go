package access

import (
	"context"
	"errors"
	desc "github.com/dimastephen/auth/grpc/pkg/access_v1"
	"github.com/dimastephen/auth/internal/service"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AccessImplementation struct {
	desc.UnimplementedAccessServer
	accessService service.AccessService
}

func NewImplementation(service service.AccessService) *AccessImplementation {
	return &AccessImplementation{accessService: service}
}

func (i *AccessImplementation) Check(ctx context.Context, request *desc.CheckRequest) (*emptypb.Empty, error) {
	endpoint := request.GetEndpointAddress()
	if endpoint == "" {
		err := errors.New("blank endpoint address")
		return nil, err
	}

	err := i.accessService.Check(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
