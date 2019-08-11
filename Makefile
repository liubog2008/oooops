ROOT := github.com/liubog2008/oooops
TARGETS := controller codectl imagectl

.PHONY: codegen compile

codegen:
	./hack/update-codegen.sh

compile: codegen
	@for target in $(TARGETS); do                                          \
		go build                                                           \
			-v                                                             \
			-o ./_output/$${target}                                        \
		./cmd/$${target};                                                  \
	done
