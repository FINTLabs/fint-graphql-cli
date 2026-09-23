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

Use Go 1.27.1 and the checked-in vendor directory:

```bash
./build.sh
./build/fint-graphql-cli --version
./build/fint-graphql-cli --tag v4.1.0 generate --exclude Fravar --exclude Fravarstype
```

Version 2.0.0 generates Spring for GraphQL controllers with explicit `@QueryMapping`,
`@SchemaMapping`, and `@Argument` names. Query annotations are emitted only for
fields in the generated root schema, including `--exclude-schema` handling.
Imports are sorted and deduplicated, with separate Java imports and resource
wildcards for five or more classes from the same package.

Both `Date` and `Long` scalars are declared. Empty object types are omitted.
List relationships keep order, limit concurrency to eight, and retain null
entries for failed or empty downstream responses. Services with an optional
Feide identifier normalize blank usernames to null, preserving the application
behavior for Elev and Skoleressurs.

Build and test a local container without publishing it:

```bash
docker build --build-arg VERSION=2.0.0 -t fint-graphql-cli:2.0.0 .
docker run --rm fint-graphql-cli:2.0.0 --version
```

The application-specific merged `PersonService` remains an override in
`fint-graphql/PersonService.txt`; that application's `generate.sh` installs it
after generation. No Python conversion is required for the CLI output.

## Author

[FINTLabs](https://fintlabs.github.io)
