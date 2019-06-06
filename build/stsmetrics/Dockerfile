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

#Build
COPY . .
RUN go build -ldflags="-w -s" -o stsmetrics ./cmd/stsmetrics
RUN strip --strip-unneeded stsmetrics
RUN upx stsmetrics

# ---

FROM ubuntu:latest 

RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /bin
COPY --from=builder /evntsrc/stsmetrics stsmetrics

CMD [ "stsmetrics" ]