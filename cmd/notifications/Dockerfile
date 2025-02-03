FROM golang:alpine

WORKDIR /api

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go test -count=1 -v github.com/shahzodshafizod/gocloud/internal/notifications/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o notifications cmd/notifications/main.go

EXPOSE 4402
CMD ["/api/notifications"]
