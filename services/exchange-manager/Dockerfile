FROM golang:1.13-alpine
WORKDIR /src
COPY . .
RUN apk add make

RUN make -C services/exchange-manager build
ENTRYPOINT ["/src/services/exchange-manager/em"]
