// Code generated by Kitex v0.10.3. DO NOT EDIT.

package register

import (
	"context"
	"errors"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	streaming "github.com/cloudwego/kitex/pkg/streaming"
	proto "google.golang.org/protobuf/proto"
	ucenter "grpc_common/kitex_gen/ucenter"
)

var errInvalidMessageType = errors.New("invalid message type for service method handler")

var serviceMethods = map[string]kitex.MethodInfo{
	"registerByPhone": kitex.NewMethodInfo(
		registerByPhoneHandler,
		newRegisterByPhoneArgs,
		newRegisterByPhoneResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingUnary),
	),
	"sendCode": kitex.NewMethodInfo(
		sendCodeHandler,
		newSendCodeArgs,
		newSendCodeResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingUnary),
	),
}

var (
	registerServiceInfo                = NewServiceInfo()
	registerServiceInfoForClient       = NewServiceInfoForClient()
	registerServiceInfoForStreamClient = NewServiceInfoForStreamClient()
)

// for server
func serviceInfo() *kitex.ServiceInfo {
	return registerServiceInfo
}

// for stream client
func serviceInfoForStreamClient() *kitex.ServiceInfo {
	return registerServiceInfoForStreamClient
}

// for client
func serviceInfoForClient() *kitex.ServiceInfo {
	return registerServiceInfoForClient
}

// NewServiceInfo creates a new ServiceInfo containing all methods
func NewServiceInfo() *kitex.ServiceInfo {
	return newServiceInfo(false, true, true)
}

// NewServiceInfo creates a new ServiceInfo containing non-streaming methods
func NewServiceInfoForClient() *kitex.ServiceInfo {
	return newServiceInfo(false, false, true)
}
func NewServiceInfoForStreamClient() *kitex.ServiceInfo {
	return newServiceInfo(true, true, false)
}

func newServiceInfo(hasStreaming bool, keepStreamingMethods bool, keepNonStreamingMethods bool) *kitex.ServiceInfo {
	serviceName := "Register"
	handlerType := (*ucenter.Register)(nil)
	methods := map[string]kitex.MethodInfo{}
	for name, m := range serviceMethods {
		if m.IsStreaming() && !keepStreamingMethods {
			continue
		}
		if !m.IsStreaming() && !keepNonStreamingMethods {
			continue
		}
		methods[name] = m
	}
	extra := map[string]interface{}{
		"PackageName": "ucenter",
	}
	if hasStreaming {
		extra["streaming"] = hasStreaming
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Protobuf,
		KiteXGenVersion: "v0.10.3",
		Extra:           extra,
	}
	return svcInfo
}

func registerByPhoneHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(ucenter.RegReq)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(ucenter.Register).RegisterByPhone(ctx, req)
		if err != nil {
			return err
		}
		return st.SendMsg(resp)
	case *RegisterByPhoneArgs:
		success, err := handler.(ucenter.Register).RegisterByPhone(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*RegisterByPhoneResult)
		realResult.Success = success
		return nil
	default:
		return errInvalidMessageType
	}
}
func newRegisterByPhoneArgs() interface{} {
	return &RegisterByPhoneArgs{}
}

func newRegisterByPhoneResult() interface{} {
	return &RegisterByPhoneResult{}
}

type RegisterByPhoneArgs struct {
	Req *ucenter.RegReq
}

func (p *RegisterByPhoneArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(ucenter.RegReq)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *RegisterByPhoneArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *RegisterByPhoneArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *RegisterByPhoneArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, nil
	}
	return proto.Marshal(p.Req)
}

func (p *RegisterByPhoneArgs) Unmarshal(in []byte) error {
	msg := new(ucenter.RegReq)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var RegisterByPhoneArgs_Req_DEFAULT *ucenter.RegReq

func (p *RegisterByPhoneArgs) GetReq() *ucenter.RegReq {
	if !p.IsSetReq() {
		return RegisterByPhoneArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *RegisterByPhoneArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *RegisterByPhoneArgs) GetFirstArgument() interface{} {
	return p.Req
}

type RegisterByPhoneResult struct {
	Success *ucenter.RegRes
}

var RegisterByPhoneResult_Success_DEFAULT *ucenter.RegRes

func (p *RegisterByPhoneResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(ucenter.RegRes)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *RegisterByPhoneResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *RegisterByPhoneResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *RegisterByPhoneResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, nil
	}
	return proto.Marshal(p.Success)
}

func (p *RegisterByPhoneResult) Unmarshal(in []byte) error {
	msg := new(ucenter.RegRes)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *RegisterByPhoneResult) GetSuccess() *ucenter.RegRes {
	if !p.IsSetSuccess() {
		return RegisterByPhoneResult_Success_DEFAULT
	}
	return p.Success
}

func (p *RegisterByPhoneResult) SetSuccess(x interface{}) {
	p.Success = x.(*ucenter.RegRes)
}

func (p *RegisterByPhoneResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *RegisterByPhoneResult) GetResult() interface{} {
	return p.Success
}

func sendCodeHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(ucenter.CodeReq)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(ucenter.Register).SendCode(ctx, req)
		if err != nil {
			return err
		}
		return st.SendMsg(resp)
	case *SendCodeArgs:
		success, err := handler.(ucenter.Register).SendCode(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*SendCodeResult)
		realResult.Success = success
		return nil
	default:
		return errInvalidMessageType
	}
}
func newSendCodeArgs() interface{} {
	return &SendCodeArgs{}
}

func newSendCodeResult() interface{} {
	return &SendCodeResult{}
}

type SendCodeArgs struct {
	Req *ucenter.CodeReq
}

func (p *SendCodeArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(ucenter.CodeReq)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *SendCodeArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *SendCodeArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *SendCodeArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, nil
	}
	return proto.Marshal(p.Req)
}

func (p *SendCodeArgs) Unmarshal(in []byte) error {
	msg := new(ucenter.CodeReq)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var SendCodeArgs_Req_DEFAULT *ucenter.CodeReq

func (p *SendCodeArgs) GetReq() *ucenter.CodeReq {
	if !p.IsSetReq() {
		return SendCodeArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *SendCodeArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *SendCodeArgs) GetFirstArgument() interface{} {
	return p.Req
}

type SendCodeResult struct {
	Success *ucenter.NoRes
}

var SendCodeResult_Success_DEFAULT *ucenter.NoRes

func (p *SendCodeResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(ucenter.NoRes)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *SendCodeResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *SendCodeResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *SendCodeResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, nil
	}
	return proto.Marshal(p.Success)
}

func (p *SendCodeResult) Unmarshal(in []byte) error {
	msg := new(ucenter.NoRes)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *SendCodeResult) GetSuccess() *ucenter.NoRes {
	if !p.IsSetSuccess() {
		return SendCodeResult_Success_DEFAULT
	}
	return p.Success
}

func (p *SendCodeResult) SetSuccess(x interface{}) {
	p.Success = x.(*ucenter.NoRes)
}

func (p *SendCodeResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *SendCodeResult) GetResult() interface{} {
	return p.Success
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) RegisterByPhone(ctx context.Context, Req *ucenter.RegReq) (r *ucenter.RegRes, err error) {
	var _args RegisterByPhoneArgs
	_args.Req = Req
	var _result RegisterByPhoneResult
	if err = p.c.Call(ctx, "registerByPhone", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) SendCode(ctx context.Context, Req *ucenter.CodeReq) (r *ucenter.NoRes, err error) {
	var _args SendCodeArgs
	_args.Req = Req
	var _result SendCodeResult
	if err = p.c.Call(ctx, "sendCode", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
