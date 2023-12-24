# 使用 golang alpine 做為建置階段的基底 image
FROM golang:1.20-alpine AS builder

# 設置工作目錄
WORKDIR /workspace

# 安裝 gcc 和其他必要的依賴項目
# go alpine 版本不包含這些項目
RUN apk add --no-cache git make gcc musl-dev

# 複製你的 go 程式碼到 container 中
COPY . .

# 使用 go build 來編譯你的應用程式
# 使用 `-ldflags -s -w` 來降低 binary 大小
RUN go build -ldflags="-w -s" -o service.app

# 創建最終階段的基底 image
FROM alpine:latest

# 更換工作目錄到 /app
WORKDIR /app/

# 從建置階段的 image 複製 binary 到最終階段的 image 中
COPY --from=builder /workspace/service.app /app/service.app

# 開放你的應用程式使用的 port
EXPOSE 5001

# 設定權限跟預留 mountpoint
RUN mkdir /datastore && \
    chown -R 1001:0 /datastore /app/service.app && \
	chmod -R g+rwX /datastore

# 不存在的使用者
USER 1001

# 設定容器啟動時執行的指令
ENTRYPOINT ["/app/service.app"]