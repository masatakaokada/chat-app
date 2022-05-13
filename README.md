# Chat App

WebSocket を用いたチャットツールのアプリケーションです。

フロントエンドは Vue.js、バックエンドは Go 言語を用いています。

Go 言語及び Echo Framework を用いた API 実装、WebSocket 通信の学習のために作成しました。

## 使用技術

client side

- Node.js 17.9
- Vue.js Vue3

server side

- Go 1.8
- Echo Framework

共通

- Firebase(Firebase Authentication)
- WebSocket
- Docker/Docker Compose
- MySQL 8
- Heroku

## 仕様

- ユーザー登録
  - client
    - Firebase Authentication で email と password を登録。
  - server
    - ユーザー情報は firebase_uid と共に自前の DB で管理する。
- ログイン・ログアウト
  - client
    - Firebase Authentication で認証し、JWT を取得。
    - JWT はローカルストレージに保存し、通信時に Authorization ヘッダを用いて認証する。
  - server
    - JWT を検証。
    - JWT の有効期限の確認も行う
- 全体チャット(WebSocket)
  - client
    - サーバーにリクエストを送ってコネクションを確立。
    - メッセージを送信する。
    - ヘッダが使用できないため QueryParameter で JWT をやり取りする。
  - server
    - クライアントのリクエストを受信してコネクションを確立。
    - メッセージを受け取って goroutine を用いて接続しているユーザー全員にメッセージを broadcast する。

## インフラ

- 開発環境・本番環境ともに Docker を使用。
- Go の本番環境はビルドファイルを Alpine のイメージに置いて実行。
- Vue の本番環境はビルドファイルを Nginx サーバーに置いて実行。
- Heroku の Container Registry を使用している。
