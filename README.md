# Electrumx Proxy Go for Atomicals

Electrumx Proxy Go is an implementation of a proxy server for Electrumx written in Go. It is designed to allow clients to connect to Electrumx servers through a high-performance, scalable Go-based intermediary. By serving as a robust buffer between clients and Electrumx servers, Electrumx Proxy Go significantly increases the stability and reliability of RPC communications within the Atomical ecosystem. This enhancement is crucial for systems that require consistent and uninterrupted access to blockchain data, ensuring that applications remain responsive and resilient to fluctuations in demand or network conditions.


## Prerequisites

Before you begin, ensure you have Go version 1.21 or higher installed on your system. You can verify your Go
installation and version by running:

```bash
go version
```

This command should output a version number. If the version number is less than 1.21, you will need to update Go to a
newer version. Visit the [official Go download page](https://golang.org/dl/) for instructions on how to install or
update Go.

## Installation

To install Electrumx Proxy Go, clone the repository, and run the main application:

```bash
git clone https://github.com/NeutronProtocol/electrumx-proxy-go.git
cd electrumx-proxy-go
```

## Configuration

Before starting the proxy, you need to configure the `config.toml` file with the correct parameters:

- `ElectrumxServer` should be set to the address of your Electrumx server (e.g., `ws://127.0.0.1:50002`).
- `ServerAddress` should be set to the address and port you want the proxy to listen on (e.g., `0.0.0.0:8080`).

Open `config.toml` in a text editor and modify the settings accordingly:

```toml
ElectrumxServer = "ws://127.0.0.1:50002"
ServerAddress = "0.0.0.0:8080"
```

Save the file after making the necessary changes.

## Starting the Proxy

Once you have configured the settings, you can start the proxy by running:

```bash
go run main.go
```

Or

```bash
go build
./electrumx-proxy-go
```

## Usage

After starting the main application, the proxy will begin listening for incoming connections on the
configured `ServerAddress`. Point your Electrumx server's client connections to this address.

## Features

- High-performance Go-based proxy.
- Easy integration with existing Electrumx servers.
- Lightweight and efficient connection handling.

## License

This project is licensed under the MIT License - see the `LICENSE` file for details.

## Contact

Project
Link: [https://github.com/NeutronProtocol/electrumx-proxy-go](https://github.com/NeutronProtocol/electrumx-proxy-go)

Website: [https://www.atomicalneutron.com/](https://www.atomicalneutron.com/)
