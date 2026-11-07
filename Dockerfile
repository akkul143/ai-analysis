FROM golang:1.23-alpine AS builder
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /out/jobfinder ./cmd/jobfinder

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /out/jobfinder /usr/local/bin/jobfinder
ENTRYPOINT ["/usr/local/bin/jobfinder"]
