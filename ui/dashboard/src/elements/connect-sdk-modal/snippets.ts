import { load } from 'js-yaml';
import androidYaml from './snippets/android.yaml?raw';
import flutterYaml from './snippets/flutter.yaml?raw';
import goYaml from './snippets/go.yaml?raw';
import iosYaml from './snippets/ios.yaml?raw';
import javascriptYaml from './snippets/javascript.yaml?raw';
import nodeYaml from './snippets/node.yaml?raw';
import reactNativeYaml from './snippets/react-native.yaml?raw';
import reactYaml from './snippets/react.yaml?raw';

export type SDKKind = 'client' | 'server';

export interface SnippetVars {
  apiKey: string;
  apiEndpoint: string;
  flagId: string;
  featureTag: string;
}

export interface SDKDefinition {
  id: string;
  label: string;
  kind: SDKKind;
  installLanguage: string;
  codeLanguage: string;
  install: string;
  code: string;
}

export const SDK_DEFINITIONS: SDKDefinition[] = [
  javascriptYaml,
  reactYaml,
  reactNativeYaml,
  androidYaml,
  iosYaml,
  flutterYaml,
  goYaml,
  nodeYaml
].map(raw => load(raw) as SDKDefinition);

export const renderSnippet = (template: string, vars: SnippetVars): string =>
  template.replace(
    /\{\{(\w+)\}\}/g,
    (match, key: string) => vars[key as keyof SnippetVars] ?? match
  );
