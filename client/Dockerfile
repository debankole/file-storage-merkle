# syntax=docker/dockerfile:1

FROM golang:1.22 AS build-stage

WORKDIR /app

COPY client/go.mod client/go.sum common/go.mod ./
RUN go mod download

COPY ./ ./

WORKDIR /app/client

RUN CGO_ENABLED=0 GOOS=linux go build -o /client

# From golang:1.22 AS build-release-stage
FROM gcr.io/distroless/base-debian11 AS build-release-stage

ENV FILE_SERVER_URL=http://host.docker.internal:8080
WORKDIR /app

COPY --from=build-stage /client client
COPY --from=build-stage /app/client/files files

# EXPOSE 8080

ENTRYPOINT ["/app/client"]