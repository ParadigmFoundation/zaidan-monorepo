FROM golang:1.13-alpine
WORKDIR /src
COPY . .
RUN apk add make

RUN make -C services/order-book-manager build
ENTRYPOINT /src/services/order-book-manager/obm
