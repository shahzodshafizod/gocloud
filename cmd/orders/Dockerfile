FROM golang:alpine AS builder

WORKDIR /api

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go test -count=1 -v github.com/shahzodshafizod/gocloud/internal/orders/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o orders cmd/orders/main.go

FROM scratch

COPY --from=builder /api/orders /

EXPOSE 4403
CMD ["/orders"]
