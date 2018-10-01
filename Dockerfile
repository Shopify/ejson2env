FROM golang:1.11-alpine AS builder
WORKDIR /go/src/github.com/Shopify/ejson2env
COPY . .
RUN go install -v github.com/Shopify/ejson2env/cmd/ejson2env

FROM scratch
COPY --from=builder /go/bin/ejson2env /
WORKDIR /tmp
ENTRYPOINT ["/ejson2env"]
