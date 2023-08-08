# 自作サイドカープロキシ

# TODO 1.

1. `init.sh` を修正して、payments のコンテナ宛の通信を sidecar-proxy のコンテナに転送するように設定してください
1. 適当な名前で `docker build` して自分のレジストリにアップロードしてください
1. `manifests/payments-with-sidecar.yaml` の init containers のイメージを2のイメージに変えてデプロイしてください
1. `kubectl -b back port-forward svc/backend 8000:80` で back コンテナと通信できるようにして、`curl http://localhost:8000/payment` を実行してください。
1. `kubectl -n payments logs ${payments の Pod 名} -c sidecar-proxy` で sidecar-proxy コンテナにアクセスログがあれば OK です

# TODO 2.

1. `main.go` を変更して、リクエストに `X-Internal-Token` ヘッダがなければエラーを返すような処理を追加してください
1. `go test ./...` でテストが通ったらデプロイしてください
