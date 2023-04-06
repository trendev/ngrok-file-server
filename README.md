# ngrok-file-server

[![Build & Save artifacts](https://github.com/trendev/ngrok-file-server/actions/workflows/build.yml/badge.svg)](https://github.com/trendev/ngrok-file-server/actions/workflows/build.yml)

[![codecov](https://codecov.io/gh/trendev/ngrok-file-server/branch/main/graph/badge.svg?token=YIWQFBITBF)](https://codecov.io/gh/trendev/ngrok-file-server)

`ngrok-file-server` is a tiny secured file server. 

You can **quickly share any content** running a simple docker container protected.

It is **fast**, **secured** and so **lite**.

## Requirements

What do you really need ?

### ngrok account
https is provided by ngrok go implementation, so, you need a ngrok token

### docker
If you're familiar with golang, you can build your own server but we recommand to use our docker image

`docker run -it --rm -e NGROK_AUTHTOKEN="YOUR_TOKEN" -v $(pwd):/shared ghcr.io/trendev/ngrok-file-server`
