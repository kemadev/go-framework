# TODO Add renovate tracking
FROM --platform=${BUILDPLATFORM} docker.io/alpine:3.21.3 AS go-deps

SHELL ["/bin/ash", "-o", "pipefail", "-c"]

# Docker buildx ARGs
ARG TARGETOS
ARG TARGETARCH

WORKDIR /deps

# hadolint ignore=DL3018
RUN apk add --no-cache curl grep tar gzip zip

ENV GOOS=${TARGETOS} GOARCH=${TARGETARCH}

RUN --mount=type=bind,source=go.mod,target=go.mod curl -sSL "https://go.dev/dl/go$(grep -oP 'go\s+\K([0-9.]+)' <go.mod).${GOOS}-${GOARCH}.tar.gz" | tar -C /usr/local -xz

# Ensure future GOPATH exists
RUN mkdir -p /root/go

# TODO Add renovate tracking
FROM --platform=${BUILDPLATFORM} docker.io/alpine:3.21.3 AS clean-deps

# Docker buildx ARGs
ARG TARGETOS
ARG TARGETARCH
# Disable C compatability, pass build arg ENABLE_CGO=1 to enable
ARG ENABLE_CGO=0

WORKDIR /deps

COPY --link --from=go-deps /usr/local/go /usr/local/go
COPY --link --from=go-deps /root/go /root/go

ENV CGO_ENABLED=${ENABLE_CGO} GOOS=${TARGETOS} GOARCH=${TARGETARCH}
ENV PATH="/root/go/bin:/usr/local/go/bin:${PATH}"

# TODO Add renovate tracking
FROM --platform=${BUILDPLATFORM} docker.io/alpine:3.21.3 AS dev

# Docker buildx ARGs
ARG TARGETOS
ARG TARGETARCH
# Disable C compatability, pass build arg ENABLE_CGO=1 to enable
ARG ENABLE_CGO=0

ENV CGO_ENABLED=${ENABLE_CGO} GOOS=${TARGETOS} GOARCH=${TARGETARCH}

# Use manually mounted volume for dev, providing hot-reload possibility
WORKDIR /app

COPY --link --from=go-deps /usr/local/go /usr/local/go

# hadolint ignore=DL3018
RUN apk add --no-cache git sudo && adduser -D -u 1000 nonroot && echo "nonroot ALL=(ALL) NOPASSWD: ALL" >/etc/sudoers.d/nonroot && chmod 0440 /etc/sudoers.d/nonroot && chown nonroot:nonroot /usr/local/go

ENV PATH="/home/nonroot/go/bin:/usr/local/go/bin:${PATH}"

USER 1000:1000

RUN go install github.com/air-verse/air@latest

ENTRYPOINT ["air", "-c", "config/air/.air.toml"]

FROM --platform=${BUILDPLATFORM} clean-deps AS base

RUN --mount=type=cache,target=/root/go/pkg \
	--mount=type=bind,source=go.mod,target=go.mod \
	--mount=type=bind,source=go.sum,target=go.sum \
	go mod download -x

FROM base AS build

WORKDIR /app

RUN --mount=type=cache,target=/root/go/pkg \
	--mount=type=cache,target=/root/.cache \
	--mount=type=bind,target=. \
	go build -ldflags='-s -w' -o /usr/local/bin/app ./cmd/main

# hadolint ignore=DL3007
FROM gcr.io/distroless/static-debian12:nonroot AS main

COPY --link --from=build /usr/local/bin/app /app

CMD ["/app"]
