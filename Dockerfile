# syntax=docker/dockerfile:1
ARG GO_VERSION=1.22.1
FROM golang:${GO_VERSION} AS build
WORKDIR /src

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    --mount=type=bind,source=Makefile,target=Makefile \
	make deps

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=Makefile,target=Makefile \
    --mount=type=bind,target=. \
	CGO_ENABLED=0 go build -o /bin/gymo .

FROM alpine:latest AS final

WORKDIR /app
RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add \
        ca-certificates \
        tzdata \
        && \
        update-ca-certificates

ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
USER appuser

COPY ./docs /app/docs
COPY --from=build /bin/gymo /app/

EXPOSE 4000

ENTRYPOINT [ "/app/gymo" ]
