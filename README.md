# fint-graphql-cli



## Description
Generates `GraphQL` schemas.

## Usage

### Exclude

To exclude classes, relations, and attributes, use the --exclude flag. It performs a case-insensitive match.

Example:
```
--exclude Fravar --exclude OTUngdom
```

To exclude from the schema, use the --exclude-schema flag.

Example:
```
--exclude-schema OTUngdom
```

## Install

### Binaries

Precompiled binaries are available as [Docker images](https://cloud.docker.com/u/fint/repository/docker/fint/graphql-cli)

Mount the directory where you want the generated source code to be written as `/src`.

Linux / MacOS:
```bash
docker run -v $(pwd):/src ghcr.io/fintlabs/fint-graphql-cli:latest <ARGS>
```

Windows PowerShell:
```ps1
docker run -v ${pwd}:/src ghcr.io/fintlabs/fint-graphql-cli:latest <ARGS>
```

### Source

To install, use `go get`:

```bash
go get -d github.com/FINTLabs/fint-graphql-cli
go install github.com/FINTLabs/fint-graphql-cli
```

## Author

[FINTLabs](https://fintlabs.github.io)
