# seccamp-2023 B3 講義資料

セキュリティ・キャンプ 全国大会 2023 で使用する資料です。

# Setup

## 必要なもの

- [kubectl](https://kubernetes.io/ja/docs/tasks/tools/install-kubectl/)
- [AWS CLI](https://docs.aws.amazon.com/ja_jp/cli/latest/userguide/getting-started-install.html)
- 手元で動作する Kubernetes クラスタ

## clone

```shell
git clone https://github.com/mrtc0/seccamp-2023-sandbox-service
```

## マイクロサービスのデプロイ

**手元のクラスタ** にデプロイしてください。EKS へはデプロイしないでください。

```shell
$ kubectx local-cluster # クラスタ名は各自環境に合わせて指定
$ ./setup.sh
```

## EKS への接続

共有した AWS クレデンシャルを設定しておいてください。

```shell
$ cat ~/.aws/credentials
[default]
aws_access_key_id = ASIA…
aws_secret_access_key = …
aws_region = ap-northeast-1
```

次のコマンドで EKS への接続設定をして、`kubectl get pods` が成功すれば OK です。

```shell
$ aws eks update-kubeconfig --name sandbox-seccamp --region ap-northeast-1 --alias eks-seccamp
Added new context eks-seccamp to /root/.kube/config

$ kubectl get pods
No resources found in default namespace.

$ kubectx # クラスタ切り替え
```

## AWS CLI コンテナを使えるようにする

攻撃で取得したクレデンシャルを使って AWS にアクセスすることがあります。手元の端末にセットアップされている AWS クレデンシャルを意図せず使ってしまっては意味がないので、まっさらな環境を使うために AWS CLI コンテナを使えるようにしておくと便利です。

```shell
docker run --rm -it --entrypoint bash amazon/aws-cli
```

攻撃で取得したクレデンシャルをワーキングディレクトリに保存してボリュームマウントすると、都度クレデンシャルをセットせずに済むので便利かもしれません。

```shell
$ cat credentials
aws_access_key_id = ASIA…
aws_secret_access_key = …
aws_session_token = ...
aws_region = ap-northeast-1

$ docker run --rm -it -v $PWD:/root/.aws eks-utils:latest bash
```

