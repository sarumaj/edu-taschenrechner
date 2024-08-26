FROM golang:latest AS builder

WORKDIR /usr/src/app

RUN go install fyne.io/fyne/v2/cmd/fyne@latest

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN version=$(git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//' | grep -Eo '^[0-9]+(\.[0-9]+){0,2}') && \
    [ -n "$version" ] || version="0.0.1"; \
    fyne package -os web -icon ./pkg/ui/icons/app.ico --appVersion "$version" --release --name "taschenrechner"

# production image
FROM nginx:alpine AS final

COPY ./scripts/entrypoint.sh /entrypoint.sh

COPY --from=builder /usr/src/app/wasm /usr/share/nginx/html

ENTRYPOINT ["/entrypoint.sh"]
