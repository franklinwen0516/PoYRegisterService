package process

import (
	"context"
	"register_service/db"
	"register_service/localutil"
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
	if len(req.AccountPublicKey) == 0 {
		errMsg := "Invalid public key"
		return c.doResponseExp(ctx, req, rsp,
			int32(protos.ERR_CODE_CODE_ERR_MISSING_PARAM), errMsg)
	}
	// TODO 调用verify接口，判断是否已经被注册
	if "verify" != "verify" {
		errMsg := "Verify fail"
		return c.doResponseExp(ctx, req, rsp,
			int32(protos.ERR_CODE_CODE_ERR_MISSING_PARAM), errMsg)
	}
	db.UserRegisterInfoInstance.UserDataWrite(req.AccountPublicKey, req.FacialImages)
	return c.doResponse(ctx, req, rsp)
}

func (c *RegisterWithImagesProcessImpl) doResponse(ctx context.Context,
	req *protos.BioRegisterRequset, rsp *protos.BioRegisterResponse) error {
	rsp.Header = &protos.CommonRspHeader{
		Ret:    0,
		Reason: "Succ",
	}
	localutil.UserRegisterLog.Infof("register success, req: %+v, rsp: %+v\n", req, rsp)
	return nil
}

func (c *RegisterWithImagesProcessImpl) doResponseExp(ctx context.Context,
	req *protos.BioRegisterRequset, rsp *protos.BioRegisterResponse,
	ret int32, errMsg string) error {
	rsp.Header = &protos.CommonRspHeader{
		Ret:    ret,
		Reason: errMsg,
	}
	localutil.UserRegisterLog.Errorf(
		"register fail, req: %+v, rsp: %+v, ret:%v, erMsg:%v\n",
		req, rsp, ret, errMsg)
	return nil
}
