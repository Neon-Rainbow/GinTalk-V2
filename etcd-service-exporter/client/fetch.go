package client

import (
	"context"
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"os"
	"time"
)

type PrometheusTarget struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

type Service struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (s *Service) toPrometheusTarget() *PrometheusTarget {
	return &PrometheusTarget{
		Targets: []string{fmt.Sprintf("%v:%v", s.Host, s.Port)},
		Labels: map[string]string{
			"job": s.Name,
			"id":  s.ID,
		},
	}
}

// fetchServicesFromEtcd 用于从 etcd 中获取服务信息
// 该函数会从 etcd 中获取所有的服务信息，然后将其转换为 PrometheusTarget
func fetchServicesFromEtcd(ctx context.Context, prefix string) ([]*PrometheusTarget, error) {
	cli := GetClient()
	resp, err := cli.Get(ctx, prefix, clientv3.WithPrefix())

	targets := make([]*PrometheusTarget, 0, 10)

	if err != nil {
		return nil, err
	}
	for _, kv := range resp.Kvs {
		// 解析服务信息
		s := &Service{}
		if err := json.Unmarshal(kv.Value, s); err != nil {
			zap.L().Error("解析服务信息失败", zap.Error(err))
			return nil, err
		}
		targets = append(targets, s.toPrometheusTarget())
	}
	return targets, nil
}

// generatePrometheusFileSd 生成Prometheus file_sd_configs使用的目标文件 (generate the target file for Prometheus file_sd_configs)
//
// 参数:
//   - targets: PrometheusTarget 列表
//   - outputFile: 输出文件
//
// 返回值:
//   - error: 错误信息
func generatePrometheusFileSd(targets []*PrometheusTarget, outputFile string) error {
	// 序列化为 JSON
	data, err := json.MarshalIndent(targets, "", "  ")
	if err != nil {
		return fmt.Errorf("无法 target 序列化为 json: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(outputFile, data, 0644); err != nil {
		return fmt.Errorf("无法写入目标文件: %w", err)
	}

	zap.L().Info("更新 Prometheus file_sd 文件成功", zap.String("file", outputFile))
	return nil
}

// AutoFetchServices 定时从ETCD中拉取服务列表并生成Prometheus file_sd文件
//
// 参数:
//   - ctx: 上下文
//
// 使用示例:
//
//	etcd.AutoFetchServices(context.Background())
func AutoFetchServices(ctx context.Context) {
	prefix := "/service/"
	outputFile := "../docker-compose/prometheus/file_sd_targets.json"
	go func() {
		for range time.Tick(10 * time.Second) {
			targets, err := fetchServicesFromEtcd(ctx, prefix)
			if err != nil {
				zap.L().Error("从ETCD中拉取服务列表失败", zap.Error(err))
				continue
			}
			if err := generatePrometheusFileSd(targets, outputFile); err != nil {
				zap.L().Error("更新Prometheus file_sd文件失败", zap.Error(err))
			}
		}
	}()
}
