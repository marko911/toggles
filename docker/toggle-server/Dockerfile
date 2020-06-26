FROM golang:1.13.8-alpine3.11

# FROM golang:1.11.2
RUN apk add git

RUN apk add --no-cache bash
ARG env
ENV ENV $env

RUN echo "Installing global utilities..."
# Install go-bindata for inlining graphql schema
RUN go get -u -v github.com/jteeuwen/go-bindata/...
# Install fswatch to rebuild on file change
RUN go get -u -v github.com/codeskyblue/fswatch
# Install mockgen to create mock interfaces for tests
RUN go get github.com/golang/mock/gomock
RUN go get github.com/golang/mock/mockgen

WORKDIR /toggle/server
COPY . /scripts

# Development script installs dependencies and development specific libraries.
CMD ["/scripts/boot.sh"]
