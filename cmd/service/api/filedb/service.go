package filedb

import (
	"context"

	"github.com/fidelfly/fxgos/cmd/service"
	"github.com/fidelfly/fxgos/cmd/service/res"
)

const (
	ServiceName = "service.file"
)

func RegisterServer(server Service) {
	service.Register(ServiceName, server)
}

type Service interface {
	Save(ctx context.Context, name string, data []byte) (int64, error)
	Read(ctx context.Context, id int64) (*res.File, error)
}
