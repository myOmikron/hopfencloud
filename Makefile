GO_COMPILER ?= go
GO_FLAGS ?= -ldflags=-w

BINARY_NAME ?= hopfencloud
BUILD_DIR ?= ./bin

INSTALL_TARGET ?= /usr/local/bin/

SYSTEMD_DIR ?= `pkg-config systemd --variable=systemdsystemunitdir`

.PHONY: all build install clean uninstall purge .purge

all: build

build:
	$(GO_COMPILER) build $(GO_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) cmd/hopfencloud/main.go

install:
	install -s -m 0755 $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_TARGET)
	install -m 0644 hopfencloud.service $(SYSTEMD_DIR)
	mkdir -p /etc/hopfencloud/ /var/lib/private/hopfencloud/
	install -m 0640 example.config.toml /etc/hopfencloud/example.config.toml
	cp -r templates/ /var/lib/private/hopfencloud/
	cp -r static/ /var/lib/private/hopfencloud/
	systemctl enable hopfencloud.service

clean:
	$(GO_COMPILER) clean
	$(RM) $(BUILD_DIR)/$(BINARY_NAME)

uninstall:
	systemctl stop hopfencloud.service |:
	$(RM) $(INSTALL_TARGET)/$(BINARY_NAME) $(SYSTEMD_DIR)/hopfencloud.service
	systemctl disable hopfencloud.service |:
	systemctl daemon-reload

purge: uninstall .purge
.purge:
	$(RM) -r /etc/hopfencloud /var/lib/hopfencloud /var/lib/private/hopfencloud \
	/var/log/hopfencloud /var/log/private/hopfencloud /run/hopfencloud
