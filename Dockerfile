# ---- Build stage ----
FROM golang:1.22-alpine AS builder

WORKDIR /build

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o carbon-emission-management ./cmd/api

# ---- Runtime stage ----
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /build/carbon-emission-management .
COPY --from=builder /build/configs ./configs

EXPOSE 8080

CMD ["./carbon-emission-management"]
