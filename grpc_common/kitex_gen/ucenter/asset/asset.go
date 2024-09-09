// Code generated by Kitex v0.10.3. DO NOT EDIT.

package asset

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
	"findWalletBySymbol": kitex.NewMethodInfo(
		findWalletBySymbolHandler,
		newFindWalletBySymbolArgs,
		newFindWalletBySymbolResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingUnary),
	),
	"findWallet": kitex.NewMethodInfo(
		findWalletHandler,
		newFindWalletArgs,
		newFindWalletResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingUnary),
	),
	"resetWalletAddress": kitex.NewMethodInfo(
		resetWalletAddressHandler,
		newResetWalletAddressArgs,
		newResetWalletAddressResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingUnary),
	),
	"findTransaction": kitex.NewMethodInfo(
		findTransactionHandler,
		newFindTransactionArgs,
		newFindTransactionResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingUnary),
	),
	"GetAddress": kitex.NewMethodInfo(
		getAddressHandler,
		newGetAddressArgs,
		newGetAddressResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingUnary),
	),
}

var (
	assetServiceInfo                = NewServiceInfo()
	assetServiceInfoForClient       = NewServiceInfoForClient()
	assetServiceInfoForStreamClient = NewServiceInfoForStreamClient()
)

// for server
func serviceInfo() *kitex.ServiceInfo {
	return assetServiceInfo
}

// for stream client
func serviceInfoForStreamClient() *kitex.ServiceInfo {
	return assetServiceInfoForStreamClient
}

// for client
func serviceInfoForClient() *kitex.ServiceInfo {
	return assetServiceInfoForClient
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
	serviceName := "Asset"
	handlerType := (*ucenter.Asset)(nil)
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

func findWalletBySymbolHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(ucenter.AssetReq)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(ucenter.Asset).FindWalletBySymbol(ctx, req)
		if err != nil {
			return err
		}
		return st.SendMsg(resp)
	case *FindWalletBySymbolArgs:
		success, err := handler.(ucenter.Asset).FindWalletBySymbol(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*FindWalletBySymbolResult)
		realResult.Success = success
		return nil
	default:
		return errInvalidMessageType
	}
}
func newFindWalletBySymbolArgs() interface{} {
	return &FindWalletBySymbolArgs{}
}

func newFindWalletBySymbolResult() interface{} {
	return &FindWalletBySymbolResult{}
}

type FindWalletBySymbolArgs struct {
	Req *ucenter.AssetReq
}

func (p *FindWalletBySymbolArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(ucenter.AssetReq)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *FindWalletBySymbolArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *FindWalletBySymbolArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *FindWalletBySymbolArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, nil
	}
	return proto.Marshal(p.Req)
}

func (p *FindWalletBySymbolArgs) Unmarshal(in []byte) error {
	msg := new(ucenter.AssetReq)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var FindWalletBySymbolArgs_Req_DEFAULT *ucenter.AssetReq

func (p *FindWalletBySymbolArgs) GetReq() *ucenter.AssetReq {
	if !p.IsSetReq() {
		return FindWalletBySymbolArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *FindWalletBySymbolArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *FindWalletBySymbolArgs) GetFirstArgument() interface{} {
	return p.Req
}

type FindWalletBySymbolResult struct {
	Success *ucenter.MemberWallet
}

var FindWalletBySymbolResult_Success_DEFAULT *ucenter.MemberWallet

func (p *FindWalletBySymbolResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(ucenter.MemberWallet)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *FindWalletBySymbolResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *FindWalletBySymbolResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *FindWalletBySymbolResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, nil
	}
	return proto.Marshal(p.Success)
}

func (p *FindWalletBySymbolResult) Unmarshal(in []byte) error {
	msg := new(ucenter.MemberWallet)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *FindWalletBySymbolResult) GetSuccess() *ucenter.MemberWallet {
	if !p.IsSetSuccess() {
		return FindWalletBySymbolResult_Success_DEFAULT
	}
	return p.Success
}

func (p *FindWalletBySymbolResult) SetSuccess(x interface{}) {
	p.Success = x.(*ucenter.MemberWallet)
}

func (p *FindWalletBySymbolResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *FindWalletBySymbolResult) GetResult() interface{} {
	return p.Success
}

func findWalletHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(ucenter.AssetReq)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(ucenter.Asset).FindWallet(ctx, req)
		if err != nil {
			return err
		}
		return st.SendMsg(resp)
	case *FindWalletArgs:
		success, err := handler.(ucenter.Asset).FindWallet(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*FindWalletResult)
		realResult.Success = success
		return nil
	default:
		return errInvalidMessageType
	}
}
func newFindWalletArgs() interface{} {
	return &FindWalletArgs{}
}

func newFindWalletResult() interface{} {
	return &FindWalletResult{}
}

type FindWalletArgs struct {
	Req *ucenter.AssetReq
}

