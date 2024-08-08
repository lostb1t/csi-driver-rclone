FROM golang:alpine as builder

RUN apk update && apk add --no-cache --progress git make
ARG GOPROXY
ADD . /build
RUN cd /build && go mod download -x
RUN cd /build && make

FROM alpine:latest

RUN apk update && apk add --no-cache --progress fuse3 rclone
COPY --from=builder /build/bin/* /bin

ENTRYPOINT ["/bin/csi-rclone"]
