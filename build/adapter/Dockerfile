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

WORKDIR /evntsrc/

#Install deps and build
ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .
RUN go mod download

#Install V8
COPY --from=libv8 /v8/include /go/pkg/mod/gopkg.in/augustoroman/v8.v1@v1.0.0/include
COPY --from=libv8 /v8/lib /go/pkg/mod/gopkg.in/augustoroman/v8.v1@v1.0.0/libv8

#Build
COPY . .
RUN go build -ldflags="-w -s" -o adapter ./cmd/adapter
RUN strip --strip-unneeded adapter
# RUN upx adapter - !! UPX has unexpected results with compressing

# ---

FROM ubuntu:latest

WORKDIR /bin
COPY --from=builder /evntsrc/adapter adapter

CMD [ "adapter", "serve" ]