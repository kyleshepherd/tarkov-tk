# Stage 1 - Binary Build
# BUILD_X args should be passed at build time as docker build args
FROM golang:1.18.5-alpine3.16 AS builder
ARG BIN_VERSION
ARG GIT_COMMIT
ARG GIT_SHA
ARG GIT_TAG
ARG GIT_DIRTY
ARG GOPROXY
ARG GONOSUMDB
ENV GOPROXY=$GOPROXY
ENV BIN_OUTDIR=./
ENV BIN_NAME=tarkovtk
RUN apk update && apk add build-base git libressl-dev
WORKDIR /usr/src/tarkovtk
# install dependencies in separate docker layer
COPY go.mod .
COPY go.sum .
RUN go mod download
# copy application source and build
COPY ./ .
RUN make static

# Stage 2 - Final Image
# The application should be statically linked
FROM alpine:latest
RUN apk update \
  && apk add --no-cache ca-certificates \
  && rm -rf /var/cache/apk/* \
  && addgroup tarkovtk \
  && adduser -D -H -G tarkovtk tarkovtk
COPY --from=builder /usr/src/tarkovtk/tarkovtk /usr/bin/tarkovtk
VOLUME ["/etc/tarkovtk"]
ENTRYPOINT ["tarkovtk"]
EXPOSE 5000
USER tarkovtk
