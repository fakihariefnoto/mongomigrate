usage: FORCE
	# See targets in Makefile (e.g. "buildlet.darwin-amd64")
	exit 1

FORCE:

mm: FORCE
	@go build ./cmd/mm
	@./mm

test:
	@echo " >> starting go test .."
	@go test -v -cover ./...
	@echo "test done.."

all: FORCE
	@echo " >> building main binaries..."
	@go build ./cmd/mm
	@echo " >> mm app has been built."


