FROM golang:latest AS builder
ENV GOPROXY="direct"
WORKDIR /app
COPY go.mod go.sum ./
RUN git config --global http.proxy http://135.245.48.34:8000
RUN http_proxy=http://135.245.48.34:8000/ https_proxy=http://135.245.48.34:8000/ go mod download -x
COPY . .
RUN GOOS=linux go build -a -o webcalc -ldflags "-X 'main.binaryType=static' -w -extldflags '-static'" .
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /app/webcalc .
USER 65532:65532
EXPOSE 8080:8080
EXPOSE 8081:8081
ENTRYPOINT ["/webcalc"]