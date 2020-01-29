FROM gcr.io/zaidan-io/golang

WORKDIR /src
COPY . .
RUN apk add make

RUN make -C services/exchange-manager build
ENTRYPOINT ["/src/services/exchange-manager/em"]
