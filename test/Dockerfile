FROM golang:1.22-alpine as builder

USER root

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -a -tags netgo,usergo -ldflags "-extldflags '-static' -s -w" -o natstest

FROM busybox:musl

COPY --from=builder /app/natstest /app/natstest

USER root
WORKDIR /app

ENTRYPOINT ["/app/natstest"]

