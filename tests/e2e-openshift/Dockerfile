# The Dockerfile's resulting image is purpose-built for executing OpenTelemetry Operator e2e tests within the OpenShift release (https://github.com/openshift/release) using Prow CI. 

FROM golang:1.21

# Copy the repository files
COPY . /tmp/opentelemetry-operator

WORKDIR /tmp

# Set the Go path and Go cache environment variables
ENV GOPATH=/tmp/go
ENV GOBIN=/tmp/go/bin
ENV GOCACHE=/tmp/.cache/go-build
ENV PATH=$PATH:$GOBIN

# Create the /tmp/go/bin and build cache directories, and grant read and write permissions to all users
RUN mkdir -p /tmp/go/bin $GOCACHE \
    && chmod -R 777 /tmp/go/bin $GOPATH $GOCACHE

# Install dependencies required by test cases and debugging
RUN apt-get update && apt-get install -y jq vim libreadline-dev

# Install kuttl e2e
RUN curl -LO https://github.com/kudobuilder/kuttl/releases/download/v0.15.0/kubectl-kuttl_0.15.0_linux_x86_64 \
    && chmod +x kubectl-kuttl_0.15.0_linux_x86_64 \
    && mv kubectl-kuttl_0.15.0_linux_x86_64 /usr/local/bin/kuttl

# Install Chainsaw e2e
RUN go install github.com/kyverno/chainsaw@v0.1.6

# Install kubectl and oc
RUN curl -LO https://mirror.openshift.com/pub/openshift-v4/x86_64/clients/ocp/latest/openshift-client-linux.tar.gz \
    && tar -xzf openshift-client-linux.tar.gz \
    && chmod +x oc kubectl \
    && mv oc kubectl /usr/local/bin/

# Set the working directory
WORKDIR /tmp/opentelemetry-operator
