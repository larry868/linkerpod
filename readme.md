# Linkerpod

LinkerPod: Unleash Your Digital Identity

Linkerpod aims to act as a personal dashboard. It can be used like a link-in-bio, or like a brand portal, or a home page for a home lab.

In the digital age, your online presence is your identity. LinkerPod is your key to unlocking and enhancing that identity. It serves as your personal hub, a central command center where you can curate, manage, and showcase every aspect of your digital life. Whether you're a developer looking to flaunt your projects, a content creator sharing your expertise, or simply someone wanting a sleek online presence, LinkerPod empowers you to take control of your digital identity.

With LinkerPod, you're not just creating a web page; you're crafting your digital identity. It's time to unleash the full potential of your online presence with LinkerPod: Unleash Your Digital Identity.

## Promises

1. **Centralized Control:** LinkerPod offers you a unified platform to manage your online presence effortlessly. From links to social profiles, public keys, personal information, and projects, it's all in one place, under your control.
1. **Simplified Link Management:** Easily add, organize, and categorize your links to various online profiles, projects, and resources. No more scattered links or outdated information.
1. **Branded Domains:** Establish a professional and memorable online identity by connecting your custom domain. With LinkerPod, your digital presence will stand out.
1. **Effortless Sharing:** Share collections of your links with a unique and simple URL. Whether it's a link-in-bio for your social media, a portfolio for your projects, or a curated list of your favorite resources, sharing is seamless.
1. **Privacy Empowerment:** Your digital identity should be as private or public as you want it to be. LinkerPod provides robust privacy controls, so you decide who can see your links and content.
1. **Developer Integration:** If you're a developer, fork LinkerPod on GitHub and make it your own.

## Key features

- Build you own page
- Group and organize many web links to a single page
- Deployed over you domain or sub-domain.
- Integration with APIs.

Linkerpod is an experimental webapp written in go with the [icecake framework](icecake.dev).

Sources of inspiration: 
- https://dashy.to/
- https://linktr.ee/

## Demo

