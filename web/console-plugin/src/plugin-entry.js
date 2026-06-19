import { diagnosticsEntry } from './diagnostics-view.js';

export const pluginName = 'cyops-console';

export const cyopsPluginManifest = {
  name: pluginName,
  version: '0.0.40',
  baseURL: '/api/plugins/cyops-console/',
  loadScripts: ['plugin-entry.js'],
  registrationMethod: 'callback',
  extensions: [diagnosticsEntry],
};

export function registerCyopsPlugin(loadPluginEntry, mountUI) {
  loadPluginEntry(pluginName, {
    init: mountUI,
    get: () => Promise.reject(new Error('CYOps does not expose module federation modules in v0.0.40')),
  });
}
