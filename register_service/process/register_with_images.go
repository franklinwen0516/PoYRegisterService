package process

import (
	"context"
	"register_service/protos"
)

type RegisterWithImagesProcess interface {
	DoProcess(ctx context.Context, req *protos.BioRegisterRequset,
		rsp *protos.BioRegisterResponse) error
}

type RegisterWithImagesProcessImpl struct {
}

var (
	NewRegisterWithImagesProcess = func() RegisterWithImagesProcess {
		return &RegisterWithImagesProcessImpl{}
	}
)

func (c *RegisterWithImagesProcessImpl) DoProcess(ctx context.Context,
	req *protos.BioRegisterRequset, rsp *protos.BioRegisterResponse) error {
	return nil
}
