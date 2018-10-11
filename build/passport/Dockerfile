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
RUN go build -ldflags="-w -s" -o passport ./cmd/passport

# ---

FROM ubuntu:latest

WORKDIR /bin
COPY --from=builder /go/src/github.com/tcfw/evntsrc/passport passport

CMD [ "passport", "serve", "--tlsdir", "/certs" ]