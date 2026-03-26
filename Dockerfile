FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /argocd-bitbucket-proxy ./cmd/proxy/

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /argocd-bitbucket-proxy /argocd-bitbucket-proxy

USER nonroot:nonroot

ENTRYPOINT ["/argocd-bitbucket-proxy"]
