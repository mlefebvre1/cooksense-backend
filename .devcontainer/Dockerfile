FROM golang:1.26-bookworm

ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID

RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m -s /bin/bash $USERNAME \
    && apt-get update \
    && apt-get install -y --no-install-recommends \
        sudo \
        git \
        make \
        curl \
        ca-certificates \
        postgresql-client \
        openssh-client \
        less \
        jq \
    && echo "$USERNAME ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/$USERNAME \
    && chmod 0440 /etc/sudoers.d/$USERNAME \
    && rm -rf /var/lib/apt/lists/*

# golangci-lint
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
        | sh -s -- -b /usr/local/bin v1.62.2

# modd (hot-reload used by `make run`)
RUN GOBIN=/usr/local/bin go install github.com/cortesi/modd/cmd/modd@latest

# Make caches owned by the dev user
RUN mkdir -p /go/pkg/mod /home/$USERNAME/.cache/go-build \
    && chown -R $USER_UID:$USER_GID /go /home/$USERNAME/.cache

USER $USERNAME
ENV GOPATH=/go \
    PATH=/go/bin:/usr/local/go/bin:$PATH
