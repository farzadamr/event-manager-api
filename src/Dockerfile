FROM golang:1.23.5-alpine AS builder

WORKDIR /app

COPY go.* ./
ENV CGO_ENABLED=0
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download

COPY . .

RUN go build -v -o server ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/config/config-docker.yml ./config/config-docker.yml
COPY --from=builder /app/storage ./storage
COPY --from=builder /app/template ./template

ENV APP_ENV=docker
ENV PORT=5005

EXPOSE 5005

CMD ["/app/server"]