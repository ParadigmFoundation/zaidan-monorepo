FROM golang:1.13-alpine
WORKDIR /src
COPY . .
RUN apk add make

RUN make -C services/hot-wallet build
ENTRYPOINT /src/services/hot-wallet/hwd
