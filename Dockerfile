FROM golang:1.26.8-bookworm AS builder

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /bin/server \
    ./cmd/server

FROM gcr.io/distroless/static-debian13:nonroot

COPY --from=builder /bin/server /server

EXPOSE 8080

ENTRYPOINT ["/server"]
