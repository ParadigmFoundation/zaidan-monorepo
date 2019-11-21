FROM golang:1.13
WORKDIR /src
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 go build -mod=vendor -o obm

FROM alpine:3.8
RUN apk --no-cache add ca-certificates
COPY --from=0 /src/obm /usr/bin
ENTRYPOINT ["obm"]
