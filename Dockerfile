FROM golang:bullseye AS builder
COPY src $GOPATH/src/concourse-autopilot/

WORKDIR $GOPATH/src/concourse-autopilot/
RUN go get -d -v
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/autopilot

FROM scratch
COPY --from=builder /go/bin/autopilot /go/bin/autopilot
