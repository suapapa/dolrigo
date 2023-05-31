# build stage
FROM golang:1.19 as builder

ARG PROGRAM_VER=dev-docker
ENV CGO_ENABLED=0

COPY . /build
WORKDIR /build

RUN go build -ldflags "-X main.programVer=${PROGRAM_VER}" -o /build/app

# ---
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/app .

# Diag http port
EXPOSE 8080

ENTRYPOINT ["./app"]