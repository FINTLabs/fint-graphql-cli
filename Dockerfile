FROM golang:alpine AS builder
WORKDIR /go/src/app/vendor/github.com/FINTLabs/fint-graphql-cli
ARG VERSION=0.0.0
COPY . .
RUN go install -v -ldflags "-X main.Version=${VERSION}"
RUN /go/bin/fint-graphql-cli --version

FROM alpine
RUN apk --update add ca-certificates
COPY --from=builder /go/bin/fint-graphql-cli /usr/bin/fint-graphql-cli
WORKDIR /src
VOLUME [ "/src" ]
ENTRYPOINT [ "/usr/bin/fint-graphql-cli" ]
