// Code generated by Kitex v0.10.3. DO NOT EDIT.
package exchangerate

import (
	server "github.com/cloudwego/kitex/server"
	market "grpc_common/kitex_gen/market"
)

// NewServer creates a server.Server with the given handler and options.
func NewServer(handler market.ExchangeRate, opts ...server.Option) server.Server {
	var options []server.Option

	options = append(options, opts...)

	svr := server.NewServer(options...)
	if err := svr.RegisterService(serviceInfo(), handler); err != nil {
		panic(err)
	}
	return svr
}

func RegisterService(svr server.Server, handler market.ExchangeRate, opts ...server.RegisterOption) error {
	return svr.RegisterService(serviceInfo(), handler, opts...)
}
