#!/bin/bash
# init コンテナでトラフィックの向きを変更するスクリプト。iptbles を叩くので特権が必要

# eth0 で TCP のトラフィックのうち、7000番ポート(payments のコンテナ)宛の通信を8000番(sidecar-proxy のコンテナ)に転送する設定。
# payments への通信は、sidecar-proxy でハンドリングを行う
#
# 旧:
#   [service] ---> payments:7000
# 新:
#   [service] ---> sidecar-proxy:8080
#
# TODO: XXXX に適切なポート番号を設定してください。--dport で「どこ宛にきた通信を」--to-port で「どこに転送するか」を指定します。
iptables -t nat -A PREROUTING -p tcp -i eth0 --dport XXXX -j REDIRECT --to-port XXXX
