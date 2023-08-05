# HPA 

## Step 1. Metrics Server をインストールする

(Kind では動かないと思うので EKS を使いますのでスキップします)

```shell
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```

```shell
$ kubectl top pod
No resources found in default namespace.
```

## Step 2. 重たいアプリケーションをデプロイする

```shell
$ kubectl apply -f app.yaml

$ kubectl get pods
NAME                         READY   STATUS    RESTARTS   AGE
heavy-app-86d84cf7dd-2fw75   1/1     Running   0          2m39s

# 別ターミナルで実行
$ ./request.sh

# 少し待って、top の結果を確認
$ kubectl top pod
NAME                         CPU(cores)   MEMORY(bytes)
heavy-app-86d84cf7dd-2fw75   97m          13Mi
```

## Step 3. HPA を設定する

```shell
$ kubectl apply -f hpa.yaml

# 別ターミナルで実行
$ ./request.sh

# 勝手にスケールしている
$ kubectl get pods
NAME                         READY   STATUS    RESTARTS   AGE
heavy-app-86d84cf7dd-2fw75   1/1     Running   0          4m29s
heavy-app-86d84cf7dd-c6bcs   1/1     Running   0          82s
heavy-app-86d84cf7dd-wt89f   1/1     Running   0          37s

$ kubectl get hpa
NAME        REFERENCE              TARGETS   MINPODS   MAXPODS   REPLICAS   AGE
heavy-app   Deployment/heavy-app   72%/80%   1         10        3          115s
```

