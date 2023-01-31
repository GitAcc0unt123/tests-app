# Step 1: Modules caching
FROM golang:1.19 AS modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.19 AS builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build

# Step 3: Final image with binary
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/config.yml /app/.env /app
COPY --from=builder /app/tests_app /app/tests_app
COPY ./static/ ./static/
EXPOSE 8080
CMD ["./tests_app"]