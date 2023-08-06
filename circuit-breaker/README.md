# サーキットブレーカーで障害を緩和しよう

## シナリオ

`web` サービスは `api` サービスと通信をしています。`api` は v1 から v2 に徐々に切り替えるために、リクエストを両方のバージョンに流して運用しています。  
しかし、v2 でバグのあるバージョンをデプロイしてしまい、一部のリクエストがエラーを返すようになりました。このとき、ユーザーが障害を体験することが少なくなるように、自動で v2 へのリクエストを止めて v1 に流すようにしてみましょう。

## Step 1. web / api v1 を起動

```shell
docker compose up
```

`web` サービスの URL `http://localhost:9091/` にリクエストを投げると、アップストリームの情報が返ってきます。`api-v1` が返ってくれば OK です。

```shell
$ curl -s 'http://localhost:9091/' | jq '.upstream_calls."http://localhost:9092".name, .code'
"api-v1"
200
```

## Step 2. api v2 をデプロイ 

v1 と v2 を同時に運用します。`web` サービスの sidecar プロキシとして動いている Consul がリクエストを `api` サービスの v1 と v2 に振り分けてくれます。

```shell
docker compose -f docker-compose.app-v2.yml up
```

Step 1 と同様に `web` サービスにリクエストを送信すると、アップストリームへのリクエストは `api-v1` と `api-v2` に振り分けられていることが確認できます。


```shell
$ curl -s 'http://localhost:9091/' | jq '.upstream_calls."http://localhost:9092".name, code'
"api-v2"
200

$ curl -s 'http://localhost:9091/' | jq '.upstream_calls."http://localhost:9092".name, code'
"api-v1"
200
```

## Step 3. バグのある api v2 をデプロイ

バグのある(90%の確率でエラーを返す) api v2 をデプロイします。  
`docker compose` のオプション `--force-recreate` をつけて今動いている `api-v2` と入れ替えます。

```shell
docker compose -f docker-compose.app-v2-fail.yaml up --force-recreate
```

`request.sh` は1秒ごとに `web` サービスにリクエストを送信するスクリプトです。実行してしばらく様子をみてください。  
`api-v2` からのレスポンスエラーが閾値を超えると `api-v1` のみにリクエストが割り振られていることが確認できます。  
その後、30秒ほど待って再度リクエストを送ると `api-v2` にリクエストが割り振られるようになります。

```shell
$ ./request.sh
"api-v2"
500

...

"api-v1"
200
"api-v1"
200

(Ctrl + c で中断して30秒ほど待つ)

$ curl -s 'http://localhost:9091/' | jq '.upstream_calls."http://localhost:9092".name, code'
"api-v2"
500
```

