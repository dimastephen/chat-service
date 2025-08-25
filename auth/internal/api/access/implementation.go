package access

import (
	"context"
	"errors"
	"github.com/dimastephen/auth/internal/logger"
	"github.com/dimastephen/auth/internal/service"
	desc "github.com/dimastephen/auth/pkg/access_v1"
	"go.uber.org/zap"
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
	logger.Info("New request for check", zap.String("endpoint", endpoint))
	if endpoint == "" {
		err := errors.New("blank endpoint address")
		return nil, err
	}

	err := i.accessService.Check(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	logger.Info("Access granted for", zap.String("endpoint", endpoint))
	return &emptypb.Empty{}, nil
}
