FROM golang:latest AS builder
ENV GOPROXY="https://repo.cci.nokia.net/proxy-golang-org"
ENV GOSUMDB=off
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download -x
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o webcalc -ldflags "-w -extldflags '-static'" .

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /app/webcalc .
EXPOSE 8080
EXPOSE 8081
USER 65535:65536
ENTRYPOINT ["/webcalc"]