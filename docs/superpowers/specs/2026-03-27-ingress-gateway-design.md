# Ingress & Gateway API Support for argocd-bitbucket-proxy Helm Chart

## Summary

Add optional Ingress and Gateway API resources to the Helm chart to allow external traffic to reach the proxy. Both are independently toggleable via `ingress.enabled` and `gateway.enabled`, disabled by default to preserve existing behavior.

## Motivation

The chart currently exposes the proxy only via a `ClusterIP` Service on port 8080. Users who want external connectivity must create Ingress or Gateway resources manually. Adding first-class support in the chart simplifies deployment and follows standard Helm patterns.

## Design

### Approach

Flat toggle pattern (industry standard): top-level `ingress` and `gateway` blocks, each with an `enabled` boolean. No mutual exclusion — users may enable both if their cluster supports it.

### Ingress (`templates/ingress.yaml`)

- **API version:** `networking.k8s.io/v1`
- **Gated by:** `ingress.enabled`
- **Ingress class:** `ingress.className` (default: `alb`)
- **Default annotations** (ALB controller):
  - `alb.ingress.kubernetes.io/scheme: internet-facing`
  - `alb.ingress.kubernetes.io/target-type: ip`
  - `alb.ingress.kubernetes.io/listen-ports: '[{"HTTPS":443}]'`
  - `alb.ingress.kubernetes.io/ssl-redirect: "443"`
- **Hosts:** List of `{host, paths[{path, pathType}]}`. All paths route to the Service on port 8080.
- **TLS:** Optional list of `{secretName, hosts[]}`. Defaults to empty (opt-in).

### Gateway API

Two templates, both gated by `gateway.enabled`:

#### `templates/gateway.yaml`

- **API version:** `gateway.networking.k8s.io/v1`
- **Gated by:** `gateway.enabled && gateway.create`
- **Gateway class:** `gateway.gatewayClassName` (default: `amazon-vpc-lattice`)
- **Listeners:** Configurable list. Default: single HTTPS listener on port 443 with TLS termination.
  - TLS config includes `mode` and `certificateRefs`.

#### `templates/httproute.yaml`

- **API version:** `gateway.networking.k8s.io/v1`
- **Always created when:** `gateway.enabled`
- **Parent ref:** If `gateway.create` is true, references the chart-created Gateway. Otherwise, references the user-provided `gateway.gatewayRef.name` and `gateway.gatewayRef.namespace`. If `gatewayRef.namespace` is empty, it defaults to the release namespace.
- **Hostnames:** Configurable list (default: `argocd-bitbucket-proxy.example.com`).
- **Rules:** Configurable list of match rules. Default: single PathPrefix `/` rule. All rules route to the Service on port 8080.

### Values Structure

```yaml
# -- Ingress configuration
ingress:
  enabled: false
  className: alb
  annotations:
    alb.ingress.kubernetes.io/scheme: internet-facing
    alb.ingress.kubernetes.io/target-type: ip
    alb.ingress.kubernetes.io/listen-ports: '[{"HTTPS":443}]'
    alb.ingress.kubernetes.io/ssl-redirect: "443"
  hosts:
    - host: argocd-bitbucket-proxy.example.com
      paths:
        - path: /
          pathType: Prefix
  tls: []

# -- Gateway API configuration
gateway:
  enabled: false
  create: true
  gatewayClassName: amazon-vpc-lattice
  listeners:
    - name: https
      protocol: HTTPS
      port: 443
      tls:
        mode: Terminate
        certificateRefs:
          - name: argocd-bitbucket-proxy-tls
  gatewayRef:
    name: ""
    namespace: ""
  httpRoute:
    hostnames:
      - argocd-bitbucket-proxy.example.com
    rules:
      - matches:
          - path:
              type: PathPrefix
              value: /
```

### What Does NOT Change

- Deployment template — unchanged
- Service template — stays `ClusterIP` on port 8080 (both Ingress and Gateway work with ClusterIP)
- Secret template — unchanged
- Helper templates — reused as-is for labels and naming

### Template Helpers

Existing helpers (`fullname`, `labels`, `selectorLabels`) are reused. No new helpers needed.

## Files to Create/Modify

| File | Action |
|------|--------|
| `chart/argocd-bitbucket-proxy/templates/ingress.yaml` | Create |
| `chart/argocd-bitbucket-proxy/templates/gateway.yaml` | Create |
| `chart/argocd-bitbucket-proxy/templates/httproute.yaml` | Create |
| `chart/argocd-bitbucket-proxy/values.yaml` | Add `ingress` and `gateway` blocks |

## Testing

- `helm template` with defaults — should produce no Ingress/Gateway resources
- `helm template` with `ingress.enabled=true` — should produce Ingress resource with ALB annotations
- `helm template` with `gateway.enabled=true` — should produce Gateway + HTTPRoute
- `helm template` with `gateway.enabled=true,gateway.create=false` — should produce only HTTPRoute with external gatewayRef
- `helm template` with both enabled — should produce all resources without conflict
