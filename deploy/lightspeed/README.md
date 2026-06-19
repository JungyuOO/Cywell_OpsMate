# Lightspeed Boundary

CYOps does not call the internal LLM directly. The internal LLM URL and model
belong in OpenShift Lightspeed `OLSConfig`; CYOps points only at the Lightspeed
API endpoint through `OpsMateConfig.spec.lightspeed.apiBaseURL`.

Use `olsconfig-cywell-llm.yaml` as the CYWELL internal LLM profile:

- provider URL: `http://cllm.cywell.co.kr/v1`
- model: `gemma-4-26b-a4b-it-awq-8bit`
- provider type: `openai` for an OpenAI-compatible `/v1` endpoint

If the internal LLM endpoint does not require authentication, create the
referenced `cywell-cllm-credentials` Secret with an empty or harmless key only
if the installed Lightspeed Operator still requires `credentialsSecretRef`.

`OpsMateConfig` should then use the Lightspeed appserver/service URL, not the
internal LLM URL.
