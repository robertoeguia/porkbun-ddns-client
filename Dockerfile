FROM golang:1.21-alpine AS builder

COPY . /app

WORKDIR /app

RUN wget "https://github.com/upx/upx/releases/download/v4.2.2/upx-4.2.2-amd64_linux.tar.xz"; \
    tar -xvf "upx-4.2.2-amd64_linux.tar.xz"; 

RUN go mod download; \
    go build -C cmd/ -o ../ddns -ldflags="-s -w"; \
    ./upx-4.2.2-amd64_linux/upx ddns;

FROM alpine

COPY --from=builder /app/ddns /app/
COPY --from=builder /app/config/ /app/config

WORKDIR /app

ENTRYPOINT [ "/app/ddns" ]