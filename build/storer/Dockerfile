FROM golang:latest as builder

WORKDIR /go/src/github.com/tcfw/evntsrc/
COPY . .

RUN go build -ldflags="-w -s" -o storer ./cmd/storer

FROM ubuntu:latest 

RUN apt-get update && apt-get install -y ca-certificates

COPY --from=builder /go/src/github.com/tcfw/evntsrc/storer /storer

CMD [ "/storer" ]