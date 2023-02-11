##
# Important! Before running any command make sure you have setup GOPATH:
# export GOPATH="$HOME/go"
# PATH="$GOPATH/bin:$PATH"

start:
	# Start the application with postgresql database
	./scripts/start.sh

mocks:
	mockgen \
		-package mocks -destination=internal/mocks/mock_task_repository.go \
		-package mocks github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/domain TaskRepository

	mockgen \
		-package mocks -destination=internal/mocks/mock_task_redis_repository.go \
		-package mocks github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/domain TaskRedisRepository

	mockgen \
		-package mocks -destination=internal/mocks/mock_task_service.go \
		-package mocks github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/domain TaskService

	mockgen \
		-package mocks -destination=internal/mocks/mock_event_service.go \
		-package mocks github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/domain EventService