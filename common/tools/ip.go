package tools

import (
	"github.com/cloudwego/hertz/pkg/app"
	"net"
)

func GetRemoteClientIp(r *app.RequestContext) string {
	remoteIp := r.RemoteAddr().String()

	if ip := string(r.GetHeader("X-Real-IP")); ip != "" {
		remoteIp = ip
	} else if ip = string(r.GetHeader("X-Forwarded-For")); ip != "" {
		remoteIp = ip
	} else {
		remoteIp, _, _ = net.SplitHostPort(remoteIp)
	}

	//本地ip
	if remoteIp == "::1" {
		remoteIp = "127.0.0.1"
	}

	return remoteIp
}
