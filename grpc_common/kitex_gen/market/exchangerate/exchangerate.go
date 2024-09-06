// Code generated by Kitex v0.10.3. DO NOT EDIT.

package exchangerate

import (
	"context"
	"errors"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	streaming "github.com/cloudwego/kitex/pkg/streaming"
	proto "google.golang.org/protobuf/proto"
	market "grpc_common/kitex_gen/market"
)

var errInvalidMessageType = errors.New("invalid message type for service method handler")

var serviceMethods = map[string]kitex.MethodInfo{
	"usdRate": kitex.NewMethodInfo(
		usdRateHandler,
		newUsdRateArgs,
		newUsdRateResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingUnary),
	),
}

var (
	exchangeRateServiceInfo                = NewServiceInfo()
	exchangeRateServiceInfoForClient       = NewServiceInfoForClient()
	exchangeRateServiceInfoForStreamClient = NewServiceInfoForStreamClient()
)

// for server
func serviceInfo() *kitex.ServiceInfo {
	return exchangeRateServiceInfo
}

// for stream client
func serviceInfoForStreamClient() *kitex.ServiceInfo {
	return exchangeRateServiceInfoForStreamClient
}

// for client
func serviceInfoForClient() *kitex.ServiceInfo {
	return exchangeRateServiceInfoForClient
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
	serviceName := "ExchangeRate"
	handlerType := (*market.ExchangeRate)(nil)
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
		"PackageName": "market",
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

func usdRateHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(market.RateReq)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(market.ExchangeRate).UsdRate(ctx, req)
		if err != nil {
			return err
		}
		return st.SendMsg(resp)
	case *UsdRateArgs:
		success, err := handler.(market.ExchangeRate).UsdRate(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*UsdRateResult)
		realResult.Success = success
		return nil
	default:
		return errInvalidMessageType
	}
}
func newUsdRateArgs() interface{} {
	return &UsdRateArgs{}
}

func newUsdRateResult() interface{} {
	return &UsdRateResult{}
}

type UsdRateArgs struct {
	Req *market.RateReq
}

func (p *UsdRateArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(market.RateReq)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *UsdRateArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *UsdRateArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *UsdRateArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, nil
	}
	return proto.Marshal(p.Req)
}

func (p *UsdRateArgs) Unmarshal(in []byte) error {
	msg := new(market.RateReq)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var UsdRateArgs_Req_DEFAULT *market.RateReq

func (p *UsdRateArgs) GetReq() *market.RateReq {
	if !p.IsSetReq() {
		return UsdRateArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *UsdRateArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *UsdRateArgs) GetFirstArgument() interface{} {
	return p.Req
}

type UsdRateResult struct {
	Success *market.RateRes
}

var UsdRateResult_Success_DEFAULT *market.RateRes

func (p *UsdRateResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(market.RateRes)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *UsdRateResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *UsdRateResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *UsdRateResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, nil
	}
	return proto.Marshal(p.Success)
}

func (p *UsdRateResult) Unmarshal(in []byte) error {
	msg := new(market.RateRes)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *UsdRateResult) GetSuccess() *market.RateRes {
	if !p.IsSetSuccess() {
		return UsdRateResult_Success_DEFAULT
	}
	return p.Success
}

func (p *UsdRateResult) SetSuccess(x interface{}) {
	p.Success = x.(*market.RateRes)
}

func (p *UsdRateResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *UsdRateResult) GetResult() interface{} {
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

func (p *kClient) UsdRate(ctx context.Context, Req *market.RateReq) (r *market.RateRes, err error) {
	var _args UsdRateArgs
	_args.Req = Req
	var _result UsdRateResult
	if err = p.c.Call(ctx, "usdRate", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
