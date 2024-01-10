# SDK additional packages that are used for development of the SDK.
SDK_CORE_PKGS=./pangea/... ./internal/...
SDK_CLIENT_PKGS=./service/...
SDK_ALL_PKGS= ${SDK_CLIENT_PKGS} ${SDK_CORE_PKGS}
TEST_TIMEOUT=-timeout 5m

###################
# Unit Testing #
###################
unit:
	@echo "Started unit tests"
	go test ${TEST_TIMEOUT} -tags=unit -v -count=1 -race ${SDK_ALL_PKGS}
	@echo "Finished unit tests"

#######################
# Integration Testing #
#######################
integration:
	@echo "Started integration tests"
	go test -count=1 -tags integration -v ./service/...
	@echo "Finished integration tests"

coverage.out: $(shell find . -type f -print | grep -v vendor | grep "\.go")
	@go test -cover -coverprofile ./coverage.out.tmp ${SDK_ALL_PKGS}
	@cat ./coverage.out.tmp | grep -v '.pb.go' | grep -v 'mock_' > ./coverage.out
	@rm ./coverage.out.tmp

cover: coverage.out
	@echo ""
	@go tool cover -func ./coverage.out ${SDK_ALL_PKGS}

cover-html: coverage.out
	@go tool cover -html=./coverage.out ${SDK_ALL_PKGS}

clean:
	@rm ./coverage.out

##################
# Linting/Verify #
##################
verify: vet

vet:
	go vet -tags "example integration" --all ${SDK_ALL_PKGS}

fmt:
	@gofmt -l -w .

##########################
# Generate Documentation #
##########################
docgen:
	@echo "Generating docs JSON from docstrings"
	@go run ./autogendoc/main.go > ./go_sdk.json || exit 1
