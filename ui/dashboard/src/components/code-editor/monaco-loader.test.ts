import { loader } from '@monaco-editor/react';
import * as monaco from 'monaco-editor/esm/vs/editor/editor.api.js';
import EditorWorker from 'monaco-editor/esm/vs/editor/editor.worker.js?worker';
import JsonWorker from 'monaco-editor/esm/vs/language/json/json.worker.js?worker';
import { describe, expect, it, vi } from 'vitest';
import { getWorker } from './monaco-loader';

// The real editor modules need a DOM; the wiring under test only needs
// stand-ins that can be told apart.
vi.mock('monaco-editor/esm/vs/editor/editor.api.js', () => ({
  editor: {},
  languages: {}
}));
vi.mock('monaco-editor/esm/vs/editor/edcore.main.js', () => ({}));
vi.mock(
  'monaco-editor/esm/vs/language/json/monaco.contribution.js',
  () => ({})
);
vi.mock(
  'monaco-editor/esm/vs/basic-languages/yaml/yaml.contribution.js',
  () => ({})
);
vi.mock('monaco-editor/esm/vs/editor/editor.worker.js?worker', () => ({
  default: class EditorWorkerStub {}
}));
vi.mock('monaco-editor/esm/vs/language/json/json.worker.js?worker', () => ({
  default: class JsonWorkerStub {}
}));

describe('monaco loader configuration', () => {
  it('resolves the bundled editor without touching the document', async () => {
    // Without a bundled instance the loader would have to append a <script>
    // tag pointing at its CDN, which is impossible here because there is no
    // document at all.
    expect(globalThis.document).toBeUndefined();

    await expect(loader.init()).resolves.toBe(monaco);
  });

  it('exposes the worker factory to the editor through MonacoEnvironment', () => {
    expect(globalThis.MonacoEnvironment?.getWorker).toBe(getWorker);
  });

  it('runs the JSON language service on the bundled JSON worker', () => {
    expect(getWorker('workerMain.js', 'json')).toBeInstanceOf(JsonWorker);
  });

  it('runs every other language on the bundled editor worker', () => {
    for (const label of ['editorWorkerService', 'yaml', 'plaintext']) {
      const worker = getWorker('workerMain.js', label);
      expect(worker).toBeInstanceOf(EditorWorker);
      expect(worker).not.toBeInstanceOf(JsonWorker);
    }
  });
});
