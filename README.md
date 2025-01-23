<br>
<p align="center">
    <img src="https://github.com/Mugen-Builders/.github/assets/153661799/7ed08d4c-89f4-4bde-a635-0b332affbd5d" align="center" width="20%">
</p>
<br>
<div align="center">
    <i>A tool for developing Cartesi Coprocessor applications</i>
</div>
<div align="center">
<b>This aims to be a development environment for Cartesi Coprocessor applications</b>
</div>
<br>
<p align="center">
	<img src="https://img.shields.io/github/license/Mugen-Builders/cartesi-coprocessor-nonodox?style=default&logo=opensourceinitiative&logoColor=white&color=00ADD8" alt="license">
	<img src="https://img.shields.io/github/last-commit/Mugen-Builders/cartesi-coprocessor-nonodox?style=default&logo=git&logoColor=white&color=000000" alt="last-commit">
</p>

## Table of Contents

- [Overview](#overview)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Running](#running)

## Overview

<div align="justify">
This is an iterative tool designed to accelerate the "debugging" and "development" process of applications using the Cartesi Coprocessor infrastructure, providing a faster path to the production environment.
</div>

## Getting Started

### Prerequisites

1. [Foundry](https://book.getfoundry.sh/getting-started/installation)
2. [Golang](https://go.dev/doc/install)
3. [Nonodo](https://github.com/Calindra/nonodo?tab=readme-ov-file#installation)

### Running

> [!WARNING]
> Before running the command below, ensure you have created a `.toml` file and set the **environment variables** correctly. Below is the structure of the content that should be included in the file:
>
> ```toml
> [anvil]
> http_url = "http://127.0.0.1:8545"
> ws_url = "ws://127.0.0.1:8545"
> private_key = ""
>
> [coprocessor]
> machine_hash = ""
> adapter_contract_address = ""
> ```

1. Install the package:

```sh
go install github.com/Mugen-Builders/cartesi-coprocessor-nonodox/cmd/nonodox@latest
```

2. Run the application:

```sh
nonodox --config <filename>.toml
```
