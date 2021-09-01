FROM curlimages/curl AS builder
ARG RELEASE_TAG
RUN curl -o /tmp/autopilot -L https://github.com/efejjota/concourse-autopilot/releases/download/${RELEASE_TAG}/autopilot-linux-amd64 \
 && chmod +x /tmp/autopilot

FROM scratch
COPY --from=builder /tmp/autopilot /go/bin/autopilot
