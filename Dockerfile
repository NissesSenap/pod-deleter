
FROM golang:1.16-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# hadolint ignore=DL3059
RUN apk add --no-cache make=4.3-r0 bash=5.1.0-r0
RUN make build

FROM alpine:3.13.5 as runtime

ARG VERSION=1.1.0
ARG BUILD_DATE=2021-05-14

LABEL \
  org.opencontainers.image.created="$BUILD_DATE" \
  org.opencontainers.image.authors="edvin.norling@gmail.com" \
  org.opencontainers.image.homepage="https://github.com/NissesSenap/pod-deleter" \
  org.opencontainers.image.documentation="https://github.com/NissesSenap/pod-deleter" \
  org.opencontainers.image.source="https://github.com/NissesSenap/pod-deleter" \
  org.opencontainers.image.version="$VERSION" \
  org.opencontainers.image.vendor="N/A" \
  org.opencontainers.image.licenses="MIT" \
  summary="pod-deleter deletes pods in k8s normally trough falco rules from falcosidekick" \
  description="pod-deleter deletes pods in k8s normally trough falco rules from falcosidekick." \
  name="pod-deleter"

USER 1001

WORKDIR /app

COPY --from=builder /app/bin/poddeleter .

CMD ["./poddeleter"]
