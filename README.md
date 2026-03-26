# argocd-bitbucket-proxy

[![CI](https://github.com/salemgolemugoo/argocd-bitbucket-proxy/actions/workflows/ci.yaml/badge.svg)](https://github.com/salemgolemugoo/argocd-bitbucket-proxy/actions/workflows/ci.yaml)
[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/argocd-bitbucket-proxy)](https://artifacthub.io/packages/search?repo=argocd-bitbucket-proxy)

A lightweight proxy that translates Bitbucket Server and Bitbucket Cloud webhook payloads into GitHub webhook format, enabling ArgoCD ApplicationSet webhook support for Bitbucket.

This is a temporary workaround until [argoproj/argo-cd#15443](https://github.com/argoproj/argo-cd/pull/15443) is merged.

## How It Works

```
Bitbucket (Server/Cloud) → argocd-bitbucket-proxy → ArgoCD ApplicationSet Controller
                           (translates to GitHub format)
```

The proxy:
1. Receives Bitbucket webhooks on `POST /webhook`
2. Validates HMAC-SHA256 signatures
3. Translates push events → GitHub `push` format (for git generators)
4. Translates PR events → GitHub `pull_request` format (for PR generators)
5. Forwards to ArgoCD's ApplicationSet webhook endpoint

## Configuration

| Variable | Required | Default | Description |
|---|---|---|---|
| `ARGOCD_WEBHOOK_URL` | No | `http://argocd-applicationset-controller.argocd.svc.cluster.local:7000/api/webhook` | ArgoCD target URL |
| `BITBUCKET_SERVER_SECRET` | No* | | HMAC secret for Bitbucket Server |
| `BITBUCKET_CLOUD_SECRET` | No* | | HMAC secret for Bitbucket Cloud |
| `PORT` | No | `8080` | Listen port |
| `LOG_LEVEL` | No | `info` | Log level (debug/info/warn/error) |

*At least one secret must be configured.

## Deployment

### Helm

```bash
helm install argocd-bitbucket-proxy chart/argocd-bitbucket-proxy/ \
  --set secrets.bitbucketServerSecret=your-secret \
  --namespace argocd
```

### Docker

```bash
docker run -p 8080:8080 \
  -e BITBUCKET_SERVER_SECRET=your-secret \
  ghcr.io/salemgolemugoo/argocd-bitbucket-proxy:latest
```

## Supported Events

| Bitbucket Event | Translated To |
|---|---|
| Server `repo:refs_changed` | GitHub `push` |
| Server `pr:opened/merged/declined/deleted` | GitHub `pull_request` |
| Cloud `repo:push` | GitHub `push` |
| Cloud `pullrequest:created/updated/fulfilled/rejected` | GitHub `pull_request` |
