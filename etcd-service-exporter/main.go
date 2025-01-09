package main

import (
	"context"
	"etcd-service-exporter/client"
	"etcd-service-exporter/logger"
)

func main() {
	err := logger.SetupGlobalLogger()
	if err != nil {
		panic("初始化全局日志失败")
	}
	client.AutoFetchServices(context.Background())
	select {} // 阻塞主 goroutine
}
