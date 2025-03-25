FROM golang:1.24.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG BUILD_VERSION
ARG BUILD_DATE
ARG BUILD_COMMIT
ARG VERSION_PACKAGE

RUN go build -ldflags="\
  -s -w \
  -X ${VERSION_PACKAGE}.buildVersion=${BUILD_VERSION} \
  -X ${VERSION_PACKAGE}.buildDate=${BUILD_DATE} \
  -X ${VERSION_PACKAGE}.buildCommit=${BUILD_COMMIT}" \
  -o /app/shortener ./cmd/shortener

FROM golang:1.24.1

WORKDIR /app

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN mkdir -p data

# Generate self-signed TLS certificate
# private key
RUN mkdir -p /etc/ssl/private /etc/ssl/certs /app/certs && \
    openssl genpkey -algorithm RSA \ 
    -out /etc/ssl/private/shortener-key.pem -pkeyopt rsa_keygen_bits:4096

# Ensure the correct permissions for private key
RUN chmod 600 /etc/ssl/private

# get self-signed certificate with pub key
RUN openssl req -x509 -new -key /etc/ssl/private/shortener-key.pem \
    -out /etc/ssl/certs/shortener-cert.pem -sha256 \
    -days 365 -nodes \
    -subj "/C=RU/L=Belgrade/O=PatrakhinDenis/OU=Dev/CN=localhost"

RUN cp /etc/ssl/certs/shortener-cert.pem /app/certs/

COPY --from=builder /app/shortener /usr/local/bin/shortener

CMD /usr/local/bin/shortener