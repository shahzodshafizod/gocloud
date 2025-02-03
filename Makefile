protos-generate:
	protoc \
		--go_out=internal \
		--go_opt=module=github.com/shahzodshafizod/gocloud/internal \
		--go-grpc_out=internal \
		--go-grpc_opt=module=github.com/shahzodshafizod/gocloud/internal \
		internal/protos/*.proto

mocks-generate:
	go generate -v ./...
	mockgen -source=internal/partners/partners_grpc.pb.go -destination=internal/partners/mocks/partners_grpc.pb.go -package=mocks
	mockgen -source=internal/orders/orders_grpc.pb.go -destination=internal/orders/mocks/orders_grpc.pb.go -package=mocks
	mockgen -source=internal/products/service.go -destination=internal/products/mocks/service.go -package=mocks

swagger-generate:
	swag fmt
	swag init -g cmd/api/main.go internal/user/handlers.go

tests-run:
	go test -count=1 ./... 2>&1 | grep -v "no test files"

tests-cover:
	go test -coverprofile=internal/gateway/test-cover.out -count=1 -v ./internal/gateway/
	go tool cover -html=internal/gateway/test-cover.out

	go test -coverprofile=internal/notifications/test-cover.out -v -count=1 ./internal/notifications/
	go tool cover -html=internal/notifications/test-cover.out

	go test -coverprofile=internal/orders/test-cover.out -v -count=1 ./internal/orders/
	go tool cover -html=internal/orders/test-cover.out

	go test -coverprofile=internal/partners/test-cover.out -v -count=1 ./internal/partners/
	go tool cover -html=internal/partners/test-cover.out

	go test -coverprofile=internal/products/test-cover.out -v -count=1 ./internal/products/
	go tool cover -html=internal/products/test-cover.out

tests-clear:
	rm internal/gateway/test-cover.out
	rm internal/notifications/test-cover.out
	rm internal/orders/test-cover.out
	rm internal/partners/test-cover.out
	rm internal/products/test-cover.out

migration-create:
	# goose -dir "migrations/notifications" create notifications_init sql
	# goose postgres "host=localhost user=postgres database=db_name password=postgres sslmode=disable" status

	# name=partners_init dir=migrations/partners make create-migration
	migrate create -ext sql -dir ${dir} -seq ${name}

images-build:
	docker build -t delivery/api -f cmd/api/Dockerfile .
	docker build -t delivery/notifications -f cmd/notifications/Dockerfile .
	docker build -t delivery/orders -f cmd/orders/Dockerfile .
	docker build -t delivery/partners -f cmd/partners/Dockerfile .
