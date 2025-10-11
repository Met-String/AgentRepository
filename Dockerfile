FROM golang:1.24.2-alpine AS build
WORKDIR /src

# Download dependencies first to leverage layer caching.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the project and build the binary.
COPY cmd ./cmd
COPY internal ./internal
COPY assets ./assets
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/gateway ./cmd

FROM alpine:3.20
WORKDIR /app

# Copy the compiled binary and the static assets.
COPY --from=build /out/gateway /usr/local/bin/gateway
COPY assets ./assets

EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/gateway"]

# docker run --rm -p 8080:8080 agentsquare