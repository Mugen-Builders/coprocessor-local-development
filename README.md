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

##  Table of Contents

- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Running](#running)

##  Getting Started

###  Prerequisites
1. [Download and Install the latest version of Golang.](https://go.dev/doc/install)

###  Running

1. Build and generate excutable from source:

```sh
go build -o mugen ./cmd 
```
   
> [!WARNING]
> Replace the argument below with the address of the previously deployed `CoprocessorCallerMock.sol`.

2. Run application:

```sh
./mugen --coprocessor-caller-address <contract-address>
```
