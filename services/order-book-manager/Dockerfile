FROM gcr.io/zaidan-io/golang

WORKDIR /src
COPY . .
RUN apk add make

RUN make -C services/order-book-manager build
ENTRYPOINT ["/src/services/order-book-manager/obm"]
