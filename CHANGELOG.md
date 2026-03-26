## [1.0.2](https://github.com/salemgolemugoo/argocd-bitbucket-proxy/compare/v1.0.1...v1.0.2) (2026-03-26)


### Bug Fixes

* **ci:** configure git identity for helm chart-releaser action ([#3](https://github.com/salemgolemugoo/argocd-bitbucket-proxy/issues/3)) ([66dcd62](https://github.com/salemgolemugoo/argocd-bitbucket-proxy/commit/66dcd625abc108b6133b1391bb773d3b6cdb4a5b))

## [1.0.1](https://github.com/salemgolemugoo/argocd-bitbucket-proxy/compare/v1.0.0...v1.0.1) (2026-03-26)


### Bug Fixes

* **asdf:** upgraded tool versions ([#1](https://github.com/salemgolemugoo/argocd-bitbucket-proxy/issues/1)) ([0b326f3](https://github.com/salemgolemugoo/argocd-bitbucket-proxy/commit/0b326f3ce5d11a3a2595925de417b22776a0790d))
* **helm:** changed chart license ([#2](https://github.com/salemgolemugoo/argocd-bitbucket-proxy/issues/2)) ([35523a1](https://github.com/salemgolemugoo/argocd-bitbucket-proxy/commit/35523a1701c9c7843375792cd64eb54ccbde1983))

# 1.0.0 (2026-03-26)


### Bug Fixes

* **ci:** update Go version to 1.26 to match go.mod ([26c0ad8](https://github.com/salemgolemugoo/argocd-bitbucket-proxy/commit/26c0ad8a6dd30c05b92fb34e0367b7a4332c2c96))


### Features

* add Bitbucket and GitHub webhook payload type definitions ([475c148](https://github.com/salemgolemugoo/argocd-bitbucket-proxy/commit/475c1481927a3d6bb415ef157d3cf834c15c6219))
* add Bitbucket Cloud push webhook translator ([323cb22](https://github.com/salemgolemugoo/argocd-bitbucket-proxy/commit/323cb22daa99d675e85913898bbe01ee3af526f0))
* add Bitbucket Server push webhook translator ([a40313a](https://github.com/salemgolemugoo/argocd-bitbucket-proxy/commit/a40313a1d01713e615d172a401cc73eab3917bea))
* add config package with env var loading and validation ([36dca5d](https://github.com/salemgolemugoo/argocd-bitbucket-proxy/commit/36dca5dc86e4f04bb94bf9d1ce4e7c811d30cdf3))
* add docker-compose for local testing ([6d27cfd](https://github.com/salemgolemugoo/argocd-bitbucket-proxy/commit/6d27cfd56dff95b839a8adf8ddd6b9d0edf478d8))
* add Helm chart for argocd-bitbucket-proxy ([acb80b5](https://github.com/salemgolemugoo/argocd-bitbucket-proxy/commit/acb80b52a54758ef6829a8d4e6c704d3709f9b84))
* add HMAC-SHA256 webhook signature validator ([9d79b34](https://github.com/salemgolemugoo/argocd-bitbucket-proxy/commit/9d79b34d860325e845310bef902fe6ffec11a048))
* add HTTP forwarder for proxying webhooks to ArgoCD ([4043799](https://github.com/salemgolemugoo/argocd-bitbucket-proxy/commit/4043799fab7a7a86d239847ece5f1cf64613ce92))
* add HTTP server with webhook handler, health checks, and routing ([2f406e5](https://github.com/salemgolemugoo/argocd-bitbucket-proxy/commit/2f406e5f7651b1b40171e93042ac7a9e37259f50))
* add multi-stage Dockerfile with distroless base ([b091688](https://github.com/salemgolemugoo/argocd-bitbucket-proxy/commit/b091688563f270328a811d8bad3c5e54a03fff11))
* add translator dispatcher for Bitbucket Server/Cloud events ([b40c91f](https://github.com/salemgolemugoo/argocd-bitbucket-proxy/commit/b40c91f4fa1b9d1d13afc6bc2a09ddeed56c8d97))
* enhance Helm chart with metadata, docs, and ArtifactHub support ([ce4c82c](https://github.com/salemgolemugoo/argocd-bitbucket-proxy/commit/ce4c82c0a3aca45611103955a64e38a1c206727c))
* wire up main.go with config, logging, and server ([8fb0e4e](https://github.com/salemgolemugoo/argocd-bitbucket-proxy/commit/8fb0e4e75a2f68f1554bf923bb407de5dbbcc57b))
