FROM golang:1.13-alpine
WORKDIR /src
COPY . .
RUN apk add make gcc libc-dev linux-headers

RUN make -C services/dealer build
ENTRYPOINT ["/src/services/dealer/dealer"]