func (p *FindWalletArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(ucenter.AssetReq)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *FindWalletArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *FindWalletArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *FindWalletArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, nil
	}
	return proto.Marshal(p.Req)
}

func (p *FindWalletArgs) Unmarshal(in []byte) error {
	msg := new(ucenter.AssetReq)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var FindWalletArgs_Req_DEFAULT *ucenter.AssetReq

func (p *FindWalletArgs) GetReq() *ucenter.AssetReq {
	if !p.IsSetReq() {
		return FindWalletArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *FindWalletArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *FindWalletArgs) GetFirstArgument() interface{} {
	return p.Req
}

type FindWalletResult struct {
	Success *ucenter.MemberWalletList
}

var FindWalletResult_Success_DEFAULT *ucenter.MemberWalletList

func (p *FindWalletResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(ucenter.MemberWalletList)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *FindWalletResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *FindWalletResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *FindWalletResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, nil
	}
	return proto.Marshal(p.Success)
}

func (p *FindWalletResult) Unmarshal(in []byte) error {
	msg := new(ucenter.MemberWalletList)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *FindWalletResult) GetSuccess() *ucenter.MemberWalletList {
	if !p.IsSetSuccess() {
		return FindWalletResult_Success_DEFAULT
	}
	return p.Success
}

func (p *FindWalletResult) SetSuccess(x interface{}) {
	p.Success = x.(*ucenter.MemberWalletList)
}

func (p *FindWalletResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *FindWalletResult) GetResult() interface{} {
	return p.Success
}

func resetWalletAddressHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(ucenter.AssetReq)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(ucenter.Asset).ResetWalletAddress(ctx, req)
		if err != nil {
			return err
		}
		return st.SendMsg(resp)
	case *ResetWalletAddressArgs:
		success, err := handler.(ucenter.Asset).ResetWalletAddress(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*ResetWalletAddressResult)
		realResult.Success = success
		return nil
	default:
		return errInvalidMessageType
	}
}
func newResetWalletAddressArgs() interface{} {
	return &ResetWalletAddressArgs{}
}

func newResetWalletAddressResult() interface{} {
	return &ResetWalletAddressResult{}
}

type ResetWalletAddressArgs struct {
	Req *ucenter.AssetReq
}

func (p *ResetWalletAddressArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(ucenter.AssetReq)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *ResetWalletAddressArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *ResetWalletAddressArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *ResetWalletAddressArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, nil
	}
	return proto.Marshal(p.Req)
}

func (p *ResetWalletAddressArgs) Unmarshal(in []byte) error {
	msg := new(ucenter.AssetReq)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var ResetWalletAddressArgs_Req_DEFAULT *ucenter.AssetReq

func (p *ResetWalletAddressArgs) GetReq() *ucenter.AssetReq {
	if !p.IsSetReq() {
		return ResetWalletAddressArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *ResetWalletAddressArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *ResetWalletAddressArgs) GetFirstArgument() interface{} {
	return p.Req
}

type ResetWalletAddressResult struct {
	Success *ucenter.AssetResp
}

var ResetWalletAddressResult_Success_DEFAULT *ucenter.AssetResp

func (p *ResetWalletAddressResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(ucenter.AssetResp)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *ResetWalletAddressResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *ResetWalletAddressResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *ResetWalletAddressResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, nil
	}
	return proto.Marshal(p.Success)
}

func (p *ResetWalletAddressResult) Unmarshal(in []byte) error {
	msg := new(ucenter.AssetResp)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *ResetWalletAddressResult) GetSuccess() *ucenter.AssetResp {
	if !p.IsSetSuccess() {
		return ResetWalletAddressResult_Success_DEFAULT
	}
	return p.Success
}

func (p *ResetWalletAddressResult) SetSuccess(x interface{}) {
	p.Success = x.(*ucenter.AssetResp)
}

func (p *ResetWalletAddressResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *ResetWalletAddressResult) GetResult() interface{} {
	return p.Success
}

func findTransactionHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(ucenter.AssetReq)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(ucenter.Asset).FindTransaction(ctx, req)
		if err != nil {
			return err
		}
		return st.SendMsg(resp)
	case *FindTransactionArgs:
		success, err := handler.(ucenter.Asset).FindTransaction(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*FindTransactionResult)
		realResult.Success = success
		return nil
	default:
		return errInvalidMessageType
	}
}
func newFindTransactionArgs() interface{} {
	return &FindTransactionArgs{}
}

func newFindTransactionResult() interface{} {
	return &FindTransactionResult{}
}

type FindTransactionArgs struct {
	Req *ucenter.AssetReq
}

func (p *FindTransactionArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(ucenter.AssetReq)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *FindTransactionArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *FindTransactionArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *FindTransactionArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, nil
	}
	return proto.Marshal(p.Req)
}

