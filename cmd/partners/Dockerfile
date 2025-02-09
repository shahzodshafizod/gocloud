FROM golang:alpine AS builder

WORKDIR /api

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go test -count=1 -v github.com/shahzodshafizod/gocloud/internal/products/
RUN go test -count=1 -v github.com/shahzodshafizod/gocloud/internal/partners/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o partners cmd/partners/main.go

FROM scratch
COPY --from=builder /api/partners /

EXPOSE 4404
CMD ["/partners"]
