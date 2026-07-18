<<<<<<< HEAD
=======
# 关键点：将版本从 1.22 升级到 1.26
>>>>>>> 0c61c4ead0f94abb72a83e688447dfd84c048591
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o nb-action .

FROM scratch
COPY --from=builder /app/nb-action /nb-action
ENTRYPOINT ["/nb-action", "server", "8080"]