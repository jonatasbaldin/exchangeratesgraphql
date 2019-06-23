FROM golang:1.12 as builder
WORKDIR /app
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o exchangeratesgraphql

FROM scratch
WORKDIR /app
# Bundling trusted certs to make HTTP requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app/exchangeratesgraphql /app/
COPY --from=builder /app/templates /app/templates
COPY --from=builder /app/static /app/static
CMD ["/app/exchangeratesgraphql", "-serve"]