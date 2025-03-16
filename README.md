# `abi`: Easy to use ABI encoding and decoding functionalities for EVM blockchains

[![Go Reference](https://pkg.go.dev/badge/github.com/omnes-tech/abi.svg)](https://pkg.go.dev/github.com/omnes-tech/abi)
[![Go Report Card](https://goreportcard.com/badge/github.com/omnes-tech/abi)](https://goreportcard.com/report/github.com/omnes-tech/abi)
[![Coverage Status](https://coveralls.io/repos/github/omnes-tech/abi/badge.svg?branch=main)](https://coveralls.io/github/omnes-tech/abi?branch=main)
[![Latest Release](https://img.shields.io/github/v/release/omnes-tech/abi)](https://github.com/omnes-tech/abi/releases/latest)
<!-- <img src="https://w3.cool/gopher.png" align="right" alt="W3 Gopher" width="158" height="224" -->

Encode to and decode from bytecode easily with `abi` package. Usage is similar to Solidity's `abi.encode`, `abi.encodePacked`, and `abi.decode` functions.

```shell
go get github.com/omnes-tech/abi
```

## At a Glance

Encode functions:
- `Encode`
- `EncodePacked`
- `EncodeSelector`
- `EncodeWithSignature`
- `EncodeWithSelector`

Decode functions:
- `Decode`
- `DecodePacked`
- `DecodeWithSignature`
- `DecodeWithSelector`
