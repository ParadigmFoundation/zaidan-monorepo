FROM golang:1.13-alpine
WORKDIR /src
COPY . .
RUN apk add make

RUN make -C services/dealer build
ENTRYPOINT ["/src/services/dealer/dealer"]
