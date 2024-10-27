test-handler:
		cd api && go test ./cmd/internal/handler/rest -v
test-relationship-ctrl:
		cd api && go test ./cmd/internal/controller/relationship -v
test-user-ctrl:
		cd api && go test ./cmd/internal/controller/user -v
test-relationship-repo:
		cd api && go test ./cmd/internal/repository/relationship -v
test-user-repo:
		cd api && go test ./cmd/internal/repository/user -v
test:
		cd api && go test ./... -v