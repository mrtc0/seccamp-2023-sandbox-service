# eBPF

## bpftrace でトレース

### Step 1. manifest を編集

```shell
vim bpftrace.yml
```

### Step 2. bpftrace.yml をデプロイ

```
kubectl apply -f bpftrace.yml
```

### Step 3. bpftrace コンテナでシステムコールをトレース

ノードで実行されたコマンド(≒ 他の人が実行したコマンドなども)取得されるはず

```shell
# コンテナ名は各自置き換える
$ kubectl -n default exec -it bpftrace-a-1 bash
(container) # bpftrace -e 'tracepoint:syscalls:sys_enter_exec* { printf("%-5d ", pid); join(args.argv); }'
```

### Step 4. bpftrace コンテナで bash の readline 関数をトレース

```shell
bpftrace -e 'uretprobe:/bin/bash:readline { printf("%-6d %s\n", pid, str(retval)); }'
```

↑ を実行しながら

```shell
# コンテナ名は各自置き換える
$ kubectl -n default exec -it bpftrace-a-1 bash
(container) # ls # とかなんでも打ってみると...
```

### Step 5. Clean up

```
$ kubectl delete -f bpftrace.yml
```
