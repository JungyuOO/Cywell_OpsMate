import { diagnosticsEntry } from './diagnostics-view.js';

export const pluginName = 'cyops-console';

export const cyopsPluginManifest = {
  name: pluginName,
  version: '0.0.44',
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
  const markEntryLoaded = () => {
    document.documentElement.setAttribute('data-cyops-plugin-entry', '0.0.44');
  };
  const cyopsLauncherFlag = () => {
    markEntryLoaded();
    mountUI();
    return Promise.resolve({ CYOPS_CONSOLE_LAUNCHER: true });
  };
  loadPluginEntry(`${pluginName}@0.0.44`, {
    cyopsLauncherFlag: () => Promise.resolve(() => cyopsLauncherFlag),
  });
  markEntryLoaded();
  mountUI();
  window.setTimeout(mountUI, 250);
  window.setTimeout(mountUI, 1000);
}
