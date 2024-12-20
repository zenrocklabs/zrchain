# zrChain


![Banner!](/docs/img/banner.png)

[![License: Source Available License](https://img.shields.io/github/license/zenrocklabs/zenrock.svg?style=flat-square)](https://github.com/zenrocklabs/zenrock/blob/main/LICENSE)
[![Version](https://img.shields.io/github/tag/Zenrock-Foundation/zenrock.svg?style=flat-square)](https://github.com/Zenrock-Foundation/zrchain/releases/latest)
![Go](https://img.shields.io/badge/go-1.23-blue.svg)
[![GitHub Super-Linter](https://img.shields.io/github/actions/workflow/status/Zenrock-Foundation/zrchain/lint.yml?style=flat-square&label=Lint)](https://github.com/marketplace/actions/super-linter)
[![Discord](https://badgen.net/badge/icon/discord?icon=discord&label)](https://discord.com/invite/zenrocklabs)
[![Twitter](https://badgen.net/badge/icon/twitter?icon=twitter&label)](https://twitter.com/OfficialZenrock)

# Overview 

Welcome to **Zenrock**, the bedrock security infrastructure that will support an omnichain future. More details to be rolled out in the coming weeks and months. 

# Table of Contents
- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)
- [Support](#support)

## Installation

### Requirements

- [go 1.23+](https://go.dev/doc/install)
- Make
- [Docker](https://docs.docker.com/get-docker/) (used to regenerate protobufs)

### Commands

Clone the repo and initialize the zenrockd binary with the following commands:

```bash
git clone git@github.com:Zenrock-Foundation/zrchain.git
cd zrchain
make install
```

Check that the zenrockd binaries have been successfully installed: 

```bash
zenrockd version
```

If the zenrockd command is not found an error message is returned, confirm that your GOPATH is correctly configured by running the following command:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

## Usage

To run your chain locally, you can either run 

```bash
zenrockd start
```

or alternatively using our custom script which starts the chain and comes with additional data populated in the genesis. 

Usage: open root at `zrchain` dir in 2 terminal tabs:

Tab 1 - `./init.sh --localnet 1`
Tab 2 - `./init.sh --localnet 2`

This will spin up 2 validators and a sidecar used by both

The script can also be used to run a single node as well as other modes using flags i.e:
- `--non-validator`
- `--no-vote-extensions`

## Contributing

We are appreciating all contributions to Zenrock and will closely review them. Find more information on how to contribute to the zrChain in the [contributing file](./CONTRIBUTING.md).

## Roadmap

Find latest information about our roadmap on [zenrocklabs.io](https://www.zenrocklabs.io/).

## License
Licensed under the Source Available License, Zenrock Foundation DAO. See LICENSE file for details.

## Support
If you want to get in contact with the Zenrock team you can either open an issue on github or join our [Discord](https://discord.com/invite/zenrocklabs) server.
