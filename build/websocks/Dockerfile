FROM golang:latest as builder

#Install XZ
RUN apt-get update && apt-get install -y \
	xz-utils \
	&& rm -rf /var/lib/apt/lists/*

#Install UPX
ADD https://github.com/upx/upx/releases/download/v3.94/upx-3.94-amd64_linux.tar.xz /usr/local
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

#Build
COPY . .
RUN go build -ldflags="-w -s" -o websocks ./cmd/websocks
RUN strip --strip-unneeded websocks
RUN upx websocks

# ---

FROM ubuntu:latest 

RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /bin
COPY --from=builder /go/src/github.com/tcfw/evntsrc/websocks websocks

CMD [ "websocks" ]