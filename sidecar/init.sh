#!/bin/bash
# init コンテナでトラフィックの向きを変更するスクリプト。iptbles を叩くので特権が必要

# eth0 で TCP のトラフィックのうち、7000番ポート(payments のコンテナ)宛の通信を8000番(sidecar-proxy のコンテナ)に転送する設定。
# payments への通信は、sidecar-proxy でハンドリングを行う
#
# 旧:
#   [service] ---> payments:7000
# 新:
#   [service] ---> sidecar-proxy:8080
iptables -t nat -A PREROUTING -p tcp -i eth0 --dport 7000 -j REDIRECT --to-port 8080
