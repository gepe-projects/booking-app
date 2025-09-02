# Stage 1: base
FROM golang:1.25.0-alpine AS base
WORKDIR /app

# Install Air untuk hot reload + utilitas
RUN apk add --no-cache bash build-base curl git \
    && go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Stage 2: dev (pakai Air)
FROM base AS dev
WORKDIR /app
CMD ["air", "-c", ".air.toml"]

# Stage 3: production (tanpa Air)
FROM golang:1.22-alpine AS prod
WORKDIR /app

COPY --from=base /app /app

RUN go build -o booking-app ./cmd/server

CMD ["./booking-app"]
