# ngrok-file-server

[![Test, Build and Save](https://github.com/trendev/ngrok-file-server/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/trendev/ngrok-file-server/actions/workflows/build.yml)

[![codecov](https://codecov.io/gh/trendev/ngrok-file-server/branch/main/graph/badge.svg?token=YIWQFBITBF)](https://codecov.io/gh/trendev/ngrok-file-server)

[![Go Reference](https://pkg.go.dev/badge/github.com/trendev/ngrok-file-server.svg)](https://pkg.go.dev/github.com/trendev/ngrok-file-server)

`ngrok-file-server` is a **tiny secured file server** :bowtie:

You can **quickly share any content** running a simple docker image !!!

It's **fast** :rocket:, **secured** :lock: and pretty **lite** :mouse2:

## :point_right: Requirements

What do you really need ?

### ngrok token
**https tunneling** (port-forwarding) is provided by [ngork](https://ngrok.com/) go implementation. So, you **need a ngrok token** :key:

... if don't have a ngrok account, don't worry, **it's free** :grimacing:

### docker
If you're familiar with golang, you can build your own server but we recommand to **use our docker image** :whale:

You can download [docker desktop](https://www.docker.com/products/docker-desktop/) and use default settings.

## :arrow_forward: Run

Ok, so now, you may have a **ngrok token** and **run docker** on your machine.

### :one: Let's start the server.

Run the following docker command in order to share the content of your local path.

`docker run -it --rm -e NGROK_AUTHTOKEN="YOUR_TOKEN" -v $(pwd):/shared ghcr.io/trendev/ngrok-file-server`
> this example is based on Linux/MacOS usage. If you run `ngrok-file-server` on Windows, you may change local path command `$(pwd)` by something like `%cd%` or direcly set the path...

> do not forget to replace "YOUR_TOKEN" by your ngrok token...

It will pull our docker image and then, build and start a new container :thumbsup:

You may see an output like this :

`ngrok ingress url:  https://2927-90-22-212-26.eu.ngrok.io`

Copy the URL :memo:

### :two: Visit the file server from anywhere

Paste the previous URL in your favorite web browser.
> if you use a free ngrok account, please, just accept the disclaimer :kissing_cat:

And here we are :white_check_mark:

You can browse your content and share the URL with anyone and access to your file server from anywhere :sunglasses:

## :cop: OAuth2 Protection

You can also control access using a oauth2 provider (like Google, Facebook, Github, Linkedin, etc) and setting an authorized domain (for eg, only `trendev.fr` users).

> you can find the supported list [here](https://ngrok.com/docs/cloud-edge/modules/oauth/#oauth-providers-supported-by-ngrok)

### Enable oauth2 authentication
`docker run -it --rm -e NGROK_AUTHTOKEN="YOUR_TOKEN" -v $(pwd):/shared ghcr.io/trendev/ngrok-file-server --provider=google`

### Enable oauth2 authentication + domain control
`docker run -it --rm -e NGROK_AUTHTOKEN="YOUR_TOKEN" -v $(pwd):/shared ghcr.io/trendev/ngrok-file-server --provider=google --domain=trendev.fr`

## :hand: Something else ?

If you need more, you can open an issue and describe your requirement... and BTW, you can also star the repo :wink:


