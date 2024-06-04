ARG VARIANT="16"
FROM mcr.microsoft.com/devcontainers/go:1-1.21-bullseye

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends golanci-lint

# [Optional] Uncomment if you want to install an additional version
#  of node using nvm
# ARG EXTRA_NODE_VERSION=18
# RUN su node -c "source /usr/local/share/nvm/nvm.sh \
#    && nvm install ${EXTRA_NODE_VERSION}"

COPY ./set_container.sh /tmp/scripts/set_container.sh
RUN apt-get update && bash /tmp/scripts/set_container.sh
