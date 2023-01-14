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
│   ├── .air.toml                   # air config file for server side live rendering when in dev mode
│   └── Taskfile.yaml               # task description to build this app
│
├── cmd
│   └── linkerpod                   # the linkerpod CLI command required to run the SPA server
│       └── main.go          
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


## MIT Licence
