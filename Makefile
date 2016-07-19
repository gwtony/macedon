all: macedon

macedon:
	@mkdir -p dist/bin dist/conf dist/log
	@sh -c "'$(CURDIR)/scripts/build.sh'"

clean:
	@rm -rf dist

.PHONY: all macedon clean test
