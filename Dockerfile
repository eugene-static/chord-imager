FROM golang:1.21 as builder
WORKDIR /app
ENV CGO_ENABLED 0
COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./
RUN go build -o chord-drawer ./app/cmd
FROM alpine:latest
COPY --from=builder /app .
CMD ["./chord-drawer"]


