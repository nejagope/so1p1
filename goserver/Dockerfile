FROM golang:alpine
RUN mkdir /goserver
ADD /pkg/server/ /goserver/
WORKDIR /goserver
RUN go build -o main .
RUN adduser -S -D -H -h /goserver goserveruser
USER goserveruser
EXPOSE 8080
CMD ["./main"]