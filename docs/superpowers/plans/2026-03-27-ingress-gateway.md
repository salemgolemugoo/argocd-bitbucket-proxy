# Ingress & Gateway API Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add optional Ingress and Gateway API (Gateway + HTTPRoute) resources to the Helm chart so users can expose the proxy externally.

**Architecture:** Three new templates (`ingress.yaml`, `gateway.yaml`, `httproute.yaml`) gated by `ingress.enabled` and `gateway.enabled` booleans. Values added to `values.yaml` with ALB controller defaults. Existing resources unchanged.

**Tech Stack:** Helm 3 templates, Kubernetes `networking.k8s.io/v1` Ingress, `gateway.networking.k8s.io/v1` Gateway/HTTPRoute

**Spec:** `docs/superpowers/specs/2026-03-27-ingress-gateway-design.md`

---

## File Structure

| File | Action | Responsibility |
|------|--------|----------------|
| `chart/argocd-bitbucket-proxy/values.yaml` | Modify | Add `ingress` and `gateway` value blocks |
| `chart/argocd-bitbucket-proxy/templates/ingress.yaml` | Create | Kubernetes Ingress resource |
| `chart/argocd-bitbucket-proxy/templates/gateway.yaml` | Create | Gateway API Gateway resource |
| `chart/argocd-bitbucket-proxy/templates/httproute.yaml` | Create | Gateway API HTTPRoute resource |

---

### Task 1: Add Ingress and Gateway values to values.yaml

**Files:**
- Modify: `chart/argocd-bitbucket-proxy/values.yaml:47` (append after `affinity: {}`)

- [ ] **Step 1: Add ingress and gateway blocks to values.yaml**

Append the following after line 47 (`affinity: {}`) in `chart/argocd-bitbucket-proxy/values.yaml`:

```yaml

# -- Ingress configuration for external access
ingress:
  # -- Enable ingress resource
  enabled: false
  # -- Ingress class name
  className: alb
  # -- Annotations for the ingress resource
  annotations:
    alb.ingress.kubernetes.io/scheme: internet-facing
    alb.ingress.kubernetes.io/target-type: ip
    alb.ingress.kubernetes.io/listen-ports: '[{"HTTPS":443}]'
    alb.ingress.kubernetes.io/ssl-redirect: "443"
  # -- Ingress host configuration
  hosts:
    - host: argocd-bitbucket-proxy.example.com
      paths:
        - path: /
          pathType: Prefix
  # -- TLS configuration
  tls: []
  #  - secretName: argocd-bitbucket-proxy-tls
  #    hosts:
  #      - argocd-bitbucket-proxy.example.com

# -- Gateway API configuration for external access
gateway:
  # -- Enable Gateway API resources
  enabled: false
  # -- Create a Gateway resource (set to false to use an existing Gateway)
  create: true
  # -- Gateway class name (used when create is true)
  gatewayClassName: amazon-vpc-lattice
  # -- Gateway listeners (used when create is true)
  listeners:
    - name: https
      protocol: HTTPS
      port: 443
      tls:
        mode: Terminate
        certificateRefs:
          - name: argocd-bitbucket-proxy-tls
  # -- Reference to an existing Gateway (used when create is false)
  gatewayRef:
    # -- Name of the existing Gateway
    name: ""
    # -- Namespace of the existing Gateway (defaults to release namespace)
    namespace: ""
  # -- HTTPRoute configuration
  httpRoute:
    # -- Hostnames for the HTTPRoute
    hostnames:
      - argocd-bitbucket-proxy.example.com
    # -- Routing rules
    rules:
      - matches:
          - path:
              type: PathPrefix
              value: /
```

- [ ] **Step 2: Commit**

```bash
git add chart/argocd-bitbucket-proxy/values.yaml
git commit -m "feat(chart): add ingress and gateway API values"
```

---

### Task 2: Create Ingress template

**Files:**
- Create: `chart/argocd-bitbucket-proxy/templates/ingress.yaml`

- [ ] **Step 1: Create the ingress template**

Create `chart/argocd-bitbucket-proxy/templates/ingress.yaml`:

