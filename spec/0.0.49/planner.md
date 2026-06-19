# v0.0.49 Planner

## 1. Goal

- Complete CRC wiring between CYOps and OpenShift Lightspeed.
- Make CYOps send Lightspeed OLS-compatible query payloads to `/v1/query`.
- Record CRC-only Lightspeed Operator install and NetworkPolicy manifests.

## 2. Architecture Overview

- Lightspeed Operator owns `OLSConfig`, appserver, postgres, and console plugin.
- `OLSConfig` points to CYWELL internal LLM `http://cllm.cywell.co.kr/v1`.
- CYOps appserver calls only `https://lightspeed-app-server.openshift-lightspeed.svc:8443/v1/query`.
- CRC NetworkPolicy allows `cywell-opsmate-olm` pods to call Lightspeed appserver.

## 3. Technology Stack

| Area | Tooling | Notes |
| --- | --- | --- |
| Lightspeed | Red Hat `lightspeed-operator` v1.1.1 | Installed from `redhat-operators` on CRC |
| CYOps provider | Go HTTP provider | Sends `query` field, parses `response` |
| CRC network | Kubernetes NetworkPolicy | Allows CYOps namespace ingress to Lightspeed appserver |
| Packaging | OLM catalog | `v0.0.49 -> v0.0.48 -> v0.0.47` |

## 4. Implementation Steps

| Phase | Scope | Status | Evidence |
| --- | --- | --- | --- |
| Phase 1 | Lightspeed OLS request payload | done | provider test |
| Phase 2 | CRC Lightspeed manifests | done | live CRC apply |
| Phase 3 | Package, publish, upgrade, and smoke | planned | pending |

Tracking issue: #208

## 5. Migration Or Operations Strategy

- Upgrade CYOps from `cywell-opsmate.v0.0.48` to `cywell-opsmate.v0.0.49`.
- Install Lightspeed Operator first, then apply `OLSConfig`.
- Patch `OpsMateConfig.spec.lightspeed.apiBaseURL` to the Lightspeed `/v1/query` service URL after v0.0.49 is installed.

## 6. Message, Communication, And Data Protocol

- CYOps sends `{"query": "<message>"}` plus existing context fields to Lightspeed.
- Lightspeed returns `LLMResponse.response`, which CYOps parses through the existing response adapter.
- CYOps still does not send model/provider selection fields.

## 7. Security Considerations

- CRC uses a local NetworkPolicy exception from `cywell-opsmate-olm` to Lightspeed appserver.
- Real clusters should use the platform-approved namespace and policy boundary.
- Internal LLM credentials remain in the Lightspeed credentials Secret.

## 8. Completion Criteria

- [x] CYOps sends `query`, not `message`, to Lightspeed.
- [x] CRC Lightspeed Operator installs and `OLSConfig` becomes Ready.
- [x] CYOps namespace can reach Lightspeed appserver readiness.
- [x] Go tests pass.
- [x] OLM dry-run passes.
- [ ] v0.0.49 installs on CRC.
- [ ] CYOps `OpsMateConfig` points to Lightspeed `/v1/query`.
