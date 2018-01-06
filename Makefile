.PHONY: test unit integration

test: unit integration

unit:
	ginkgo -r --skipPackage=integration_tests

integration:
	ginkgo -r integration_tests/
