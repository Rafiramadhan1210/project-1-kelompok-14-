FROM golang:1.25-bookworm AS build
WORKDIR /src
COPY backend/go.mod backend/go.sum ./backend/
WORKDIR /src/backend
RUN go mod download
WORKDIR /src
COPY backend ./backend
WORKDIR /src/backend
RUN CGO_ENABLED=0 go build -tags netgo -ldflags '-s -w' -o /app/app .

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=build /app/app ./app
COPY frontend/public ./frontend/public
EXPOSE 8080
CMD ["./app"]