```yaml
{{- if .Values.ingress.enabled -}}
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {{ include "argocd-bitbucket-proxy.fullname" . }}
  labels:
    {{- include "argocd-bitbucket-proxy.labels" . | nindent 4 }}
  {{- with .Values.ingress.annotations }}
  annotations:
    {{- toYaml . | nindent 4 }}
  {{- end }}
spec:
  {{- with .Values.ingress.className }}
  ingressClassName: {{ . }}
  {{- end }}
  {{- if .Values.ingress.tls }}
  tls:
    {{- range .Values.ingress.tls }}
    - hosts:
        {{- range .hosts }}
        - {{ . | quote }}
        {{- end }}
      secretName: {{ .secretName }}
    {{- end }}
  {{- end }}
  rules:
    {{- range .Values.ingress.hosts }}
    - host: {{ .host | quote }}
      http:
        paths:
          {{- range .paths }}
          - path: {{ .path }}
            pathType: {{ .pathType }}
            backend:
              service:
                name: {{ include "argocd-bitbucket-proxy.fullname" $ }}
                port:
                  number: {{ $.Values.service.port }}
          {{- end }}
    {{- end }}
{{- end }}
```

- [ ] **Step 2: Verify template renders with ingress enabled**

```bash
helm template test chart/argocd-bitbucket-proxy --set ingress.enabled=true
```

Expected: Output includes an `Ingress` resource with ALB annotations, host `argocd-bitbucket-proxy.example.com`, path `/`, and backend pointing to the service on port 8080.

- [ ] **Step 3: Verify template does NOT render with defaults**

```bash
helm template test chart/argocd-bitbucket-proxy
```

Expected: No `Ingress` resource in output. Only Deployment, Service (and Secret if secrets are set).

- [ ] **Step 4: Verify template renders with TLS**

```bash
helm template test chart/argocd-bitbucket-proxy \
  --set ingress.enabled=true \
  --set ingress.tls[0].secretName=my-tls \
  --set ingress.tls[0].hosts[0]=proxy.example.com
```

Expected: Ingress includes a `tls` block with `secretName: my-tls` and host `proxy.example.com`.

- [ ] **Step 5: Commit**

```bash
git add chart/argocd-bitbucket-proxy/templates/ingress.yaml
git commit -m "feat(chart): add ingress template"
```

---

### Task 3: Create Gateway template

**Files:**
- Create: `chart/argocd-bitbucket-proxy/templates/gateway.yaml`

- [ ] **Step 1: Create the gateway template**

Create `chart/argocd-bitbucket-proxy/templates/gateway.yaml`:

```yaml
{{- if and .Values.gateway.enabled .Values.gateway.create -}}
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: {{ include "argocd-bitbucket-proxy.fullname" . }}
  labels:
    {{- include "argocd-bitbucket-proxy.labels" . | nindent 4 }}
spec:
  gatewayClassName: {{ .Values.gateway.gatewayClassName }}
  listeners:
    {{- range .Values.gateway.listeners }}
    - name: {{ .name }}
      protocol: {{ .protocol }}
      port: {{ .port }}
      {{- with .tls }}
      tls:
        mode: {{ .mode }}
        {{- with .certificateRefs }}
        certificateRefs:
          {{- toYaml . | nindent 10 }}
        {{- end }}
      {{- end }}
    {{- end }}
{{- end }}
```

- [ ] **Step 2: Verify template renders with gateway enabled**

```bash
helm template test chart/argocd-bitbucket-proxy \
  --set gateway.enabled=true \
  --set gateway.create=true
```

Expected: Output includes a `Gateway` resource with `gatewayClassName: amazon-vpc-lattice` and an HTTPS listener on port 443.

- [ ] **Step 3: Verify template does NOT render when create is false**

```bash
helm template test chart/argocd-bitbucket-proxy \
  --set gateway.enabled=true \
  --set gateway.create=false
```

Expected: No `Gateway` resource in output.

- [ ] **Step 4: Verify template does NOT render with defaults**

```bash
helm template test chart/argocd-bitbucket-proxy
```

Expected: No `Gateway` resource in output.

