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
  type: 'console.navigation/href',
  properties: {
    id: 'cyops-diagnostics',
    name: 'CYOps Diagnostics',
    href: '/console-plugin/diagnostics',
    section: 'cyops',
    perspective: 'admin',
  },
};
