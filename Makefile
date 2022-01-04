.PHONY: build

build:
	goreleaser release --snapshot --rm-dist

release: build
	docker push vanhtuan/dnsstat:armv7
	docker push vanhtuan/dnsstat:arm64
	docker push vanhtuan/dnsstat:amd64
	docker manifest create vanhtuan/dnsstat \
		vanhtuan/dnsstat:amd64 \
		vanhtuan/dnsstat:armv7 \
		vanhtuan/dnsstat:arm64
	docker manifest push vanhtuan/dnsstat
