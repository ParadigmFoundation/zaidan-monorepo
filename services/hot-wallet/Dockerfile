FROM golang:1.13
WORKDIR /src
COPY . .
RUN make -C services/hot-wallet build

FROM alpine:3.8
RUN apk --no-cache add ca-certificates
COPY --from=0 /src/services/hot-wallet/hwd /usr/bin
ENTRYPOINT ["hwd"]
