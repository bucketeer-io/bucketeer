// Serves the Monaco editor from the application bundle instead of letting
// @monaco-editor/react download it from its default CDN at runtime.
//
// Only the editor core (with every editor feature) and the two languages the
// dashboard edits are bundled: JSON, with its language service, and YAML.
// Importing the package root instead would also bundle the CSS/HTML/TypeScript
// language services and emit their multi-megabyte web workers even though
// nothing ever requests them.
import { loader } from '@monaco-editor/react';
import 'monaco-editor/esm/vs/basic-languages/yaml/yaml.contribution.js';
import 'monaco-editor/esm/vs/editor/edcore.main.js';
import * as monaco from 'monaco-editor/esm/vs/editor/editor.api.js';
import editorWorker from 'monaco-editor/esm/vs/editor/editor.worker.js?worker';
import jsonWorker from 'monaco-editor/esm/vs/language/json/json.worker.js?worker';
import 'monaco-editor/esm/vs/language/json/monaco.contribution.js';

// Creates the bundled web worker for a Monaco language service.
// Only the JSON language service needs its own worker; every other
// language (e.g. YAML) runs on the base editor worker.
export const getWorker = (_workerId: string, label: string): Worker => {
  if (label === 'json') {
    return new jsonWorker();
  }
  return new editorWorker();
};

globalThis.MonacoEnvironment = { getWorker };

loader.config({ monaco });
