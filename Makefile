# Go Gardener Lab
#
#   make test EX=01        run the tests of exercise 01 (your code)
#   make test EX=12        exercise 12 needs envtest: downloads kube-apiserver + etcd once
#   make solution EX=01    run the tests against the reference solution
#   make test-all          all your exercises
#   make solutions         all reference solutions (should always pass)
#   make check             gofmt + go vet + tests-in-sync

ENVTEST_K8S_VERSION ?= 1.37.0
ENVTEST := $(shell go env GOPATH)/bin/setup-envtest

.PHONY: test solution test-all solutions check envtest

envtest:
	@test -x $(ENVTEST) || go install sigs.k8s.io/controller-runtime/tools/setup-envtest@release-0.24

test: envtest
	@test -n "$(EX)" || (echo "usage: make test EX=01" && exit 1)
	KUBEBUILDER_ASSETS="$$($(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path)" go test -race -count=1 $$(ls -d ./exercises/$(EX)-*)/...

solution: envtest
	@test -n "$(EX)" || (echo "usage: make solution EX=01" && exit 1)
	KUBEBUILDER_ASSETS="$$($(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path)" go test -race -count=1 $$(ls -d ./solutions/$(EX)-*)/...

test-all: envtest
	KUBEBUILDER_ASSETS="$$($(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path)" go test -count=1 ./exercises/...

solutions: envtest
	KUBEBUILDER_ASSETS="$$($(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path)" go test -race -count=1 ./solutions/...

check:
	@test -z "$$(gofmt -l .)" || (echo "gofmt needed:" && gofmt -l . && exit 1)
	go vet ./...
	./hack/check-tests-in-sync.sh
