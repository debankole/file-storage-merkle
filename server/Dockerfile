# syntax=docker/dockerfile:1

FROM golang:1.22 AS build-stage

WORKDIR /app

COPY server/go.mod server/go.sum common/go.mod ./
RUN go mod download

COPY ./ ./

WORKDIR /app/server

RUN CGO_ENABLED=0 GOOS=linux go build -o /server

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /app

COPY --from=build-stage /server server

EXPOSE 8080

ENTRYPOINT ["/app/server"]