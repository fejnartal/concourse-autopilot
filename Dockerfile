FROM alpine AS autopilot-cli
ARG RELEASE_TAG
RUN apk add --no-cache curl tar
RUN curl -o /tmp/autopilot.tar.gz -L https://github.com/efejjota/concourse-autopilot/releases/download/${RELEASE_TAG}/autopilot-linux-amd64.tar.gz \
 && tar -xvf /tmp/autopilot.tar.gz -C /tmp/ \
 && chmod +x /tmp/autopilot

FROM scratch
COPY --from=autopilot-cli /tmp/autopilot /go/bin/autopilot
