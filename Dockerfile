FROM golang:1.22.4 AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:latest
RUN adduser -D appuser
USER appuser

WORKDIR /app

COPY --from=builder /build/main .

EXPOSE 8089

CMD ["./main"]
