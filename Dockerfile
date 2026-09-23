FROM golang:1.27.1 AS builder
ENV CGO_ENABLED=0
WORKDIR /src
ARG VERSION=2.0.0
COPY . .
RUN go test -mod=vendor ./...
RUN go install -mod=vendor -trimpath -ldflags "-X main.Version=${VERSION}"
RUN /go/bin/fint-graphql-cli --version

FROM gcr.io/distroless/static
VOLUME [ "/src" ]
WORKDIR /src
COPY --from=builder /go/bin/fint-graphql-cli /usr/bin/fint-graphql-cli
ENTRYPOINT [ "/usr/bin/fint-graphql-cli" ]
