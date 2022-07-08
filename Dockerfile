FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download -x
COPY . .
RUN GOOS=linux go build -a -o webcalc -ldflags "-X 'main.binaryType=static' -w -extldflags '-static'" .
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /app/webcalc .
USER 65532:65532
EXPOSE 8080
EXPOSE 8081
ENTRYPOINT ["/webcalc"]