FROM golang:1.18-alpine as builder
COPY go.mod go.sum /app/
WORKDIR /app/
RUN go mod download
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o pdfrender

FROM chromedp/headless-shell:latest
COPY --from=builder /app/pdfrender ./
ENTRYPOINT ["/pdfrender"]