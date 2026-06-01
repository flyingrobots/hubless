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
ENV GO_VERSION=1.26.3
RUN curl -fsSLo /tmp/go.tar.gz https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf /tmp/go.tar.gz \
    && rm /tmp/go.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOTOOLCHAIN=local
RUN go env -w GOTOOLCHAIN=local

# Markdown transclusion CLI
ENV MARKDOWN_TRANSCLUSION_VERSION=1.2.0
RUN npm install -g markdown-transclusion@${MARKDOWN_TRANSCLUSION_VERSION}

# golangci-lint binary
ENV GOLANGCI_LINT_VERSION=2.12.2
RUN curl -sSfL https://github.com/golangci/golangci-lint/releases/download/v${GOLANGCI_LINT_VERSION}/golangci-lint-${GOLANGCI_LINT_VERSION}-linux-amd64.tar.gz \
    | tar -xz -C /usr/local/bin --strip-components=1 golangci-lint-${GOLANGCI_LINT_VERSION}-linux-amd64/golangci-lint

WORKDIR /workspace

CMD ["/bin/bash", "/workspace/scripts/run-ci-inside.sh"]
