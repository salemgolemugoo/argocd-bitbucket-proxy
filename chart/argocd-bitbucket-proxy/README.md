# argocd-bitbucket-proxy

Proxy that translates Bitbucket webhooks to GitHub format for ArgoCD ApplicationSet

**Homepage:** <https://github.com/salemgolemugoo/argocd-bitbucket-proxy>

## Overview

This chart deploys a lightweight proxy that translates Bitbucket Server and Bitbucket Cloud webhook payloads into GitHub webhook format, enabling ArgoCD ApplicationSet webhook support for Bitbucket.

This is a temporary workaround until [argoproj/argo-cd#15443](https://github.com/argoproj/argo-cd/pull/15443) is merged.

## Installation

```bash
helm install argocd-bitbucket-proxy chart/argocd-bitbucket-proxy/ \
  --set secrets.bitbucketServerSecret=your-secret \
  --namespace argocd
```

At least one of `secrets.bitbucketServerSecret` or `secrets.bitbucketCloudSecret` must be set.

## Supported Events

| Bitbucket Event | Translated To |
|---|---|
| Server `repo:refs_changed` | GitHub `push` |
| Server `pr:opened/merged/declined/deleted` | GitHub `pull_request` |
| Cloud `repo:push` | GitHub `push` |
| Cloud `pullrequest:created/updated/fulfilled/rejected` | GitHub `pull_request` |

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| salemgolemugoo |  |  |

## Source Code

* <https://github.com/salemgolemugoo/argocd-bitbucket-proxy>

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` | Affinity rules for pod scheduling |
| env.argocdWebhookURL | string | `"http://argocd-applicationset-controller.argocd.svc.cluster.local:7000/api/webhook"` | ArgoCD ApplicationSet webhook URL to forward translated payloads to |
| env.logLevel | string | `"info"` | Log level (debug, info, warn, error) |
| gateway | object | `{"create":true,"enabled":false,"gatewayClassName":"amazon-vpc-lattice","gatewayRef":{"name":"","namespace":""},"httpRoute":{"hostnames":["argocd-bitbucket-proxy.example.com"],"rules":[{"matches":[{"path":{"type":"PathPrefix","value":"/"}}]}]},"listeners":[{"name":"https","port":443,"protocol":"HTTPS","tls":{"certificateRefs":[{"name":"argocd-bitbucket-proxy-tls"}],"mode":"Terminate"}}]}` | Gateway API configuration for external access |
| gateway.create | bool | `true` | Create a Gateway resource (set to false to use an existing Gateway) |
| gateway.enabled | bool | `false` | Enable Gateway API resources |
| gateway.gatewayClassName | string | `"amazon-vpc-lattice"` | Gateway class name (used when create is true) |
| gateway.gatewayRef | object | `{"name":"","namespace":""}` | Reference to an existing Gateway (used when create is false) |
| gateway.gatewayRef.name | string | `""` | Name of the existing Gateway |
| gateway.gatewayRef.namespace | string | `""` | Namespace of the existing Gateway (defaults to release namespace) |
| gateway.httpRoute | object | `{"hostnames":["argocd-bitbucket-proxy.example.com"],"rules":[{"matches":[{"path":{"type":"PathPrefix","value":"/"}}]}]}` | HTTPRoute configuration |
| gateway.httpRoute.hostnames | list | `["argocd-bitbucket-proxy.example.com"]` | Hostnames for the HTTPRoute |
| gateway.httpRoute.rules | list | `[{"matches":[{"path":{"type":"PathPrefix","value":"/"}}]}]` | Routing rules |
| gateway.listeners | list | `[{"name":"https","port":443,"protocol":"HTTPS","tls":{"certificateRefs":[{"name":"argocd-bitbucket-proxy-tls"}],"mode":"Terminate"}}]` | Gateway listeners (used when create is true) |
| image.pullPolicy | string | `"IfNotPresent"` | Image pull policy |
| image.repository | string | `"ghcr.io/salemgolemugoo/argocd-bitbucket-proxy"` | Container image repository |
| image.tag | string | `"latest"` | Container image tag |
| ingress | object | `{"annotations":{"alb.ingress.kubernetes.io/listen-ports":"[{\"HTTPS\":443}]","alb.ingress.kubernetes.io/scheme":"internet-facing","alb.ingress.kubernetes.io/ssl-redirect":"443","alb.ingress.kubernetes.io/target-type":"ip"},"className":"alb","enabled":false,"hosts":[{"host":"argocd-bitbucket-proxy.example.com","paths":[{"path":"/","pathType":"Prefix"}]}],"tls":[]}` | Ingress configuration for external access |
| ingress.annotations | object | `{"alb.ingress.kubernetes.io/listen-ports":"[{\"HTTPS\":443}]","alb.ingress.kubernetes.io/scheme":"internet-facing","alb.ingress.kubernetes.io/ssl-redirect":"443","alb.ingress.kubernetes.io/target-type":"ip"}` | Annotations for the ingress resource |
| ingress.className | string | `"alb"` | Ingress class name |
| ingress.enabled | bool | `false` | Enable ingress resource |
| ingress.hosts | list | `[{"host":"argocd-bitbucket-proxy.example.com","paths":[{"path":"/","pathType":"Prefix"}]}]` | Ingress host configuration |
| ingress.tls | list | `[]` | TLS configuration |
| nodeSelector | object | `{}` | Node selector for pod scheduling |
| replicaCount | int | `1` | Number of replicas |
| resources.limits.cpu | string | `"200m"` | CPU limit |
| resources.limits.memory | string | `"128Mi"` | Memory limit |
| resources.requests.cpu | string | `"50m"` | CPU request |
| resources.requests.memory | string | `"64Mi"` | Memory request |
| secrets.bitbucketCloudSecret | string | `""` | HMAC secret for validating Bitbucket Cloud webhooks |
| secrets.bitbucketServerSecret | string | `""` | HMAC secret for validating Bitbucket Server webhooks |
| service.port | int | `8080` | Service port |
| service.type | string | `"ClusterIP"` | Kubernetes service type |
| tolerations | list | `[]` | Tolerations for pod scheduling |
