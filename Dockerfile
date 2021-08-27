FROM alpine:latest as curl
RUN apk --update --no-cache add curl \
 && rm -rf /var/cache/apk/*

FROM curl as yq_cli
RUN curl -o yq -L https://github.com/mikefarah/yq/releases/download/v4.12.0/yq_linux_amd64 \
 && chmod +x yq

FROM curl as jq_cli
RUN curl -o jq -L https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64 \
 && chmod +x jq

FROM ubuntu

COPY --from=yq_cli yq /usr/local/bin/yq
COPY --from=jq_cli jq /usr/local/bin/jq

COPY scripts/gen-autopilot-extended.sh /usr/local/bin/autopilot
