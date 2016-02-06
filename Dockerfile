FROM golang:latest
COPY ./ $GOPATH/src/wxBot
WORKDIR $GOPATH/src/wxBot
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN go get
RUN go build
EXPOSE 4000
CMD ["./wxBot"]