package main

import (
	"context"
	"register_service/process"
	"register_service/protos"
)

func (s *RegisterServiceImpl) RegisterWithBioKey(ctx context.Context,
	req *protos.BioRegisterRequset) (*protos.BioRegisterResponse, error) {
	var processor = process.NewRegisterWithImagesProcess()
	rsp := &protos.BioRegisterResponse{}
	err := processor.DoProcess(ctx, req, rsp)
	return rsp, err
}
