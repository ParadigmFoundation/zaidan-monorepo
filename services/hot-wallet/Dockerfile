FROM gcr.io/zaidan-io/golang

WORKDIR /src
COPY . .
RUN apk add make

RUN make -C services/hot-wallet build
ENTRYPOINT ["/src/services/hot-wallet/hwd"]
