curl -X POST https://localhost:8443/ \
  --cacert deployments/.certs/shortener-cert.pem \
  -H "Content-Type: text/plain" \
  -d "https://google.com/"