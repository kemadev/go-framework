ARG ALPINE_VERSION=3.21.3

FROM --platform=${BUILDPLATFORM} docker.io/alpine:${ALPINE_VERSION} AS alpine-base

FROM alpine-base AS go-deps

SHELL ["/bin/ash", "-o", "pipefail", "-c"]

# Docker buildx ARGs
ARG BUILDPLATFORM

WORKDIR /deps

# hadolint ignore=DL3018
RUN apk add --no-cache curl grep tar gzip zip

ENV BUILDPLATFORM=${BUILDPLATFORM}

RUN --mount=type=bind,source=go.mod,target=go.mod curl -fsSL "https://go.dev/dl/go$(grep -oP 'go\s+\K([0-9.]+)' <go.mod).$(echo ${BUILDPLATFORM} | sed 's|/|-|g').tar.gz" | tar -C /usr/local -xz

# Ensure future GOPATH exists
RUN mkdir -p /root/go

FROM alpine-base AS clean-deps

WORKDIR /deps

COPY --link --from=go-deps /usr/local/go /usr/local/go
COPY --link --from=go-deps /root/go /root/go

ENV PATH="/root/go/bin:/usr/local/go/bin:${PATH}"

FROM --platform=${BUILDPLATFORM} clean-deps AS base

WORKDIR /src

# Docker buildx ARGs
ARG TARGETOS
ARG TARGETARCH

ARG CGO_ENABLED=0
ARG GOAMD64=v3
ARG GOARM64=v8.2

ENV CGO_ENABLED=${CGO_ENABLED} GOOS=${TARGETOS} GOARCH=${TARGETARCH} GOARM64=${GOARM64} GOAMD64=${GOAMD64}

RUN --mount=type=secret,id=netrc \
	--mount=type=cache,target=/root/go/pkg \
	--mount=type=bind,source=go.mod,target=go.mod \
	--mount=type=bind,source=go.sum,target=go.sum \
	cp /run/secrets/netrc /root/.netrc && \
	go mod download -x

FROM --platform=${BUILDPLATFORM} base AS build

RUN --mount=type=cache,target=/root/go/pkg \
	--mount=type=cache,target=/root/.cache \
	--mount=type=bind,source=.,target=. \
	go build -ldflags='-s -w' -o /usr/local/bin/app ./cmd/main

FROM --platform=${BUILDPLATFORM} base AS debug-build

RUN --mount=type=cache,target=/root/go/pkg \
		--mount=type=cache,target=/root/.cache \
		--mount=type=bind,source=.,target=. \
		go build -gcflags="all=-N -l" -o /usr/local/bin/app ./cmd/main

FROM gcr.io/distroless/static-debian12:nonroot AS goreleaser

LABEL org.opencontainers.image.authors="kemadev"
LABEL org.opencontainers.image.url="ghcr.io/kemadev/go-framework"
LABEL org.opencontainers.image.documentation="https://github.com/kemadev/go-framework"
LABEL org.opencontainers.image.source="https://github.com/kemadev/go-framework"
LABEL org.opencontainers.image.vendor="kemadev"
LABEL org.opencontainers.image.licenses="MPL-2.0"

COPY go-framework /app

ENTRYPOINT ["/app"]

FROM alpine-base AS dev

WORKDIR /src

# hadolint ignore=DL3018
RUN apk add --no-cache git sudo && \
	adduser -D -u 1000 nonroot && \
	echo "nonroot ALL=(ALL) NOPASSWD: ALL" >/etc/sudoers.d/nonroot && \
	chmod 0440 /etc/sudoers.d/nonroot

COPY --link --from=go-deps /usr/local/go /usr/local/go
COPY --link --chown=nonroot:nonroot --from=go-deps /root/go /home/nonroot/go

ENV PATH="/home/nonroot/go/bin:/usr/local/go/bin:${PATH}"

USER 1000:1000

RUN go install github.com/go-delve/delve/cmd/dlv@latest && \
	git config --global --add safe.directory /src

COPY --link --from=debug-build /usr/local/bin/app /app
COPY --link ./tool/dev/docker-entrypoint-dev.sh /src/tool/dev/docker-entrypoint-dev.sh

ENTRYPOINT ["/src/tool/dev/docker-entrypoint-dev.sh"]

# hadolint ignore=DL3007
FROM gcr.io/distroless/static-debian12:nonroot AS main

LABEL org.opencontainers.image.authors="kemadev"
LABEL org.opencontainers.image.url="ghcr.io/kemadev/go-framework"
LABEL org.opencontainers.image.documentation="https://github.com/kemadev/go-framework"
LABEL org.opencontainers.image.source="https://github.com/kemadev/go-framework"
LABEL org.opencontainers.image.vendor="kemadev"
LABEL org.opencontainers.image.licenses="MPL-2.0"

COPY --link --from=build /usr/local/bin/app /app

ENTRYPOINT ["/app"]
