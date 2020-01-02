package collector

import (
	"context"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/prometheus/client_golang/prometheus"
)

type DockerContainers struct {
	*client.Client
}

var _ prometheus.Collector = (*DockerContainers)(nil)

// 指标描述
var (
	containerUptime = prometheus.NewDesc(
		"docker_ps_container_up",
		"Whether docker container is up, as reported by docker ps command.",
		[]string{"container_name", "container_id", "state", "image"}, nil,
	)
)

// 向管道填充数据
func (c DockerContainers) Describe(ch chan<- *prometheus.Desc) {
	ch <- containerUptime
}

// 从docker采集数据
func (c DockerContainers) Collect(ch chan<- prometheus.Metric) {
	// docker ps -a
	containers, err := c.Client.ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		log.Printf("Error while fetching container list: %s", err)
		return
	}

	for _, container := range containers {
		for _, name := range container.Names {
			up := isContainerUp(container)
			ch <- prometheus.MustNewConstMetric(
				containerUptime,
				prometheus.GaugeValue,
				boolToGaugeValue(up),
				name,
				container.ID,
				container.State,
				container.Image,
			)
		}
	}
}

// 判断docker 容器是否在运行状态
func isContainerUp(container types.Container) bool {
	return strings.EqualFold(container.State, "running")
}

func boolToGaugeValue(val bool) float64 {
	if val {
		return 1
	}
	return 0
}
