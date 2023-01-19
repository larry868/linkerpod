# Linkerpod web3 app

Linkerpod allows to reference a set of URL links (aka. social media profiles), personnal info, and web3 public keys to a single sharable link.

Linkerpod aims to act as a personal dashboard portal or a link-in-bio or a brand portal in the web3 era.

- build your own page of links in seconds
- stop spending time to lookup your public keys or your IDs 

Linkerpod is an experimental DApp aiming to discover and implement some Web3 technologies such as :
- decentralized storage
- authentication with a wallet
- display NFT
- work with avatar
- encrypt data
- reward with tokens

but also to discover and implement latest best-of-the-bread web technologies such as :
- web assembly [see wasm doc](https://developer.mozilla.org/fr/docs/WebAssembly)
- fullstack in go (no JS) [see this post to use wasm in GO](https://tutorialedge.net/golang/writing-frontend-web-framework-webassembly-go/)
- PWA with installation, off-line mode, notifications [see PWA doc](https://developer.mozilla.org/fr/docs/Web/Progressive_web_apps)
- HTTP/2 and gRPC with protobuff rather than RESTFull API [see this post about protocols](https://getstream.io/blog/communication-protocols/)
- redis strorage

Sources of inspiration: 
- https://dashy.to/
- https://github.com/maxence-charriere/go-app


## Backlog

- [ ] "hello world" wasm served by a SPA server, with dev environment setup.


## Tech

- Go 1.19
- CSS responsive framework, without any JS code: [Bulma](https://bulma.io/)


## Project layout

The Directory structure follows [go ecosystem recommendation](https://github.com/golang-standards/project-layout).

```bash
linkerpod
│   ├── readme.md
│   └── .gitignore
│
├── configs                         # congigurations files, loaded at server startup
│   ├── dev.env                 
│   └── prod.env          
│
├── build                           # build scripts
│   └── Taskfile.yaml               # building task configuration, ic. autobuild the front
│
├── cmd
│   └── linkerpod                   # the linkerpod CLI command required to run the SPA server
│       └── linkerpod.go          
│
├── pkg
│   ├── spa                         # SPA main package
│   │   ├── webserver.go                   
│   │   └── middleware.go                   
│   ├── api                         # API package, code to serve API requests
│   │   └── [*.go]                   
│   ├── sdk                         # SDK package, for any client willing to call APIs
│   │   └── [*.go]                   
│
├── web                             # source codes and assets required by the front, even server side templates
│   ├── component
│   │   ├── {componentname}.go      # every component in single go file
│   │   └── {componentname}.html    # html template of the component, optional, can also be directly described in the go file
│   ├── static
│   │   ├── wasm_exec.js            # this file is mandatory and is provided by the go compiler
│   │   └── [*.*]                   # any img, js 
│   └── wasm
│       └── main.go                 # the front app entry point, uses components
│
├── website                         # the self sufficient dir to serve the app in production, built with prod tasks (see Taskfile.yaml)
│   ├── *.*
│
├── tmp                             # built with dev tasks (see Taskfile.yaml)
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

## Testing

Useful read about [go data race detector](https://go.dev/doc/articles/race_detector#How_To_Use)

To be able to test wasm code on the browser, you need to install [wasmbrowsertest](https://github.com/agnivade/wasmbrowsertest):

```bash
$ go install github.com/agnivade/wasmbrowsertest@latest
$ mv $(go env GOBIN)/wasmbrowsertest $(go env GOBIN)/go_js_wasm_exec
```

Run the `unit_test` task to run both testing pkg and wasm:

```bash
$ task -t ./build/Taskfile.yaml unit_test
```

## Licence

[LICENCE](LICENCE)
 