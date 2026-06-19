# v0.0.48 Phase 2 - Lightspeed OLSConfig Example

## Tasks

- [x] Add a Lightspeed `OLSConfig` example for the CYWELL internal LLM endpoint.
- [x] Document that CYOps must point at Lightspeed, not the internal LLM.
- [x] Update the `OpsMateConfig` sample to use a Lightspeed service-style endpoint.

## Verification

- Manifest review against local `lightspeed-operator` `OLSConfig` types:
  - `spec.llm.providers[].url`
  - `spec.llm.providers[].models[].name`
  - `spec.ols.defaultProvider`
  - `spec.ols.defaultModel`

## Remaining Scope

- Validate against an installed Lightspeed Operator on a real OCP cluster.
