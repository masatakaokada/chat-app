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
- ログイン・ログアウト
- ルーム作成
- 全体チャット(WebSocket)
- ルームチャット(WebSocket)

## インフラ

- 開発環境・本番環境ともに Docker を使用。
- Go の本番環境はビルドファイルを Alpine のイメージに置いて実行。
- Vue の本番環境はビルドファイルを Nginx サーバーに置いて実行。
- Heroku の Container Registry を使用している。

## Migration after deploy

```
cd app && goose -env production status
goose -env production up
```
