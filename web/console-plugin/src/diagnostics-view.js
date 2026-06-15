export async function loadDiagnostics(fetchJSON) {
  const [schema, diagnostics] = await Promise.all([
    fetchJSON('/api/ops/diagnostics/schema'),
    fetchJSON('/api/ops/diagnostics'),
  ]);
  return {
    schema,
    diagnostics,
    primaryEntry: diagnostics.ui?.primaryEntry || schema.primaryEntry,
  };
}

export const diagnosticsEntry = {
  id: 'cyops.diagnostics',
  title: 'CYOps Diagnostics',
  path: '/console-plugin/diagnostics',
  primaryEntry: 'openshift-web-console',
};
