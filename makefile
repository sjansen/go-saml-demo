.PHONY:  default  boltdb  runserver  test  test-coverage  test-docker

default: runserver


boltdb:
	GSD_SESSION_STORE=boltdb \
	GSD_BOLTSTORE_PREFIX=data/ \
	go run main.go runserver


dynamodb:
	@scripts/docker-up-dynamodb


runserver:
	go run main.go runserver


test:
	@scripts/run-all-tests
	@echo ========================================
	@git grep TODO  -- '**.go' || true
	@git grep FIXME -- '**.go' || true


test-coverage: test-docker
	go tool cover -html=dist/coverage.txt


test-docker:
	@scripts/docker-up-test
