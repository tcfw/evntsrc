FROM golang:latest as builder

#Install glide for installing deps
RUN curl https://glide.sh/get | sh

WORKDIR /go/src/github.com/tcfw/evntsrc/
COPY . .

#Install deps and build
RUN glide install --strip-vendor
RUN go build -ldflags="-w -s" -o streams ./cmd/streams

# ---

FROM ubuntu:latest

WORKDIR /bin
COPY --from=builder /go/src/github.com/tcfw/evntsrc/streams streams

CMD [ "streams", "serve" ]