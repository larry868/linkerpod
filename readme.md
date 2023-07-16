# Linkerpod web3 app

Linkerpod allows to reference a set of URL links (aka. social media profiles), personnal info, and web3 public keys to a single sharable link.

Linkerpod aims to act as a personal dashboard. It can be used like a link-in-bio, or like a brand portal, or a home page for a home lab.

- build your own page of links in seconds
- stop spending time to look for your links
- get you public keys or your IDs instantly
- share your collection of inks

Linkerpod is an experimental webapp written in go with the [icecake framework](icecake.dev).

Sources of inspiration: 
- https://dashy.to/
- https://github.com/maxence-charriere/go-app

## Usage

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
├── pkg
│
├── web                             # source codes and assets required by the front
│   ├── static
│   │   ├── wasm_exec.js            # this file is mandatory and is provided by the go compiler
│   │   └── [*.*]                   # any img, js 
│   └── wasm
│       └── main.go                 # the front app entry point, uses components
│
├── website                         # the self sufficient dir to serve the app in production, built with prod tasks (see Taskfile.yaml)
│   ├── *.*

```

### About Web Assembly with go

Some documentation available here https://tinygo.org/docs/guides/webassembly/ and here https://github.com/golang/go/wiki/WebAssembly

Go provides a specific js file called `wasm_exec.js` that need to be served by your webpapp. This file mustbe part of the static assets to be served by the server. To get the latest version you can _extract_ it from you go installation: `cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ./web/static`

## Development

We use ``task`` as task runner. See [Taskfile Installation](https://taskfile.dev/installation/) doc to install it.

In development mode run the `dev_front` task from the root path with the `--watch flag`:

```bash
$ task -t ./build/Taskfile.yaml dev_front --watch
```

This task:
1. moves any changed files in ``./web/static/`` to ``./tmp/website/``
1. builds/rebuilds any frontend components and the .wasm file
1. builds/rebuilds the ``./tmp/website/spa.wasm`` file according to changes in the ``web/wasm/main.go``

Start the server either in debug mode with the `F5` in vscode, or by running the `dev_back` task:

```bash
$ task -t ./build/Taskfile.yaml dev_back
```

## Licence

[LICENCE](LICENCE)
 