ROOT := github.com/liubog2008/oooops
TARGETS := controller codectl imagectl

.PHONY: codegen compile

codegen:
	./hack/update-codegen.sh

compile: codegen
	@for target in $(TARGETS); do                                              \
		docker run --rm                                                        \
			-v $(PWD):/go/src/$(ROOT)                                          \
		    -w /go/src/$(ROOT)                                                 \
		    -e GOOS=linux                                                      \
		    -e GOARCH=amd64                                                    \
		    -e GOPATH=/go                                                      \
		golang:1.12.4-alpine3.9                                                \
		    go build                                                           \
				-v                                                             \
				-o ./_output/$${target}                                        \
		    ./cmd/$${target};                                                  \
	done