- [linkerpod homepage](https://linkerpod.net) is a linkerpod made with linkerpod !
- [larry868.github.io](https://larry868.github.io/linkerpod) is a simple static page built with linkerpod.

## Usage : make your own linkerpod

### Install on GitHub Pages

Fork this repo and activate [GitHub Pages](https://pages.github.com/):

1. fork linkerpod 
1. copy the default setup file `/docs/linkerpod_default.yaml` to `/docs/linkerpod.yaml`
1. customize the `/docs/linkerpod.yaml` file with your own links
1. if you want to display website favicons rather than icons you've selected:
    - download favicons and set them in a cache with the following command:
    1. install linkerpod with `go install github.com/larry868/linkerpod`
    1. go to docs in your forked repo: `cd {your_lnkerpod_repo}/docs`
    1. run `linkerpod -loadfavicons`
1. activate GitHub Pages on your repo, and specify `deploy from a branch` in `master` and `/docs`

### Setting up `/linkerpod.yaml` config file

All you need is links grouped in minipods.

```yaml
minipods:
    {minipod ID}:
        Links:
            {link ID}:
                link: {link url}
                [icon: {link icon}]
```

## To do

- [x] run saas transpiler with dev task
- [ ] header section of links
- [ ] save layout setup in the localstorage
- [ ] implement seo features
- [ ] keep comments in yaml setup file when running -loadfavicons
- [ ] make it a pwa with disconnected mode
- [ ] handle additional information on cards, not only link
- [ ] encrypt private information

## Roadmap

In a futur version linkerpod could implement some Web3 technologies such as :
- web3 authentication 
- decentralized storage
- work with avatar
- display NFT
- use IA to get the icons for minipods

## Tech

- Go 1.21
- CSS responsive framework, without any JS code: [Bulma](https://bulma.io/)
- [icecake framework](icecake.dev)
    - web assembly [see wasm doc](https://developer.mozilla.org/fr/docs/WebAssembly)
    - fullstack in go (no JS) [see this post to use wasm in GO](https://tutorialedge.net/golang/writing-frontend-web-framework-webassembly-go/)

### About Web Assembly with go

Some documentation available here https://tinygo.org/docs/guides/webassembly/ and here https://github.com/golang/go/wiki/WebAssembly

Go provides a specific js file called `wasm_exec.js` that need to be served by your webpapp. This file mustbe part of the static assets to be served by the server. To get the latest version you can _extract_ it from you go installation: `cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ./web/static`

## Project directory structure

The Directory structure follows [go ecosystem recommendation](https://github.com/golang-standards/project-layout).

```bash
linkerpod
├── build                           # build scripts
│   └── Taskfile.yaml               # building task configuration, ic. autobuild the front
│
├── cmd
│   └── linkerpod                 
│       └── linkerpod.go            # source for the linkerpod CLI
│
├── examples                        # examples of linkerpod yaml setup files
│
├── pkg                             # common sources for the command line and the wasm code
│
├── web                             # source codes and assets required by the front
│   ├── bulma-0.9.4                 # bulma saas files*
│   ├── saas
│   │       └── [*.*]               # any saas files
│   │
│   ├── static                      # 
│   │   ├── assets
│   │   │   ├── wasm_exec.js        # this file is mandatory and is provided by the go compiler
│   │   │   ├── wasm_loader.js      # this file is required to load the wasm code
│   │   │   └── [*.*]               # any img, js and other assets
│   │   ├── linkerpod.html          # The single and unique linkerpod html file
│   │   └── linkerpod.yaml          # The default linkerpod setup file
│   │
│   └── wasm
│       └── webapp.go               # the front app entry point
│           └── [*.*]               # front app wasm components
│
├── website OR docs                 # the self sufficient dir to serve the app in production, built with prod tasks (see Taskfile.yaml)
│   ├── *.*

```

_\* The `./web/bulma-0.9.4` directory is not sync with git. It need to be downloaded on the [bulma site](https://bulma.io/documentation/customize/with-sass-cli/)_


## Development

If you want to work on the layout, to debug, or to add features, consider the following

### setting dev env

We use [taskfile.dev](https://taskfile.dev) as task runner. See [Taskfile Installation](https://taskfile.dev/installation/) doc to install it.

Run the `devinit` task from the root path to setup the `tmp` directory for the first use:

```bash
task -t ./build/Taskfile.yaml devinit
```

We use the official saas transpiler written in node so we need to set it up. Alternativelly you can install the liveSassCompile VS Code extension but it will be up to you to run it, it won't be run automatically.

```bash
# install nvm
wget -qO- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.5/install.sh | bash
source ~/.bashrc

# Verify Installation
command -v nvm
nvm -v

## install node
nvm install node

## install sass
npm install -g sass less

```

Then you need to [install bulma](https://bulma.io/) locally.
Download the Bulma zip within your `/web/` directory like `/web/bulma-0.9.4/bulma/`

We're also using liverserver to test the static version of linkerpod page.

### running dev env

Run the `dev` task with the `--watch flag` to live reload your `linkedpod.html` page then launch liveServer on port 5511.

```bash
task -t ./build/Taskfile.yaml dev --watch
```

This task runs: 
- dev_sass : rebuilds `./tmp/linkerpos.css` based on `./web/sass/**/*`
- dev_static : moves `./web/static/**/*` to `./tmp`
- dev_wasm : rebuilds `./tmp/webapp.wasm` from `./web/wasm` including any `.pkg/**/*` changes

## Rebuild

:warning: Rebuilding the `/docs` directory clears all its previous content. Backup your linkerpod.yaml setup file before to run it.

```bash
task -t ./build/Taskfile.yaml rebuild
```


## Licence

[LICENCE](LICENCE)
 