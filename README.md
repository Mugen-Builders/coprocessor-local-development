<br>
<p align="center">
    <img src="" align="center" width="20%">
</p>
<br>
<div align="center">
    <i></i>
</div>
<div align="center">
<b></b>
</div>
<br>
<p align="center">
	<img src="https://img.shields.io/github/license/tribeshq/tribes?style=default&logo=opensourceinitiative&logoColor=white&color=959CD0" alt="license">
	<img src="https://img.shields.io/github/last-commit/tribeshq/tribes?style=default&logo=git&logoColor=white&color=D1DCCB" alt="last-commit">
</p>

##  Table of Contents

- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Running](#running)

##  Getting Started

###  Prerequisites
1. [Install Docker Desktop for your operating system](https://www.docker.com/products/docker-desktop/).

    To install Docker RISC-V support without using Docker Desktop, run the following command:
    
   ```shell
   ❯ docker run --privileged --rm tonistiigi/binfmt --install all
   ```

2. [Download and install the latest version of Node.js](https://nodejs.org/en/download).

3. Cartesi CLI is an easy-to-use tool to build and deploy your dApps. To install it, run:

   ```shell
   ❯ npm i -g @cartesi/cli
   ```

> [!IMPORTANT]
>  To run the system in development mode, it is required to install:
>
> 1. [Download and Install the latest version of Golang.](https://go.dev/doc/install)
> 2. Install development node:
>
>   ```shell
>   ❯ npm i -g nonodo
>   ```


> [!TIP]
> Before start the application, export `COPROCESSOR_CALLER_MOCK_ADDRESS` as a enviroment variable:
>
>   ```shell
>   ❯ export COPROCESSOR_CALLER_MOCK_ADDRESS=<contract-address>
>   ```

###  Running

1. Build and generate excutable from source:

   ```sh
   ❯ go build -o solver ./cmd 
   ```

2. Run application:

   ```sh
   ❯ ./solver
   ```

3. Run validator node:

   ```sh
   ❯ cartesi run
   ```