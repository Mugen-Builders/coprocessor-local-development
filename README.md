<br>
<p align="center">
    <img src="https://github.com/Mugen-Builders/.github/assets/153661799/7ed08d4c-89f4-4bde-a635-0b332affbd5d" align="center" width="20%">
</p>
<br>
<div align="center">
<b>This application aims to be a local option to improve the development environment for applications built for the Cartesi Coprocessor</b>
</div>
<br>
<p align="center">
	<img src="https://img.shields.io/github/license/tribeshq/tribes?style=default&logo=opensourceinitiative&logoColor=white&color=00F6FF" alt="license">
	<img src="https://img.shields.io/github/last-commit/tribeshq/tribes?style=default&logo=git&logoColor=white&color=CCFDFF" alt="last-commit">
</p>

##  Table of Contents

- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Running](#running)

##  Getting Started

###  Prerequisites
1. [Download and Install the latest version of Golang.](https://go.dev/doc/install)

2. Install development node:
	```shell
	❯ npm i -g nonodo
	```

> [!TIP]
> Before start the application, export `COPROCESSOR_CALLER_MOCK_ADDRESS` as a enviroment variable:
>
>   ```shell
>   ❯ export COPROCESSOR_CALLER_MOCK_ADDRESS=<contract-address>
>   ```

###  Running

1. Start nonodo in a separate terminal:

   ```sh
   ❯ nonodo
   ```

2. Build and generate excutable from source:

   ```sh
   ❯ go build -o solver ./cmd 
   ```
   

3. Run application:

   ```sh
   ❯ ./solver
   ```
