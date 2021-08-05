FROM alpine:3.14 AS cliDownloader

ARG DOCTL_VERSION=1.62.0
ARG KUBECTL_VERSION=1.21.3
ARG TERRAFORM_VERSION=1.0.3

WORKDIR /tools/

RUN apk add git

# Download kubectx
RUN git clone https://github.com/ahmetb/kubectx opt/kubectx && \
    mv opt/kubectx/kubectx ./kubectx && chmod +x ./kubectx && \
    mv opt/kubectx/kubens ./kubens && chmod +x ./kubens && \
    rm -rf opt

# Download doctl
RUN wget -qO- "https://github.com/digitalocean/doctl/releases/download/v${DOCTL_VERSION}/doctl-${DOCTL_VERSION}-linux-amd64.tar.gz" | tar -xzv

# Download kubectl
RUN wget "https://dl.k8s.io/release/v${KUBECTL_VERSION}/bin/linux/amd64/kubectl" && \
    wget "https://dl.k8s.io/v${KUBECTL_VERSION}/bin/linux/amd64/kubectl.sha256" && \ 
    echo "$(cat kubectl.sha256)  kubectl" | sha256sum -c && \
    chmod +x kubectl && \
    rm kubectl.sha256

# Download terraform
RUN wget "https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip" && \
    wget "https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_SHA256SUMS" && \
    grep "linux_amd64.zip" terraform_${TERRAFORM_VERSION}_SHA256SUMS  | sha256sum -c && \
    unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    rm terraform_${TERRAFORM_VERSION}_SHA256SUMS terraform_${TERRAFORM_VERSION}_linux_amd64.zip

# Download kubebuilder
RUN wget -O kubebuilder "https://go.kubebuilder.io/dl/latest/linux/amd64" && \
    chmod +x kubebuilder


# ## BUILD FRONT END

# FROM node:14.16-alpine as fe-builder

# WORKDIR /frontend/

# COPY slides-ts/ /frontend/

# RUN yarn install --frozen-lockfile && \
#     yarn run build

## BUILD SERVER

FROM golang:1.16-buster as be-builder

WORKDIR /backend/

COPY server/ /backend/

RUN GOOS=linux GOARCH=amd64 go build -o server main.go

## FINAL

FROM golang:1.16-buster

RUN apt update && apt install -y make

COPY --from=cliDownloader /tools/ /usr/bin/

# RUN adduser -D nonroot
RUN adduser --disabled-password --gecos '' nonroot
USER nonroot

WORKDIR /home/nonroot/presentation

# COPY --from=fe-builder /frontend/build/ /home/nonroot/presentation/static/
COPY slides-ts/build/ /home/nonroot/presentation/static/
COPY --from=be-builder /backend/server /home/nonroot/presentation/server

ENTRYPOINT [ "./server", "--address", ":8082", "--folder", "./static" ]
