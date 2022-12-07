FROM golang:1.19.4-alpine AS builder
WORKDIR /go/src/github.com/Shopify/ejson2env
COPY . .
RUN apk add --no-cache git && \
    go install -v github.com/Shopify/ejson2env/v2/cmd/ejson2env

FROM scratch
COPY --from=builder /go/bin/ejson2env /
WORKDIR /tmp
ENTRYPOINT ["/ejson2env"]
