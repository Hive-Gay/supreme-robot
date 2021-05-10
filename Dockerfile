FROM golang:1.16 AS builder
RUN go get github.com/markbates/pkger/cmd/pkger

WORKDIR /go/src/github.com/Hive-Gay/supreme-robot

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN pkger && \
    CGO_ENABLED=0 go build -a -installsuffix cgo -o supreme-robot

FROM scratch

COPY --from=builder /etc/ssl /etc/ssl
COPY --from=builder /go/src/github.com/Hive-Gay/supreme-robot/supreme-robot /supreme-robot

WORKDIR /
CMD ["/supreme-robot"]