- [ ] **Step 5: Commit**

```bash
git add chart/argocd-bitbucket-proxy/templates/gateway.yaml
git commit -m "feat(chart): add gateway template"
```

---

### Task 4: Create HTTPRoute template

**Files:**
- Create: `chart/argocd-bitbucket-proxy/templates/httproute.yaml`

- [ ] **Step 1: Create the httproute template**

Create `chart/argocd-bitbucket-proxy/templates/httproute.yaml`:

```yaml
{{- if .Values.gateway.enabled -}}
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: {{ include "argocd-bitbucket-proxy.fullname" . }}
  labels:
    {{- include "argocd-bitbucket-proxy.labels" . | nindent 4 }}
spec:
  parentRefs:
    {{- if .Values.gateway.create }}
    - name: {{ include "argocd-bitbucket-proxy.fullname" . }}
    {{- else }}
    - name: {{ .Values.gateway.gatewayRef.name }}
      {{- with .Values.gateway.gatewayRef.namespace }}
      namespace: {{ . }}
      {{- end }}
    {{- end }}
  {{- with .Values.gateway.httpRoute.hostnames }}
  hostnames:
    {{- toYaml . | nindent 4 }}
  {{- end }}
  rules:
    {{- range .Values.gateway.httpRoute.rules }}
    - matches:
        {{- toYaml .matches | nindent 8 }}
      backendRefs:
        - name: {{ include "argocd-bitbucket-proxy.fullname" $ }}
          port: {{ $.Values.service.port }}
    {{- end }}
{{- end }}
```

- [ ] **Step 2: Verify HTTPRoute renders with gateway enabled and create=true**

```bash
helm template test chart/argocd-bitbucket-proxy \
  --set gateway.enabled=true \
  --set gateway.create=true
```

Expected: Output includes an `HTTPRoute` with `parentRefs` pointing to the chart-created Gateway name, hostname `argocd-bitbucket-proxy.example.com`, PathPrefix `/` match, and backendRef to the service on port 8080.

- [ ] **Step 3: Verify HTTPRoute renders with external gateway ref**

```bash
helm template test chart/argocd-bitbucket-proxy \
  --set gateway.enabled=true \
  --set gateway.create=false \
  --set gateway.gatewayRef.name=shared-gateway \
  --set gateway.gatewayRef.namespace=infra
```

Expected: Output includes an `HTTPRoute` with `parentRefs` pointing to `name: shared-gateway`, `namespace: infra`. No `Gateway` resource in output.

- [ ] **Step 4: Verify HTTPRoute does NOT render with defaults**

```bash
helm template test chart/argocd-bitbucket-proxy
```

Expected: No `HTTPRoute` resource in output.

- [ ] **Step 5: Commit**

```bash
git add chart/argocd-bitbucket-proxy/templates/httproute.yaml
git commit -m "feat(chart): add httproute template"
```

---

### Task 5: Integration verification

- [ ] **Step 1: Verify all resources render together**

```bash
helm template test chart/argocd-bitbucket-proxy \
  --set ingress.enabled=true \
  --set gateway.enabled=true
```

Expected: Output contains Deployment, Service, Ingress, Gateway, and HTTPRoute — all with consistent labels and names. No template errors.

- [ ] **Step 2: Verify defaults produce no new resources**

```bash
helm template test chart/argocd-bitbucket-proxy 2>&1 | grep "^kind:"
```

Expected: Only `Deployment` and `Service` (no Secret since secrets are empty by default).

- [ ] **Step 3: Run helm lint**

```bash
helm lint chart/argocd-bitbucket-proxy
```

Expected: `0 chart(s) failed` — no errors or warnings.

- [ ] **Step 4: Run helm lint with all features enabled**

```bash
helm lint chart/argocd-bitbucket-proxy \
  --set ingress.enabled=true \
  --set gateway.enabled=true
```

Expected: `0 chart(s) failed`.

- [ ] **Step 5: Commit any fixes if needed, then final commit**

If no fixes needed, skip this step. Otherwise:

```bash
git add -A
git commit -m "fix(chart): address lint issues"
```
