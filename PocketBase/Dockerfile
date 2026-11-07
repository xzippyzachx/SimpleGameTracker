# Multi-stage build: compile the custom PocketBase binary, then run it on a minimal image.
# BUILDPLATFORM/TARGETOS/TARGETARCH are set by buildx/compose when --platform is given,
# so building linux/arm64 on an amd64 host cross-compiles instead of emulating.
FROM --platform=$BUILDPLATFORM golang:1.27-alpine AS build

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /pb .

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S pb && adduser -S pb -G pb

COPY --from=build /pb /pb

# Static frontend (optional): mount the built PWA dist/ at /pb_public
# and set PB_PUBLIC_DIR=/pb_public to serve it from PocketBase itself.
ENV PB_PUBLIC_DIR=

USER pb
EXPOSE 8090
VOLUME /pb_data

ENTRYPOINT ["/pb"]
CMD ["serve", "--http=0.0.0.0:8090"]