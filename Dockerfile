FROM golang:1.22-alpine AS builder

WORKDIR /build
COPY . .
RUN go mod download 
RUN go build -o ./shortener-url

FROM gcr.io/distroless/base-debian12 

WORKDIR /app
COPY --from=builder /build/shortener-url ./shortener-url
CMD ["/app/shortener-url"]
