# 汎用メンテナンスサーバー

HTTP 503 ステータス (Service Unavailable) を返すだけのサーバーです。  
サーバーのメンテナンスでダウンタイムが発生するときに利用する想定です。  

## Features

- すべてのリクエスト (`/assets/`を除く) に対して `HTTP 503` のレスポンスを返します。
- `Accept` ヘッダを見て、 `/json` が含まれていたら `503.json` のレスポンスを返します。
- リクエストのパスが `.json` で終わる場合も `503.json` のレスポンスを返します。
- それ以外は `503.html` の内容を返します。
- HTTPS 対応しています。複数ドメインに対応しています。

## Usage

### Edit your contents

`503.html` と `503.json` をユースケースに合わせて編集してください。

#### Deploy assets 

必要に応じて `503.html` で利用するファイルを `assets` ディレクトリに格納してください。

### Start server

Listen ポートは環境変数の `PORT` で指定できます。  
指定しなかった場合は `80` で Listen しますが、 Linux などでは root 以外 80 番ポートが利用できないのでご注意ください。

```sh
PORT=8080 ./maint-server
```

### Listen HTTPS

HTTPS を利用する場合は、以下のようなディレクトリ構成にして実行してください。 

```
.
├── assets/
├── 503.html
├── 503.json
├── maint-server
└── ssl
    ├── your-domain-1
    │   ├── your-domain-1-fullchain.crt
    │   └── your-domain-1-private.key
    ├── your-domain-2
    │   ├── your-domain-2-fullchain.crt
    │   └── your-domain-2-private.key
    └── ...
```

- 複数のドメインに対応しています。  
- 証明書は、必ず中間証明書などを結合した Chain 証明書にしてください。
- 拡張子でファイルを判別します。証明書を `.crt`, 秘密鍵を `.key` として格納してください。
  - 両方が揃っていない場合、そのドメインの証明書は読み込まれません。
- `ssl` ディレクトリに証明書がない場合は HTTPS のリスナーは起動しません。

#### Specified listen port

HTTPS のポートは `HTTPS_PORT` で指定してください。指定しなかった場合は 443 ポートが利用されますが、Linux (ry

```sh
PORT=8080 HTTPS_PORT=8443 ./maint-server
```

## Use systemd service unit

systemd のサービスとして登録する場合は以下のようにしてください。

```toml:/usr/lib/systemd/system/maint-server.service
[Unit]
Description=メンテナンスサーバー
After=network-online.target

[Service]
ExecStart=/path/to/maint-server/maint-server
ExecStop=/bin/kill -INT ${MAINPID}
Restart=always
WorkingDirectory=/path/to/maint-server

[Install]
WantedBy=multi-user.target
```

## Build

```sh
make linux # <- linux-amd64
make cross # <- cross platform build
```
