FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o nb-action .

FROM scratch
COPY --from=builder /app/nb-action /nb-action
ENTRYPOINT ["/nb-action", "server", "8080"]