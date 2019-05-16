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

#Build
COPY . .
RUN go build -ldflags="-w -s" -o interconnect ./cmd/interconnect
RUN strip --strip-unneeded interconnect
RUN upx interconnect

# ---

FROM ubuntu:latest

WORKDIR /bin
COPY --from=builder /go/src/github.com/tcfw/evntsrc/interconnect interconnect

CMD [ "interconnect", "serve" ]