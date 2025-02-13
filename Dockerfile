FROM golang:alpine AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd

FROM alpine
COPY --from=builder /app/app /usr/bin/app
EXPOSE 50051
CMD ["app"]