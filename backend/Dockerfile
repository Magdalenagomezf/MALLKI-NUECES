# ---- Etapa 1: build ----
FROM golang:1.24-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /app/server .

# ---- Etapa 2: runtime ----
FROM alpine:3.20
WORKDIR /app
COPY --from=build /app/server .
EXPOSE 8080
ENTRYPOINT ["./server"]