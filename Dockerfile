FROM golang:1.25-alpine AS builder

ENV CGO_ENABLED=0

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o /bin/app .

FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /bin/app /app/app

EXPOSE 9000

CMD ["./app"]