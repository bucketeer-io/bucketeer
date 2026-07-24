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
  install: (v: SnippetVars) => string;
  code: (v: SnippetVars) => string;
}

export const SDK_DEFINITIONS: SDKDefinition[] = [
  {
    id: 'javascript',
    label: 'JavaScript',
    kind: 'client',
    installLanguage: 'shell',
    codeLanguage: 'javascript',
    install: () => 'npm install @bucketeer/js-client-sdk',
    code: v => `import {
  defineBKTConfig,
  defineBKTUser,
  getBKTClient,
  initializeBKTClient
} from '@bucketeer/js-client-sdk';

const config = defineBKTConfig({
  apiKey: '${v.apiKey}',
  apiEndpoint: '${v.apiEndpoint}',
  featureTag: '${v.featureTag}',
  appVersion: '1.0.0'
});

const user = defineBKTUser({ id: 'USER_ID' });

await initializeBKTClient(config, user);
const client = getBKTClient();

const showNewFeature = client?.booleanVariation('${v.flagId}', false);
if (showNewFeature) {
  // Show the new feature
} else {
  // Run the default code
}`
  },
  {
    id: 'react',
    label: 'React',
    kind: 'client',
    installLanguage: 'shell',
    codeLanguage: 'javascript',
    install: () => 'npm install @bucketeer/react-client-sdk',
    code: v => `import {
  BucketeerProvider,
  defineBKTConfigForReact,
  defineBKTUser,
  getBKTClient,
  initializeBKTClient,
  useBooleanVariation
} from '@bucketeer/react-client-sdk';

const config = defineBKTConfigForReact({
  apiKey: '${v.apiKey}',
  apiEndpoint: '${v.apiEndpoint}',
  featureTag: '${v.featureTag}',
  appVersion: '1.0.0'
});

const user = defineBKTUser({ id: 'USER_ID' });

// Initialize once in your root component, then wrap your app:
// await initializeBKTClient(config, user);
// <BucketeerProvider client={getBKTClient()}>...</BucketeerProvider>

function MyComponent() {
  const showNewFeature = useBooleanVariation('${v.flagId}', false);
  return showNewFeature ? <NewFeature /> : <DefaultFeature />;
}`
  },
  {
    id: 'react-native',
    label: 'React Native',
    kind: 'client',
    installLanguage: 'shell',
    codeLanguage: 'javascript',
    install: () => 'npm install @bucketeer/react-native-client-sdk',
    code: v => `import {
  BucketeerProvider,
  defineBKTConfigForReactNative,
  defineBKTUser,
  getBKTClient,
  initializeBKTClient,
  useBooleanVariation
} from '@bucketeer/react-native-client-sdk';

const config = defineBKTConfigForReactNative({
  apiKey: '${v.apiKey}',
  apiEndpoint: '${v.apiEndpoint}',
  featureTag: '${v.featureTag}',
  appVersion: '1.0.0'
});

const user = defineBKTUser({ id: 'USER_ID' });

// Initialize once in your root component, then wrap your app:
// await initializeBKTClient(config, user);
// <BucketeerProvider client={getBKTClient()}>...</BucketeerProvider>

function MyScreen() {
  const showNewFeature = useBooleanVariation('${v.flagId}', false);
  return showNewFeature ? <NewFeatureScreen /> : <DefaultScreen />;
}`
  },
  {
    id: 'android',
    label: 'Android',
    kind: 'client',
    installLanguage: 'groovy',
    codeLanguage: 'kotlin',
    install: () => `dependencies {
  implementation 'io.bucketeer:android-client-sdk:LATEST_VERSION'
}`,
    code: v => `import io.bucketeer.sdk.android.*

val config = BKTConfig.builder()
  .apiKey("${v.apiKey}")
  .apiEndpoint("${v.apiEndpoint}")
  .featureTag("${v.featureTag}")
  .appVersion(BuildConfig.VERSION_NAME)
  .build()

val user = BKTUser.builder()
  .id("USER_ID")
  .build()

BKTClient.initialize(this.application, config, user)
val client = BKTClient.getInstance()

val showNewFeature = client.booleanVariation("${v.flagId}", false)
if (showNewFeature) {
  // Show the new feature
} else {
  // Run the default code
}`
  },
  {
    id: 'ios',
    label: 'iOS',
    kind: 'client',
    installLanguage: 'ruby',
    codeLanguage: 'swift',
    install: () => `# CocoaPods
pod 'Bucketeer', 'LATEST_VERSION'

# or Swift Package Manager
# .package(url: "https://github.com/bucketeer-io/ios-client-sdk.git", exact: "LATEST_VERSION")`,
    code: v => `import Bucketeer

do {
  let config = try BKTConfig.Builder()
    .with(apiKey: "${v.apiKey}")
    .with(apiEndpoint: "${v.apiEndpoint}")
    .with(featureTag: "${v.featureTag}")
    .with(appVersion: "1.0.0")
    .build()

  let user = try BKTUser.Builder()
    .with(id: "USER_ID")
    .build()

  try BKTClient.initialize(config: config, user: user)
} catch {
  // Handle the error
}

let client = try? BKTClient.shared
let showNewFeature = client?.boolVariation(featureId: "${v.flagId}", defaultValue: false) ?? false
if showNewFeature {
  // Show the new feature
} else {
  // Run the default code
}`
  },
  {
    id: 'flutter',
    label: 'Flutter',
    kind: 'client',
    installLanguage: 'yaml',
    codeLanguage: 'dart',
    install: () => `dependencies:
  bucketeer_flutter_client_sdk: LATEST_VERSION`,
    code: v => `import 'package:bucketeer_flutter_client_sdk/bucketeer_flutter_client_sdk.dart';

final config = BKTConfigBuilder()
  .apiKey("${v.apiKey}")
  .apiEndpoint("${v.apiEndpoint}")
  .featureTag("${v.featureTag}")
  .appVersion("1.0.0")
  .build();

final user = BKTUserBuilder()
  .id("USER_ID")
  .build();

await BKTClient.initialize(config: config, user: user);
final client = BKTClient.instance;

final showNewFeature = await client.boolVariation("${v.flagId}", false);
if (showNewFeature) {
  // Show the new feature
} else {
  // Run the default code
}`
  },
  {
    id: 'go',
    label: 'Go',
    kind: 'server',
    installLanguage: 'shell',
    codeLanguage: 'go',
    install: () => 'go get github.com/bucketeer-io/go-server-sdk',
    code: v => `import (
  "github.com/bucketeer-io/go-server-sdk/pkg/bucketeer"
  "github.com/bucketeer-io/go-server-sdk/pkg/bucketeer/user"
)

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

client, err := bucketeer.NewSDK(
  ctx,
  bucketeer.WithAPIKey("${v.apiKey}"),
  bucketeer.WithAPIEndpoint("${v.apiEndpoint}"),
  bucketeer.WithTag("${v.featureTag}"),
)
if err != nil {
  log.Fatalf("Failed to initialize the client: %v", err)
}

endUser := user.NewUser("END_USER_ID", nil)
showNewFeature := client.BoolVariation(ctx, endUser, "${v.flagId}", false)
if showNewFeature {
  // Show the new feature
} else {
  // Run the default code
}`
  },
  {
    id: 'node',
    label: 'Node.js',
    kind: 'server',
    installLanguage: 'shell',
    codeLanguage: 'javascript',
    install: () => 'npm install @bucketeer/node-server-sdk',
    code: v => `import { defineBKTConfig, initializeBKTClient } from '@bucketeer/node-server-sdk';

const config = defineBKTConfig({
  apiKey: '${v.apiKey}',
  apiEndpoint: '${v.apiEndpoint}',
  featureTag: '${v.featureTag}'
});

const client = initializeBKTClient(config);

const user = { id: 'END_USER_ID', data: {} };
const showNewFeature = await client.getBoolVariation(user, '${v.flagId}', false);
if (showNewFeature) {
  // Show the new feature
} else {
  // Run the default code
}`
  }
];
