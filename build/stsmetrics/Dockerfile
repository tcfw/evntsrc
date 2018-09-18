FROM golang:latest as builder

#Install glide for installing deps
RUN curl https://glide.sh/get | sh

WORKDIR /go/src/github.com/tcfw/evntsrc/
COPY . .

#Install deps and build
RUN glide install --strip-vendor
RUN go build -ldflags="-w -s" -o stsmetrics ./cmd/stsmetrics

# ---

FROM ubuntu:latest 

RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /bin
COPY --from=builder /go/src/github.com/tcfw/evntsrc/stsmetrics stsmetrics

CMD [ "stsmetrics" ]