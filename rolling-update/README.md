# Rolling Update

## Rolling Update の様子を観察

### Step 1. nginx の Deployment をデプロイ

```yaml
kubectl apply -f nginx.yaml
```

### Step 2. 別ターミナルで Pod を監視

```shell
watch -n 1 kubectl get pods
```

`watch` コマンドがなければ...

```shell
while true; do kubectl get pods; sleep 1; done
```

### Step 3. nginx のイメージを更新してデプロイ

```yaml
- name: nginx
  image: nginx:1.25 # nginx:1.24 から変更する
  imagePullPolicy: Always
  ...
```

```shell
kubectl apply -f nginx.yaml
```

### Step 4. Pod の入れ替わりを確認

一時的に replicas 数を超過する Pod が2個(20%)になっていることを確認

## 安全なアップデート

もし、設定ミスやバグで起動しない場合...

### Step 5. nginx を起動しないように設定して適用

```shell
containers:
- name: nginx
  image: nginx:1.25
  imagePullPolicy: Always
  command: ["nginx", "-no-exists-option"] # 追加。-no-exists-option は存在しないので nginx は起動に失敗する
  ports:
  - containerPort: 80
```

```shell
kubectl apply -f nginx.yaml
```

### Step 6. Pod の入れ替わりが止まっていることを確認しよう

```shell
kubectl get pods
```
