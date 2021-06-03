FROM golang:1.14.3-alpine
ENV GO111MODULE=on
ENV BROWSERTASK=/go/src/github.com/chain5j/sync_eth
ENV GOPROXY=https://goproxy.cn,direct
COPY . $BROWSERTASK
WORKDIR $BROWSERTASK
ADD go.mod .
RUN cd $BROWSERTASK  && go mod download
RUN cd $BROWSERTASK/ && GOOS=linux CGO_ENABLED=0 go build -o sync_eth

FROM alpine:edge
ENV BROWSERTASK=/go/src/github.com/chain5j/sync_eth
COPY --from=0  $BROWSERTASK/sync_eth /usr/bin
WORKDIR /data
CMD  ["sync_eth"]
