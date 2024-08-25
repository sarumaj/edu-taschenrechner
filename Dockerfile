FROM golang:latest AS builder

WORKDIR /usr/src/app

RUN go install fyne.io/fyne/v2/cmd/fyne@latest

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN fyne package -os web -icon ./pkg/ui/icons/app.ico

# production image
FROM nginx:alpine AS final

COPY ./scripts/entrypoint.sh /entrypoint.sh

COPY --from=builder /usr/src/app/wasm /usr/share/nginx/html

ENTRYPOINT ["/entrypoint.sh"]
