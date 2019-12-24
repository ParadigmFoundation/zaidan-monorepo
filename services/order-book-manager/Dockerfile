FROM golang:1.13
WORKDIR /src
COPY . .
RUN make -C services/order-book-manager build
ENTRYPOINT ["/src/services/order-book-manager/obm"]