func (p *FindTransactionArgs) Unmarshal(in []byte) error {
	msg := new(ucenter.AssetReq)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var FindTransactionArgs_Req_DEFAULT *ucenter.AssetReq

func (p *FindTransactionArgs) GetReq() *ucenter.AssetReq {
	if !p.IsSetReq() {
		return FindTransactionArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *FindTransactionArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *FindTransactionArgs) GetFirstArgument() interface{} {
	return p.Req
}

type FindTransactionResult struct {
	Success *ucenter.MemberTransactionList
}

var FindTransactionResult_Success_DEFAULT *ucenter.MemberTransactionList

func (p *FindTransactionResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(ucenter.MemberTransactionList)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *FindTransactionResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *FindTransactionResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *FindTransactionResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, nil
	}
	return proto.Marshal(p.Success)
}

func (p *FindTransactionResult) Unmarshal(in []byte) error {
	msg := new(ucenter.MemberTransactionList)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *FindTransactionResult) GetSuccess() *ucenter.MemberTransactionList {
	if !p.IsSetSuccess() {
		return FindTransactionResult_Success_DEFAULT
	}
	return p.Success
}

func (p *FindTransactionResult) SetSuccess(x interface{}) {
	p.Success = x.(*ucenter.MemberTransactionList)
}

func (p *FindTransactionResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *FindTransactionResult) GetResult() interface{} {
	return p.Success
}

func getAddressHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(ucenter.AssetReq)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(ucenter.Asset).GetAddress(ctx, req)
		if err != nil {
			return err
		}
		return st.SendMsg(resp)
	case *GetAddressArgs:
		success, err := handler.(ucenter.Asset).GetAddress(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*GetAddressResult)
		realResult.Success = success
		return nil
	default:
		return errInvalidMessageType
	}
}
func newGetAddressArgs() interface{} {
	return &GetAddressArgs{}
}

func newGetAddressResult() interface{} {
	return &GetAddressResult{}
}

type GetAddressArgs struct {
	Req *ucenter.AssetReq
}

func (p *GetAddressArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(ucenter.AssetReq)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *GetAddressArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *GetAddressArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *GetAddressArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, nil
	}
	return proto.Marshal(p.Req)
}

func (p *GetAddressArgs) Unmarshal(in []byte) error {
	msg := new(ucenter.AssetReq)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var GetAddressArgs_Req_DEFAULT *ucenter.AssetReq

func (p *GetAddressArgs) GetReq() *ucenter.AssetReq {
	if !p.IsSetReq() {
		return GetAddressArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *GetAddressArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *GetAddressArgs) GetFirstArgument() interface{} {
	return p.Req
}

type GetAddressResult struct {
	Success *ucenter.AddressList
}

var GetAddressResult_Success_DEFAULT *ucenter.AddressList

func (p *GetAddressResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(ucenter.AddressList)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *GetAddressResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *GetAddressResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *GetAddressResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, nil
	}
	return proto.Marshal(p.Success)
}

func (p *GetAddressResult) Unmarshal(in []byte) error {
	msg := new(ucenter.AddressList)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *GetAddressResult) GetSuccess() *ucenter.AddressList {
	if !p.IsSetSuccess() {
		return GetAddressResult_Success_DEFAULT
	}
	return p.Success
}

func (p *GetAddressResult) SetSuccess(x interface{}) {
	p.Success = x.(*ucenter.AddressList)
}

func (p *GetAddressResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *GetAddressResult) GetResult() interface{} {
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

func (p *kClient) FindWalletBySymbol(ctx context.Context, Req *ucenter.AssetReq) (r *ucenter.MemberWallet, err error) {
	var _args FindWalletBySymbolArgs
	_args.Req = Req
	var _result FindWalletBySymbolResult
	if err = p.c.Call(ctx, "findWalletBySymbol", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) FindWallet(ctx context.Context, Req *ucenter.AssetReq) (r *ucenter.MemberWalletList, err error) {
	var _args FindWalletArgs
	_args.Req = Req
	var _result FindWalletResult
	if err = p.c.Call(ctx, "findWallet", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) ResetWalletAddress(ctx context.Context, Req *ucenter.AssetReq) (r *ucenter.AssetResp, err error) {
	var _args ResetWalletAddressArgs
	_args.Req = Req
	var _result ResetWalletAddressResult
	if err = p.c.Call(ctx, "resetWalletAddress", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) FindTransaction(ctx context.Context, Req *ucenter.AssetReq) (r *ucenter.MemberTransactionList, err error) {
	var _args FindTransactionArgs
	_args.Req = Req
	var _result FindTransactionResult
	if err = p.c.Call(ctx, "findTransaction", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetAddress(ctx context.Context, Req *ucenter.AssetReq) (r *ucenter.AddressList, err error) {
	var _args GetAddressArgs
	_args.Req = Req
	var _result GetAddressResult
	if err = p.c.Call(ctx, "GetAddress", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
