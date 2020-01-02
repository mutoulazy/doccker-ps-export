# doccker-ps-export
提供docker ps的命令结果指标给prometheus,方便监控集群外的容器情况

## docker中使用
docker run --name docker-export --volume "/var/run/docker.sock":"/var/run/docker.sock" -p 9491:9491 {images}
