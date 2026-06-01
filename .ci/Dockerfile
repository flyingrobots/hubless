FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        ca-certificates \
        curl \
        git \
        build-essential \
        ripgrep \
        gnupg \
    && curl -fsSL https://deb.nodesource.com/setup_20.x | bash - \
    && apt-get install -y --no-install-recommends nodejs \
    && rm -rf /var/lib/apt/lists/*

# Install Go 1.26.3
ARG TARGETARCH
ENV GO_VERSION=1.26.3
RUN set -eu; \
    go_arch="${TARGETARCH:-$(dpkg --print-architecture)}"; \
    case "$go_arch" in \
        amd64|arm64) ;; \
        *) echo "unsupported Go architecture: $go_arch" >&2; exit 1 ;; \
    esac; \
    curl -fsSLo /tmp/go.tar.gz "https://go.dev/dl/go${GO_VERSION}.linux-${go_arch}.tar.gz" \
    && tar -C /usr/local -xzf /tmp/go.tar.gz \
    && rm /tmp/go.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOTOOLCHAIN=local
RUN go env -w GOTOOLCHAIN=local

# Markdown transclusion CLI
ENV MARKDOWN_TRANSCLUSION_VERSION=1.2.0
RUN npm install -g markdown-transclusion@${MARKDOWN_TRANSCLUSION_VERSION}

# golangci-lint built with the active Go toolchain
ENV GOLANGCI_LINT_VERSION=2.12.2
RUN GOBIN=/usr/local/bin go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v${GOLANGCI_LINT_VERSION}

WORKDIR /workspace

CMD ["/bin/bash", "/workspace/scripts/run-ci-inside.sh"]
