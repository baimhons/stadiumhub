FROM golang:1.24.3-alpine AS builder

RUN apk add --no-cache git ca-certificates
WORKDIR /app
COPY ../../../backend/go.mod ../../../backend/go.sum ./
RUN go mod download && go mod verify
COPY ../../../backend/cmd/server/ ./cmd/server
COPY ../../../backend/internal/ ./internal/
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

FROM scratch AS runner
WORKDIR /app

# copy binary + CA certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/main .

CMD ["./main"]
