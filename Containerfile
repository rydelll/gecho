FROM docker.io/golang:1.23.3-alpine3.20 AS build
WORKDIR /app

ENV GOOS=linux \
    GOARCH=amd64 \
    CGO_ENABLED=0 \
    GOCACHE=/go/cache \
    GOMODCACHE=/go/modcache

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/modcache \
    go mod download

COPY . .
RUN --mount=type=cache,target=/go/cache \
    --mount=type=cache,target=/go/modcache \
    go build -o gecho ./main.go

FROM scratch
USER 1000:1000
WORKDIR /app

COPY --from=build /app/gecho /usr/bin/gecho

EXPOSE 7777
ENTRYPOINT ["gecho", "-port", "7777"]
