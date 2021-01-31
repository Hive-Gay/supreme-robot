FROM golang:1.15 AS builder
RUN go get github.com/markbates/pkger/cmd/pkger

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go test -v ./...
RUN pkger && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o supreme_robot

FROM scratch
COPY --from=builder /etc/ssl /etc/ssl

COPY --from=builder /app/supreme_robot /supreme_robot
CMD ["/supreme_robot"]