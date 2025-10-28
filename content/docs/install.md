---
title: 'Install'
weight: 20
slug: install
description: "Instructions for installing FauxRPC using pre-built binaries, Docker, or from source."
icon: "download"
---

## Pre-built binaries
Binaries are built for several platforms for each release. See the latest ones on [the releases page](https://github.com/sudorandom/fauxrpc/releases/latest).

## Docker Image
You can also run FauxRPC using the docker image, available at [hub.docker.com/r/sudorandom/fauxrpc](https://hub.docker.com/r/sudorandom/fauxrpc).

```shell
docker run -it --rm -p 6660:6660 docker.io/sudorandom/fauxrpc:latest run --empty --addr=:6660
```

## Install via source
This is not recommended long-term but this is sufficient for playing around with FauxRPC.
```
go install github.com/sudorandom/fauxrpc/cmd/fauxrpc@latest
```

## Run with `go run`
This is also not recommended long-term but this is sufficient for playing around with FauxRPC.
```
go run github.com/sudorandom/fauxrpc/cmd/fauxrpc@latest
```
