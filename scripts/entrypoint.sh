#!/bin/sh

# Generate the nginx.conf file using the PORT environment variable
cat <<EOF | envsubst "PORT=${PORT:-8080}" >/etc/nginx/nginx.conf
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
            root /usr/share/nginx/html;
            gzip on;
        }
    }
}
EOF

cat /etc/nginx/nginx.conf

# Start Nginx
nginx -g "daemon off;"
