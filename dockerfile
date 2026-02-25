FROM golang:1.25-alpine AS go-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go


FROM alpine:3.19 AS production

WORKDIR /app

COPY --from=go-builder /app/main /main
COPY --from=go-builder /app/migrations /app/migrations
COPY .env .env

ENTRYPOINT ["/main"]