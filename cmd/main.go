package main

import (
	"doccker-ps-export/internal/collector"
	"github.com/alecthomas/kingpin"
	"github.com/docker/docker/client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

var (
	listenAddress = kingpin.Flag("web.listen-address", "Address on which to expose metrics.").Default(":9491").String()
	metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
)

func boot() error {
	// 新建docker 客户端
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	coll := collector.DockerContainers{Client: dockerClient}
	if err := prometheus.Register(&coll); err != nil {
		return err
	}

	http.Handle(*metricsPath, promhttp.Handler())
	log.Println("docker export server start!")
	return http.ListenAndServe(*listenAddress, nil)
}

func main() {
	kingpin.Parse()
	kingpin.FatalIfError(boot(), "")
}
