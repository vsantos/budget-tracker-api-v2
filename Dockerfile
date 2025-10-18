# Stage 1: build
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o budget-tracker-api-v2 .

FROM gcr.io/distroless/base-debian10:nonroot

WORKDIR /app

COPY --from=builder /app/swagger ./swagger/
COPY --from=builder /app/budget-tracker-api-v2 .

CMD ["./budget-tracker-api-v2"]
