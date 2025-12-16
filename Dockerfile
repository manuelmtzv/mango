
FROM node:22-alpine AS node-builder
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npx tailwindcss -i ./web/static/css/input.css -o ./web/static/css/output.css --minify

FROM golang:1.25-alpine AS go-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/server/main.go

FROM gcr.io/distroless/static-debian12
COPY --from=go-builder /app/main /
COPY --from=go-builder /app/web/templates /web/templates
COPY --from=go-builder /app/web/locales /web/locales

COPY --from=go-builder /app/web/static /web/static

COPY --from=node-builder /app/web/static/css/output.css /web/static/css/output.css

EXPOSE 8080
CMD ["/main"]