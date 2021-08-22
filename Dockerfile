FROM ubuntu

RUN apt-get update \
 && apt-get install -y --no-install-recommends \
      curl \
      jq \
      git \
      ca-certificates \
 && rm -rf /var/lib/apt/lists/* \
 && curl -o /usr/local/bin/yq -L https://github.com/mikefarah/yq/releases/download/v4.12.0/yq_linux_amd64 \
 && chmod +x /usr/local/bin/yq

COPY . /opt/resource
