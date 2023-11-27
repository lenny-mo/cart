package utils

import (
	"net/http"
	"strconv"

	"github.com/ngaut/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func PrometheusBoot(port int) {
	// 在 "/metrics" 路径上注册一个处理器，用于 Prometheus 的数据抓取
	http.Handle("/metrics", promhttp.Handler())

	// 使用 goroutine 异步启动 HTTP 服务器
	go func() {
		// 构造监听地址和端口，启动 HTTP 服务
		err := http.ListenAndServe("0.0.0.0:"+strconv.Itoa(port), nil)
		// 如果启动失败，记录致命错误并退出
		if err != nil {
			log.Fatal("启动失败: ", err)
		}
	}()

	// 记录日志信息，表明监控服务已启动
	log.Info("监控启动，端口为：" + strconv.Itoa(port))
}