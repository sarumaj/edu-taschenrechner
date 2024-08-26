#!/bin/sh

set -e

CONFIG_FILE="/etc/nginx/nginx.conf"
PORT=${PORT:-8080}
TARGET_DIR="/usr/share/nginx/html"

cat <<EOF >"$CONFIG_FILE"
worker_processes 1;

events {
  worker_connections 1024;
}

http {
  include       mime.types;
  default_type  application/octet-stream;

  # Enable gzip compression
  gzip on;
  gzip_types text/plain text/css text/javascript application/javascript application/json application/wasm image/png image/gif;
  gzip_proxied any;
  gzip_min_length 256;

  # Limit the number of requests per second
  # allow 30 requests per second per IP
  # use 10 MB of memory to store the state (IP addresses)
  limit_req_zone \$binary_remote_addr zone=one:10m rate=30r/s;

  server {
    listen $PORT;

    server_name taschenrechner.sarumaj.com;

    location / {
      # Limit the number of requests per second
      # allow 20 burst requests per second per IP
      limit_req zone=one burst=20 nodelay;
      root $TARGET_DIR;
      gzip on;
    }
  }
}
EOF

echo "Validating Nginx configuration"
nginx -t -c "$CONFIG_FILE"
if [ $? -ne 0 ]; then
  echo "Nginx configuration is invalid. Exiting."
  exit 1
fi

# Start Nginx
nginx -g "daemon off;"
