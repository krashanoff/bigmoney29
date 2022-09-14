#
# Builds your installation to ./dist.
#

GO_TAGS=sqlite_foreign_keys
DIST_DIR=./dist

default: host

package: host
	tar czf host.tar.gz -C dist/ .

.PHONY: host
host: frontend backend doc

.PHONY: frontend
frontend:
	cd ui; \
	yarn; \
	BUILD_PATH='../$(DIST_DIR)/build' yarn build

.PHONY: backend
backend:
	cd backend; \
	go build -o '../$(DIST_DIR)/bigmoney29' --tags $(GO_TAGS)

.PHONY: doc
doc:
	mkdir -p $(DIST_DIR)/docs; \
	cp docs/*.md $(DIST_DIR)/docs

.PHONY: test
test:
	cd ui; \
	yarn run test; \
	cd ../backend; \
	go test
