all: main

main:
	@mkdir -p dist/bin dist/conf dist/log
	@sh -c "'$(CURDIR)/scripts/build.sh'"
debug:
	@mkdir -p dist/bin dist/conf dist/log
	@sh -c "'$(CURDIR)/scripts/build.sh' debug"
test:
	@sh -c "'$(CURDIR)/scripts/test.sh'"
cover:
	@sh -c "'$(CURDIR)/scripts/test.sh' cover"


clean:
	@rm -rf dist

.PHONY: all main clean test debug
