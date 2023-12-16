package global

import "github.com/micro/go-micro/v2"

var (
	GlobalRPCService micro.Service
)

func GetGlobalRPCService() micro.Service {
	return GlobalRPCService
}
