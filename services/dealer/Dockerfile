FROM golang:1.13
WORKDIR /src
COPY . .
RUN make -C services/dealer build

FROM alpine:3.8
RUN apk --no-cache add ca-certificates
COPY --from=0 /src/services/dealer/dealer /usr/bin
ENTRYPOINT ["/usr/bin/dealer"]
