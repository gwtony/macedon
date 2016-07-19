all: napi

napi:
	@mkdir -p dist/bin dist/conf dist/log
	@sh -c "'$(CURDIR)/scripts/build.sh'"
test:
	@sh -c "'$(CURDIR)/scripts/test.sh'"
cover:
	@sh -c "'$(CURDIR)/scripts/test.sh' cover"


clean:
	@rm -rf dist

.PHONY: all napi clean test
