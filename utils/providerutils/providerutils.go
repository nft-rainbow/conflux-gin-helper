package providerutils

import (
	"context"

	providers "github.com/openweb3/go-rpc-provider/provider_wrapper"
	"github.com/sirupsen/logrus"
)

func GetCountStatisticMiddleware(network string) providers.CallContextMiddleware {
	return func(call providers.CallContextFunc) providers.CallContextFunc {
		rpcCalls := make(map[string]uint)
		return func(ctx context.Context, resultPtr interface{}, method string, args ...interface{}) error {
			rpcCalls[method]++
			logrus.WithField("rpc calls", rpcCalls).WithField("network", network).Info("rpc call")
			return call(ctx, resultPtr, method, args...)
		}
	}
}
