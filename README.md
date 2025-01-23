<br>
<p align="center">
    <img src="https://github.com/Mugen-Builders/.github/assets/153661799/7ed08d4c-89f4-4bde-a635-0b332affbd5d" align="center" width="20%">
</p>
<br>
<div align="center">
<b>This application aims to be a local option to improve the development environment <br> for applications built for the Cartesi Coprocessor</b>
</div>
<br>
<p align="center">
	<img src="https://img.shields.io/github/license/henriquemarlon/coprocessor-local-development?style=default&logo=opensourceinitiative&logoColor=white&color=00F6FF" alt="license">
	<img src="https://img.shields.io/github/last-commit/henriquemarlon/coprocessor-local-development?style=default&logo=git&logoColor=white&color=000000" alt="last-commit">
</p>

## Table of Contents

- [Overview](#overview)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Running](#running)


## Overview

<div align="justify">
cartesi
</div>

##  Getting Started

###  Prerequisites

1. 
2. 
3. 

###  Running

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

2. Running:

```sh
nonodox --config <filename>.toml
```
