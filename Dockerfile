FROM node:24-alpine AS web
RUN corepack enable
WORKDIR /app
COPY ./web/package.json ./web/pnpm-lock.yaml ./web/
RUN cd web && pnpm install --frozen-lockfile
COPY ./web ./web
RUN cd web && pnpm build

FROM golang:1.26-alpine AS build
RUN apk add --no-cache gcc musl-dev
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
COPY --from=web /app/internal/infrastructure/web/dist ./internal/infrastructure/web/dist
RUN CGO_ENABLED=1 GOOS=linux go build -o /out/dummy-smtp ./cmd/dummy-smtp

FROM alpine:3.20
RUN adduser -D -H -u 10001 app && mkdir /data && chown app /data
COPY --from=build /out/dummy-smtp /usr/local/bin/dummy-smtp
USER app
WORKDIR /data
EXPOSE 1025 8025
ENV STORAGE=memory SMTP_ADDR=0.0.0.0:1025 HTTP_ADDR=0.0.0.0:8025
ENTRYPOINT ["dummy-smtp"]
