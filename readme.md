# Linkerpod

Linkerpod allows referencing a set of URL links (aka. social media profiles), personnal info like licence number or web3 public keys, and monitoring informations to a single page and sharable link.

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

### Install on GitHub Pages

Fork this repo and activate [GitHub Pages](https://pages.github.com/):

1. fork linkerpod 
1. copy the default setup file `/docs/linkerpod_default.yaml` to `/docs/linkerpod.yaml`
1. customize the `/docs/linkerpod.yaml` file with your own links
1. if you want to display website favicons rather than icons you've selected:
    - download favicons and set them in a cache with the following command:
    1. install linkerpod with `go install github.com/lolorenzo777/linkerpod`
    1. go to docs in your forked repo: `cd {your_lnkerpod_repo}/docs`
    1. run `linkerpod -loadfavicons`
1. activate GitHub Pages on your repo, and specify `deploy from a branch` in `master` and `/docs`

## Roadmap

- [ ] make it a pwa
- [ ] handle additional information on cards, not only link
- [ ] encrypt private information
- [ ] keep comments in yaml setup file when running -loadfavicons

In a futur version linkerpod could implement some Web3 technologies such as :
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
├── cmd
│   └── linkerpod                 
│       └── linkerpod.go            # source for the linkerpod CLI
│
├── examples                        # examples of linkerpod setup files
│
├── pkg                             # common sources for the command line and the wasm code
│
├── web                             # source codes and assets required by the front
│   ├── bulma-0.9.4                 # bulma saas files*
│   ├── saas
│   │       └── [*.*]               # any saas files
│   │
│   ├── static
│   │   ├── assets
│   │   │   ├── wasm_exec.js        # this file is mandatory and is provided by the go compiler
│   │   │   ├── wasm_loader.js      # this file is required to load the wasm code
│   │   │   └── [*.*]               # any img, js and other assets
│   │   ├── linkerpod.html          # The single and unique linkerpod html file
│   │   └── linkerpod.yaml          # The default linkerpod setup file
│   │
│   └── wasm
│       └── webapp.go               # the front app entry point, uses components
│
├── website OR docs                 # the self sufficient dir to serve the app in production, built with prod tasks (see Taskfile.yaml)
│   ├── *.*

```

_\* The `./web/bulma-0.9.4` directory is not sync with git. It need to be downloaded on the [bulma site](https://bulma.io/documentation/customize/with-sass-cli/)_


### About Web Assembly with go

Some documentation available here https://tinygo.org/docs/guides/webassembly/ and here https://github.com/golang/go/wiki/WebAssembly

Go provides a specific js file called `wasm_exec.js` that need to be served by your webpapp. This file mustbe part of the static assets to be served by the server. To get the latest version you can _extract_ it from you go installation: `cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ./web/static`

## Development

We use ``task`` as task runner. See [Taskfile Installation](https://taskfile.dev/installation/) doc to install it.

In development mode run the `dev` task from the root path with the `--watch flag`:

```bash
task -t ./build/Taskfile.yaml dev --watch
```

Run the `build2docs` task to compile the linkerdpod wasm file and to rebuild the `/docs` directory.

```bash
task -t ./build/Taskfile.yaml build2docs
```

:warning: This rebuilds the `/docs` directory and clears all its content. Backup your linkerpod setup file before to run it.

If you've tuned `/web/saas/linkerpod.scss` you need to rebuild `.css` files according to your saas configuration.

## Licence

[LICENCE](LICENCE)
 