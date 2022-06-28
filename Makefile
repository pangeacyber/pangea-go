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
	go test ${TEST_TIMEOUT} -v -count=1 -race ${SDK_ALL_PKGS}
	@echo "Finished unit tests"

#######################
# Integration Testing #
#######################
integration:
	@echo "Started integration tests"
	go test -count=1 -tags "integration" -v -run '^Test_Integration' ./service/...
	@echo "Finished integration tests"
