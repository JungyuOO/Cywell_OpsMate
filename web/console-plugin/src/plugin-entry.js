import { diagnosticsEntry } from './diagnostics-view.js';

export const pluginName = 'cyops-console';

export const cyopsPluginManifest = {
  name: pluginName,
  version: '0.0.43',
  baseURL: '/api/plugins/cyops-console/',
  loadScripts: ['plugin-entry.js'],
  registrationMethod: 'callback',
  extensions: [
    diagnosticsEntry,
    {
      type: 'console.flag',
      properties: {
        handler: { $codeRef: 'cyopsLauncherFlag' },
      },
    },
  ],
};

export function registerCyopsPlugin(loadPluginEntry, mountUI) {
  const cyopsLauncherFlag = () => {
    mountUI();
    return Promise.resolve({ CYOPS_CONSOLE_LAUNCHER: true });
  };
  loadPluginEntry(`${pluginName}@0.0.43`, {
    cyopsLauncherFlag: () => Promise.resolve(() => cyopsLauncherFlag),
  });
}
