FROM gcr.io/zaidan-io/golang

WORKDIR /src
COPY . .
RUN apk add make

RUN make -C services/watcher build
ENTRYPOINT ["/src/services/watcher/watcher"]
