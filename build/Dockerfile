FROM quay.io/app-sre/golangci-lint:latest AS builder
ADD . /go/src/github.com/openshift/compliance-audit-router
WORKDIR /go/src/github.com/openshift/compliance-audit-router
RUN make build

FROM registry.access.redhat.com/ubi8-micro:latest
COPY --from=builder /go/src/github.com/openshift/compliance-audit-router/build/_output/compliance-audit-router /bin/compliance-audit-router
ENTRYPOINT ["/bin/compliance-audit-router"]
