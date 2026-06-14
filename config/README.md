# Config Scaffold

This directory will follow the Kubebuilder/Operator SDK layout used by OpenShift Lightspeed Operator:

- `crd/`: generated CRDs for `OpsMateConfig`
- `manager/`: controller manager deployment patches
- `rbac/`: service account, roles, and bindings
- `samples/`: sample `OpsMateConfig` resources
- `default/`: kustomize entrypoint for local installation
- `manifests/`: OLM CSV generation inputs
- `scorecard/`: operator-sdk scorecard configuration
