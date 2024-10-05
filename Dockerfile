FROM golang:1.22

COPY . /build
WORKDIR /build

RUN go build -o /opt/service /build/cmd/service/main.go

ENTRYPOINT ["/opt/service"]