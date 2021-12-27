# syntax=docker/dockerfile:1

ARG BASE_IMG=golang:1.17
ARG RUN_IMG=alpine:3.14

#############
# Base stage
#############
FROM ${BASE_IMG} as base

RUN DEBIAN_FRONTEND=noninteractive \
  apt-get update && apt-get install -y tzdata && \
  update-ca-certificates

# Add an unprivileged user
ENV USER=appuser
ENV UID=10001
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --no-create-home \ 
    --shell "/sbin/nologin" \        
    --uid "${UID}" \    
    "${USER}"


#############
# Compile stage
#############
FROM ${BASE_IMG} as compile

WORKDIR /go/src/github.com/massdriver-cloud/xo
COPY . .

RUN git config --global --add url."ssh://git@github.com/".insteadOf "https://github.com/"

# Fetch github's SSH host keys and compare them to the published
# ones at https://help.github.com/en/articles/githubs-ssh-key-fingerprints
RUN ["/bin/bash", "-c", "set -euo pipefail && mkdir -p -m 0600 ~/.ssh && ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts && ssh-keygen -F github.com -l -E sha256 | grep -q \"SHA256:nThbg6kXUpJWGl7E1IGOCspRomTxdCARLviKw6E5SY8\""]

RUN --mount=type=ssh GOPRIVATE=github.com/massdriver-cloud CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /usr/bin/xo .

ENTRYPOINT ["/usr/bin/xo"]


#############
# Run stage
#############
FROM ${RUN_IMG}

# Get tzdata
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo

# Get updated certs
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Use unprivileged user
COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /etc/group /etc/group
USER appuser:appuser

COPY --from=compile /usr/bin/xo /usr/bin/xo

ENTRYPOINT ["/usr/bin/xo"]
