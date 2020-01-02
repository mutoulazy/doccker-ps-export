FROM golang:1.13 AS builder

WORKDIR $GOPATH/doccker-ps-export

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o /usr/local/bin/docker_ps_exporter ./cmd/main.go

FROM alpine:3

COPY --from=builder /usr/local/bin/docker_ps_exporter /usr/local/bin/docker_ps_exporter

EXPOSE 9491

CMD [ "/usr/local/bin/docker_ps_exporter" ]