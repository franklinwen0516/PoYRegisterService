package main

import (
	"register_service/protos"
	"register_service/process"
	"context"
)

func (s *RegisterServiceImpl) RegisterWithBioKey(ctx context.Context,
	req *protos.BioRegisterRequset) (*protos.BioRegisterResponse, error) {
	var processor = process.NewRegisterWithImagesProcess()
	rsp := &protos.BioRegisterResponse{}
	err := processor.DoProcess(ctx, req, rsp)
	return rsp, err
}
