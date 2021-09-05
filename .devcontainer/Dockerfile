# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.183.0/containers/go/.devcontainer/base.Dockerfile

# [Choice] Go version: 1, 1.16, 1.15
ARG VARIANT="1.16"
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

RUN git config --global --add url."ssh://git@github.com/".insteadOf "https://github.com/"
RUN ["/bin/bash", "-c", "set -euo pipefail && mkdir -p -m 0600 ~/.ssh && ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts && ssh-keygen -F github.com -l -E sha256 | grep -q \"SHA256:nThbg6kXUpJWGl7E1IGOCspRomTxdCARLviKw6E5SY8\""]
ENV GOPRIVATE=github.com/massdriver-cloud

# [Optional] Uncomment this section to install additional OS packages.
# RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
#     && apt-get -y install --no-install-recommends 

# [Optional] Uncomment the next line to use go get to install anything else you need
# RUN go get -x <your-dependency-or-tool>

# [Optional] Uncomment this line to install global node packages.
# RUN su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install -g <your-package-here>" 2>&1