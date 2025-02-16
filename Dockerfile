FROM golang:1.21.0 as builder

WORKDIR /app
COPY . /app

RUN go mod download
RUN CGO_ENABLED=0 go build -o app

FROM golang:1.21-alpine as store

RUN apk add -U tzdata

WORKDIR /app
COPY --from=builder /app .

CMD ["./app"]
