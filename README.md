# Connect-RPC を使ってみる

Connect-RPC を使ってストリーム通信を試してみる。

gRPC を使おうとしたが、プロキシの設定などで躓く。そこで Web との通信が簡単な Connect-RPC を使ってみることに下。

## 使い方

vsCode の dev container での開発を前提。
./ssl に証明書を置く。

1. dev container を立ち上げる
2. `make run` Go サーバが立ち上がる
3. `cd web && yarn dev` React が立ち上がる

## proto

proto ファイルは、`./api` 以下にある。コード生成のためのパッケージは、dev container を立ち上げた後の `./script/init-devcontainer.sh` でインストールされる。

コードの生成は、`make genbuf`

## 中身

クライアントから日付範囲とキーワードをリクエストする。

サーバは、受け取ったデータから国会議事録を取得し、形態素解析で名詞の共起リストをストリームで返す。

クライアントは全てを受け取った後、上位の共起ペアをネットワークグラフで表示する。
