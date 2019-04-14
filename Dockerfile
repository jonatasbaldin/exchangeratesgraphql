FROM golang:1.12
WORKDIR /app

COPY . .
RUN make build-static

CMD ["sh", "-c", "make scrape && make run"]