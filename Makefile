GOBIN := $(shell pwd)/bin
VERSION ?= $(shell git rev-parse --abbrev-ref HEAD)
COMMIT ?= $(shell git describe --match=NeVeRmAtCh --always --abbrev=40 --dirty)
BUILDFLAGS ?= -gcflags '$(GCFLAGS)' -ldflags '$(LDFLAGS) -X main.Version=$(VERSION) -X main.Commit=$(COMMIT)' -tags '$(BUILD_TAGS)'

.PHONY: cmd_%

default: cmd_csi-rclone

cmd_%: OUTPUT=$(patsubst cmd_%,./bin/%,$@)
cmd_%: SOURCE=$(patsubst cmd_%,./cmd/%,$@)
cmd_%:
	go build $(BUILDFLAGS) -o $(OUTPUT) $(SOURCE)

docker-release:
	docker buildx build . --platform linux/arm64,linux/amd64 -f docker/Dockerfile -t ghcr.io/lostb1t/csi-driver-rclone --push
