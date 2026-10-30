import { Editor } from '@monaco-editor/react';
import './monaco-loader';

// Loaded on demand by `components/code-editor`, so the bundled editor is only
// downloaded once an editor is actually rendered.
export default Editor;
