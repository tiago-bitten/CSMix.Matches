# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS build
WORKDIR /src

# go.sum is optional on purpose: this service has no dependencies yet, so the
# file does not exist. The glob keeps the COPY working either way, and the
# moment a dependency arrives this layer starts earning its cache.
COPY go.mod go.sum* ./
RUN go mod download

COPY . .

# CGO off is what makes the binary static, which is what lets the final stage be
# a base image with no libc in it at all. -trimpath keeps build paths out of the
# binary; -s -w drops the symbol and DWARF tables.
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

# Postgres, Kafka and CSMix.GameServers are reached from here, so the CA
# bundle in distroless/static is what makes TLS to any of them work.
FROM gcr.io/distroless/static-debian12:nonroot AS final

COPY --from=build /out/api /api

USER nonroot:nonroot

EXPOSE 8082

# No shell in here, which is the point: nothing to exec into and nothing for an
# injected command to run. Use the :debug tag of the same image when you need
# one. Probing is left to whatever runs the container.
ENTRYPOINT ["/api"]
