# fint-graphql-cli



## Description
Generates `GraphQL` schemas.

## Usage



## Install

### Binaries

Precompiled binaries are available as [Docker images](https://cloud.docker.com/u/fint/repository/docker/fint/graphql-cli)

Mount the directory where you want the generated source code to be written as `/src`.

Linux / MacOS:
```bash
docker run -v $(pwd):/src fint/graphql-cli:latest <ARGS>
```

Windows PowerShell:
```ps1
docker run -v ${pwd}:/src fint/graphql-cli:latest <ARGS>
```

### Source

To install, use `go get`:

```bash
go get -d github.com/FINTLabs/fint-graphql-cli
go install github.com/FINTLabs/fint-graphql-cli
```

## Author

[FINTLabs](https://fintlabs.github.io)
