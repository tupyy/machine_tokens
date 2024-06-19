FROM docker.io/golang:1.22 AS build

WORKDIR /app

COPY mis/ .

ARG build_args
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -extldflags "static"' -o mis main.go


################
#   Run step   #
################
FROM gcr.io/distroless/base

COPY --from=build /app/mis /usr/bin/mis

ENTRYPOINT ["/usr/bin/mis"]
