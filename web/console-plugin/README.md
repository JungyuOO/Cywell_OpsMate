# CYOps Console Plugin Bundle Skeleton

This directory is the planned home for the OpenShift ConsolePlugin frontend
bundle.

v0.0.5 intentionally does not introduce a JavaScript toolchain. The first
frontend implementation should add the smallest buildable bundle that provides:

- CYOps floating launcher
- CYOps chat drawer shell
- chat composer
- document upload/list side panel placeholder
- calls to `/api/chat` and `/api/documents`

The bundle should keep the Red Hat OpenShift Lightspeed UI separate. CYOps owns
its own branding, drawer behavior, and document-management workflow.
