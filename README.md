# Connect-RPC を使ってみる

Connect-RPC を使ってストリーム通信を試してみる。

gRPC を使おうとしたが、プロキシの設定などで躓く。そこで Web との通信が簡単な Connect-RPC を使ってみることにした。

## 中身

![使用感](./docs/images/example.png)

クライアントから日付範囲とキーワードをリクエストする。

サーバは、受け取ったデータから国会議事録APIを叩き、リクエストされたキーワードに関連する共起ペアリストを返す。このとき作業の進捗状況もサーバから返す。

クライアントは進捗を表示し、共起ペアリストをネットワークグラフで描画する。

ネットワークグラフ上のノード（単語）をクリックすると、そのノードに関連するノードがサーバから取得される。

## 起動方法

``docker compose up --build``

## 開発環境

vsCode の dev container での開発を前提。

1. dev container を立ち上げる
2. `make run` Go サーバが立ち上がる
3. `cd web && yarn dev` React が立ち上がる

### proto

proto ファイルは、`./api` 以下にある。コード生成のためのパッケージは、dev container を立ち上げた後の `./script/init-devcontainer.sh` でインストールされる。

コードの生成は、`make genbuf`

## 問題点

TF-IDFでノードを全体の1割に削減。

ネットワーク中心性でノードの重要度を計算、表現。

~~nodeとedgeが場合によって多すぎる。~~

~~ノードクリック時のリクエストは、初期リクエストとは別ストリームのためオンメモリ＋fileに永続化して対応しているが、多すぎるとメモリーリークする気がする。~~

メモリや永続化をしていない。

~~あとは、単純にネットワークグラフが混雑して何も分からない。~~

~~なにかしらのアルゴリズムでより重要な共起ペアリストをフィルタするのが良いかもしれない。~~
