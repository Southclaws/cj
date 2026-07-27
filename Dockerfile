# -
# Build workspace
# -
FROM golang:1.25 AS compile

RUN apt-get update -y && apt-get install --no-install-recommends -y -q build-essential ca-certificates nodejs npm
RUN go install github.com/go-task/task/v3/cmd/task@latest

WORKDIR /cj
ADD . .
RUN task static

# -
# Runtime
# -
FROM scratch

COPY --from=compile /cj/cj /bin/cj
COPY --from=compile /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT ["cj"]
