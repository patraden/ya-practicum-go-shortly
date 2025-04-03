buf curl \
  --data '{"url": "http://ya.ru"}' \
  --schema . \
  --protocol grpc \
  --http2-prior-knowledge \
  --verbose \
  http://localhost:3200/shortener.v1.URLShortenerService/ShortenURL