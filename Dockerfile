from golang:1.18-alpine

WORKDIR /app
COPY ./dist/linux-amd64 ./
CMD ["./puppet-github-metrics"]
