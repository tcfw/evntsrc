FROM augustoroman/v8-lib:6.7.290 as libv8
FROM golang:latest as builder

#Install XZ
RUN apt-get update && apt-get install -y \
	xz-utils \
	&& rm -rf /var/lib/apt/lists/*

#Install UPX
COPY ./build/upx-3.94-amd64_linux.tar.xz /usr/local
RUN xz -d -c /usr/local/upx-3.94-amd64_linux.tar.xz | \
	tar -xOf - upx-3.94-amd64_linux/upx > /bin/upx && \
	chmod a+x /bin/upx

#Install glide for installing deps
RUN curl https://glide.sh/get | sh

WORKDIR /go/src/github.com/tcfw/evntsrc/

#Install deps and build
COPY glide.yaml glide.yaml
COPY glide.lock glide.lock
RUN glide install --strip-vendor

#Install V8
COPY --from=libv8 /v8/include ./vendor/github.com/augustoroman/v8/include
COPY --from=libv8 /v8/lib ./vendor/github.com/augustoroman/v8/libv8

#Build
COPY . .
RUN go build -ldflags="-w -s" -o adapter ./cmd/adapter
RUN strip --strip-unneeded adapter
# RUN upx adapter - UPX has unexpected results with compressing

# ---

FROM ubuntu:latest

WORKDIR /bin
COPY --from=builder /go/src/github.com/tcfw/evntsrc/adapter adapter

CMD [ "adapter", "serve" ]