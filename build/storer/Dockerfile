FROM golang:latest as builder

#Install glide for installing deps
RUN curl https://glide.sh/get | sh

WORKDIR /go/src/github.com/tcfw/evntsrc/
COPY . .

#Install deps and build
RUN glide install --strip-vendor
RUN go build -ldflags="-w -s" -o storer ./cmd/storer

# ---

FROM ubuntu:latest

COPY --from=builder /go/src/github.com/tcfw/evntsrc/storer /storer

CMD [ "/storer" ]