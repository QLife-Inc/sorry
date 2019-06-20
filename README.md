# General Purpose Sorry Server

This is a server that only returns HTTP 503 status codes (Service Unavailable).  
It is assumed to be used when downtime occurs in server maintenance.  

## Features

- Responses with `HTTP 503` for all requests (excluding `/assets/*`).
- Responses with `503.json` when the `Accept` header contains `/json`.
- Responses with `503.json` when the request path ends with `.json`.
- Returns the contents of `503.html` for all other requests.
- Supports `HTTPS` and multiple domains.
- Allows specification of the `Retry-After` header through an environment variable.

## Usage

### Customize response contents

Edit `503.html` and` 503.json` according to your use case.

#### Deploy assets 

If necessary, save the files used by `503.html` in the `assets` directory.

### Start server

The listen port can be specified by the environment variable `PORT`.  
If not specified, it will listen on port `80`.  

> Note that on Linux, only the root user can use port 80.

```sh
PORT=8080 ./maint-server
```

### Listen HTTPS

Use the following directory structure if you want to use HTTPS.
 

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

- Supports multiple domains.
- The certificate file must be a chained certificate combining intermediate certificates.
- Files are determined by extension.
    Save the certificate file with extension `.crt` and the private key file with extension` .key`.
  - The certificate for that domain will not be read if it dose not have both.
- The HTTPS listener will not start if there is no `ssl` directory or certificate file.

#### Specify listen port

The listen port can be specified by the environment variable `HTTPS_PORT`.
If not specified, it will listen on port `443`.

> Note that on Linux, only the root user can use port 443.

```sh
PORT=8080 HTTPS_PORT=8443 ./maint-server
```

### Specify Retry-After

You can include a `Retry-After` header in your response by specifying the date and time in the `RETRY_AFTER` environment variable in the form `yyyy-MM-dd hh:mm:ss+0000`.

```sh
RETRY_AFTER="2019-06-20 23:59:59+0900" ./maint-server
```

## Use systemd service unit

Edit the following files in the `examples` directory and deploy them on the server.

- [/etc/sysconfig/sorry](examples/etc/sysconfig/sorry)
- [/usr/lib/systemd/system/sorry.service](examples/usr/lib/systemd/system/sorry.service)

After deployment, start the service with the following command.

```sh
sudo systemctl daemon-reload
sudo systemctl enable sorry
sudo systemctl start sorry
```

## Use Upstart init config

Edit the following files in the `examples` directory and deploy them on the server.

- [/etc/sysconfig/sorry](examples/etc/sysconfig/sorry)
- [/etc/init/sorry.conf](examples/etc/init/sorry.conf)

After deployment, start the service with the following command.

```sh
sudo initctl reload-configuration
sudo initctl start sorry
```

## Build

```sh
make linux # <- linux-amd64
make cross # <- cross platform build
```
