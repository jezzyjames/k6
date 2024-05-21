FROM golang:alpine AS builder

RUN go install go.k6.io/xk6/cmd/xk6@latest
RUN xk6 build --with github.com/grafana/xk6-output-influxdb


FROM alpine
COPY --from=builder /go/k6 .