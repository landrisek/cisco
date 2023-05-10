FROM golang:1.13.6
WORKDIR /
ADD . ./
RUN go get .
RUN go build -o container
RUN mkdir -p /data
ENTRYPOINT ./container