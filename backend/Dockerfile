# ビルドステージ: Goの公式イメージを使用 (バージョン1.21, Alpine Linuxベース)
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 依存関係のダウンロードに必要なgitをインストール
RUN apk add --no-cache git

# ソースコードをコンテナ内にコピー
COPY . .

# go.sum が存在しない場合があるため、生成と依存関係の整理を行う
ENV GOPROXY=https://proxy.golang.org,direct
RUN go mod tidy
RUN go mod download

# アプリケーションのビルド (出力ファイル名: main)
RUN go build -o main src/main.go

# 実行ステージ: 軽量なAlpine Linuxを使用
FROM alpine:latest
WORKDIR /root/

# ビルドステージからコンパイル済みのバイナリのみをコピー
COPY --from=builder /app/main .

# ポート8080を公開
EXPOSE 8080

# アプリケーションの実行
CMD ["./main"]
