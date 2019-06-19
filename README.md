# Generic maintenance server

It is a server that only returns HTTP 503 status (Service Unavailable).  
It is assumed to be used when downtime occurs in server maintenance.  

## Features

- Responses `HTTP 503` for all request (without `/assets/*`).
- Responses `503.json` when contains `/json` in `Accept` header.
- Responses `503.json` when ends with `.json` for request path.
- Otherwise it responses the contents of `503.html`.
- Supported `HTTPS`, supported multi domain.

## Usage

### Edit your contents

Edit `503.html` and` 503.json` according to your use case.

#### Deploy assets 

If necessary, save the files used by `503.html` in the `assets` directory.

### Start server

Listen port can be specified by the environment variable `PORT`.  
If not specified, it will listen on `80` port.  

> Notice, on Linux etc, only root user can use 80 port.

```sh
PORT=8080 ./maint-server
```

### Listen HTTPS

Use the following directory structure, if you want to use HTTPS.
 

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

- Supported multi domain.
- Certificate file must be a chained certificate combining intermediate certificates.
- Determine the file by extension.
    Save the certificate file extension `.crt` and the private key file extension` .key`.
  - Certificate for that domain will not be read if not have both.
- HTTPS listener will not start if there is no `ssl` directory or certificate file.

#### Specify listen port

Listen port can be specified by the environment variable `HTTPS_PORT`.
If not specified, it will listen on `443` port.

> Notice, on Linux etc, only root user can use 443 port.

```sh
PORT=8080 HTTPS_PORT=8443 ./maint-server
```

## Systemd service unit example

```toml
# /usr/lib/systemd/system/maint-server.service
[Unit]
Description=Maintenance Server
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
