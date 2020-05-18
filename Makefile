ROOT := github.com/liubog2008/oooops
TARGETS := operator
REGISTRY := registry.cn-hangzhou.aliyuncs.com
GROUP := liubog2008
PROJECT := oooops
NAMESPACE := oooops
VERSION := `./hack/version.sh DOCKER_VERSION`
LDFLAGS := `./hack/version.sh`

.PHONY: crd codegen compile build push deploy load-to-kind test reload-operator logs-operator

crd:
	controller-gen crd:crdVersions=v1,preserveUnknownFields=false paths=$(PWD)/pkg/apis/... output:crd:dir=$(PWD)/crd/

codegen:
	./hack/codegen.sh

build:
	rm -rf _output
	mkdir _output
	@for target in $(TARGETS); do                                       \
		go build                                                        \
			-v                                                          \
			--ldflags "$(LDFLAGS)"                                      \
			-o ./_output/$${target}                                     \
		./cmd/$${target};                                               \
	done

container:
	rm -rf _output
	mkdir _output
	@for target in $(TARGETS); do                                            \
		docker run                                                           \
			--rm                                                             \
			-w /go/src/$(ROOT)                                               \
			-v $(PWD):/go/src/$(ROOT)                                        \
			-v $(GOCACHE):/go/.cache                                         \
			-v $(GOPATH)/pkg/mod:/go/pkg/mod                                 \
			-e GO111MODULE=on                                                \
			-e GOCACHE=/go/.cache                                            \
			-e GOPROXY=https://goproxy.io                                    \
			golang:1.12.5-alpine3.9                                          \
			go build                                                         \
				-o _output/$${target}                                        \
				-v                                                           \
				--ldflags "$(LDFLAGS)"                                       \
				./cmd/$${target};                                            \
		docker build                                                         \
			-t $(REGISTRY)/$(GROUP)/$(PROJECT)-$${target}:$(VERSION)     \
			-f $(PWD)/build/$${target}/Dockerfile .;                         \
	done

push:
	@for target in $(TARGETS); do                                       \
		docker push                                                     \
			$(REGISTRY)/$(GROUP)/$(PROJECT)-$${target}:$(VERSION);  \
	done

load-to-kind:
	@for target in $(TARGETS); do \
		kind load docker-image --name=test \
			$(REGISTRY)/$(GROUP)/$(PROJECT)-$${target}:$(VERSION); \
	done

deploy:
	cat $(PWD)/deploy/namespace.yaml | \
		NAMESPACE=$(NAMESPACE) \
		envsubst | \
		kubectl apply -f -
	kubectl apply -f $(PWD)/crd/
	@for target in $(TARGETS); do                      \
		cat $(PWD)/deploy/$${target}/$${target}.yaml | \
			VERSION=$(VERSION)                         \
			REGISTRY=$(REGISTRY)                       \
			GROUP=$(GROUP)                             \
			PROJECT=$(PROJECT)                         \
			NAMESPACE=$(NAMESPACE)                     \
			envsubst |                                 \
			kubectl apply -f -;                        \
	done

test:
	kubectl delete pipes --all -n $(NAMESPACE)
	kubectl delete flows --all -n $(NAMESPACE)
	kubectl delete events.mario.oooops.com --all -n $(NAMESPACE)
	kubectl delete jobs --all -n $(NAMESPACE)
	cat $(PWD)/test/testdata/test-pipe.yaml | \
		NAMESPACE=$(NAMESPACE) \
		envsubst | \
		kubectl apply -f -
	cat $(PWD)/test/testdata/test-event.yaml | \
		NAMESPACE=$(NAMESPACE) \
		envsubst | \
		kubectl apply -f -
	cat $(PWD)/test/testdata/test-flow.yaml | \
		NAMESPACE=$(NAMESPACE) \
		envsubst | \
		kubectl apply -f -

reload-operator:
	kubectl delete pod -n $(NAMESPACE) \
		`kubectl get pods -n $(NAMESPACE) | awk '/operator/{ print $$1 }'`
logs-operator:
	kubectl logs -f -n $(NAMESPACE) \
		`kubectl get pods -n $(NAMESPACE) | awk '/operator/{ print $$1 }'`
