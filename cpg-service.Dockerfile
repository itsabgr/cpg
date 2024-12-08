FROM golang AS build
WORKDIR /app
COPY . .
RUN mkdir -p build
ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct
ENV CGO_ENABLED=0
RUN go build -ldflags="-linkmode external -extldflags -static" -tags netgo -o build/cpg-service ./cmd/cpg-service

FROM alpine
WORKDIR /
COPY --from=build /app/build/cpg-service /cpg-service

ENV ASSETS_CONFIG=assets.json
ENV SALT_KEYRING=salt.txt
ENV BACKUP_KEYRING=backup.txt
ENV GRPC_SERVER=0.0.0.0:9090
ENV RATE_LIMITER=MEMORY
ENV PG_URI=postgres://postgres/postgres?sslmode=disable

ENTRYPOINT ["/cpg-service"]
