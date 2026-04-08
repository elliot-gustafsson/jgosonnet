.PHONY:
SHELL := /bin/bash
# MAKEFLAGS += --no-print-directory


build:
	cd cmd/jgosonnet && CGO_ENABLED=0 GOOS=linux go build

cpu-profile:
	go tool pprof -http=:8080 cpu.prof

# localrun:
# 	go run main.go

# generate:
# 	go generate main.go

test:
	go test ./... -count=1 -cover

test-coverage:
	go test ./... -count=1 -coverprofile=coverage.out
	go tool cover -func=coverage.out

benchmark:
	mkdir -p benchmarks/out
	cd benchmarks && go test -benchmem -bench=. -cpuprofile=out/cpu.prof -memprofile=out/mem.prof -o=out/benchmarks.test -count=1 -v


benchmark-prof-cpu:
	go tool pprof -http=:8080 benchmarks/out/cpu.prof

benchmark-prof-mem:
	go tool pprof -http=:8080 benchmarks/out/mem.prof
