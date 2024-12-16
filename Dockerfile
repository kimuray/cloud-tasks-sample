# ベースとなるビルド環境イメージ
FROM golang:1.23 AS builder

# 作業ディレクトリ設定
WORKDIR /app

# Goの依存ファイルをキャッシュするため、先にgo.mod, go.sumをコピー
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# ソースコードをコピー
COPY . .

# アプリケーションをビルド (static build)
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server ./cmd/worker.go

# 実行用の軽量なベースイメージを使用
# Distrolessイメージを利用する場合
FROM gcr.io/distroless/base-debian11

# コンテナ内で実行するアプリケーションバイナリをコピー
COPY --from=builder /app/server /server

# ポート指定（Cloud Runは$PORTを使うため明示指定しなくても可、必要ならEXPOSE）
EXPOSE 8080

# 実行コマンド
CMD ["/server"]
