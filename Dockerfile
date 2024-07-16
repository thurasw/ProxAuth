# Build stage: Compile go app
FROM golang:1.22.4-alpine AS builder
RUN apk update \
    && apk add --update gcc musl-dev

WORKDIR /app
ENV CGO_ENABLED=1 GOOS=linux

# Leverage docker layer cache by downloading dependencies first
COPY go.* .
RUN --mount=type=cache,target=/go/pkg/mod/ \
    go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod/ \
    go build -o "bin/ProxAuth" "./src/cmd"


# Build stage: Web app production build
FROM node:20-alpine AS web-builder
WORKDIR /app
COPY ./web .
RUN --mount=type=cache,target=/usr/local/share/.cache/yarn/v6 yarn --immutable
RUN yarn build

# Dev stage: Copy static executable for go app
FROM alpine AS development
# Create empty dir for DB path
WORKDIR /app/bin
COPY --from=builder /app/bin .

# These are inherited in prod stage below as well
EXPOSE 3000
ENV PORT=3000 WEB_ROOT_PATH=/app/web
ENTRYPOINT ["/app/bin/ProxAuth"]

# Prod stage: Copy production files for web app
FROM development AS production
WORKDIR /app/web
COPY --from=web-builder /app/dist .

ENTRYPOINT ["/app/bin/ProxAuth"]
