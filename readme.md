# Linkerpod

Linkerpod allows to reference a set of URL links (aka. social media profiles), personnal info like web3 public keys, and monitoring informations to a single page and sharable link.

Linkerpod aims to act as a personal dashboard. It can be used like a link-in-bio, or like a brand portal, or a home page for a home lab.

## Promises

- build your own page of links in seconds
- stop spending time to look for your links
- get you public keys or your IDs instantly
- share collections of links
- monitor your favorite website or homelab
- display information from API queries

Linkerpod is an experimental webapp written in go with the [icecake framework](icecake.dev).

Sources of inspiration: 
- https://dashy.to/
- https://linktr.ee/

## Demo

[linkerpod demo](https://lolorenzo777.github.io/linkerpod/)

## Usage

The simplest way is to fork this repo and to activate [GitHub Pages](https://pages.github.com/):

1. fork linkerpod 
1. update the `/docs/linkerpod.yaml` file with your own links
1. activate GitHub Pages on your repo, and specify `deploy from a branch` and choose `master` and `/docs`

## Roadmap

In a futur version, linkerpod could implement some Web3 technologies such as :
- decentralized storage
- authentication with a wallet
- display NFT
- work with avatar
- encrypt data
- reward with tokens

## Tech

- Go 1.20
- CSS responsive framework, without any JS code: [Bulma](https://bulma.io/)
- [icecake framework](icecake.dev)
    - web assembly [see wasm doc](https://developer.mozilla.org/fr/docs/WebAssembly)
    - fullstack in go (no JS) [see this post to use wasm in GO](https://tutorialedge.net/golang/writing-frontend-web-framework-webassembly-go/)

## Project directory structure

The Directory structure follows [go ecosystem recommendation](https://github.com/golang-standards/project-layout).

```bash
linkerpod
│   ├── readme.md
│   └── .gitignore
│
├── build                           # build scripts
│   └── Taskfile.yaml               # building task configuration, ic. autobuild the front
│
├── web                             # source codes and assets required by the front
│   ├── bulma-0.9.4                 # bulma saas files
│   ├── saas
│   │   ├── ick.scss                # required by icecake
│   │   └── linkerpod.scss          # customized saas
│   │
│   ├── static
│   │   ├── assets
│   │   │   ├── wasm_exec.js        # this file is mandatory and is provided by the go compiler
│   │   │   ├── wasm_loader.js      # this file is required to load the wasm code
│   │   │   └── [*.*]               # any img, js and other assets
│   │   └── index.html              # The single and unique linkerpod html file
│   │
│   └── wasm
│       └── webapp.go               # the front app entry point, uses components
│
├── website                         # the self sufficient dir to serve the app in production, built with prod tasks (see Taskfile.yaml)
│   ├── *.*

```

The `./web/bulma-0.9.4` directory is not sync with git. It need to be downloaded on the [bulma site](https://bulma.io/documentation/customize/with-sass-cli/)


### About Web Assembly with go

Some documentation available here https://tinygo.org/docs/guides/webassembly/ and here https://github.com/golang/go/wiki/WebAssembly

Go provides a specific js file called `wasm_exec.js` that need to be served by your webpapp. This file mustbe part of the static assets to be served by the server. To get the latest version you can _extract_ it from you go installation: `cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ./web/static`

## Development

We use ``task`` as task runner. See [Taskfile Installation](https://taskfile.dev/installation/) doc to install it.

In development mode run the `dev` task from the root path with the `--watch flag`:

```bash
$ task -t ./build/Taskfile.yaml dev --watch
```

## Licence

[LICENCE](LICENCE)
 