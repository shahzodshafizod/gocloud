FROM golang:alpine as builder

WORKDIR /api

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go test -count=1 -v github.com/shahzodshafizod/gocloud/internal/gateway/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api cmd/api/main.go

EXPOSE 4401
CMD ["/api/api"]

# FROM scratch

# COPY --from=builder /api/api /

# EXPOSE 4401
# ENTRYPOINT ["/api"]
