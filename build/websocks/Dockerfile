FROM golang:latest as builder

#Install glide for installing deps
RUN curl https://glide.sh/get | sh

WORKDIR /go/src/github.com/tcfw/evntsrc/

#Install deps and build
ADD glide.yaml glide.yaml
ADD glide.lock glide.lock
RUN glide install --strip-vendor

#Build
COPY . .
RUN go build -ldflags="-w -s" -o websocks ./cmd/websocks

# ---

FROM ubuntu:latest 

RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /bin
COPY --from=builder /go/src/github.com/tcfw/evntsrc/websocks websocks

CMD [ "websocks" ]