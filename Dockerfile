FROM golang:1.20-alpine AS builder

LABEL stage=gobuilder

ENV GO111MODULE=on
ENV CGO_ENABLED 1
ENV GOPROXY https://goproxy.cn,direct

RUN echo -e http://mirrors.aliyun.com/alpine/v3.18/main/ > /etc/apk/repositories

ENV TZ=Asia/Shanghai
RUN apk add --no-cache tzdata && ln -sf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apk add --no-cache \
    # Important: required for go-sqlite3
    gcc \
    g++ \
    # Required for Alpine
    musl-dev

WORKDIR /build

COPY . .
RUN go mod tidy
RUN go build -ldflags="-s -w" -o /app/server main.go


FROM alpine

WORKDIR /app
# 需要先本地编译，手动 GOOS=linux GOARCH=amd64 go build -o grade
# COPY config/config.json /app/config/config.json
COPY --from=builder /usr/lib/libstdc++.so.6 /usr/lib/libstdc++.so.6
COPY --from=builder /usr/lib/libgcc_s.so.1 /usr/lib/libgcc_s.so.1
COPY --from=builder /go/pkg/mod/github.com/yanyiwu/gojieba@v1.3.0/dict/ /go/pkg/mod/github.com/yanyiwu/gojieba@v1.3.0/dict/
COPY --from=builder /app/server /app/server

# RUN echo -e http://mirrors.aliyun.com/alpine/v3.18/main/ > /etc/apk/repositories
# RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
#     && echo "Asia/Shanghai" > /etc/timezone \
#     && apk del tzdata

EXPOSE "$PORT"

CMD ["./server", "index", "-c", "/app/config/config.json"]
