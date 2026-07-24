import { useState } from 'react';
import { useQueryFeature } from '@queries/feature-details';
import { getCurrentEnvironment, useAuth } from 'auth';
import { urls } from 'configs';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import { copyToClipBoard } from 'utils/function';
import { cn } from 'utils/style';
import { IconCopy } from '@icons';
import Button from 'components/button';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';
import { SDK_DEFINITIONS, SnippetVars } from './snippets';

const CodeBlock = ({ code, onCopy }: { code: string; onCopy: () => void }) => (
  <div className="relative rounded-lg border border-gray-200 bg-gray-100">
    <button
      type="button"
      onClick={onCopy}
      className="absolute right-2 top-2 flex-center rounded-md p-1.5 text-gray-500 hover:bg-gray-200 hover:text-gray-700"
    >
      <Icon icon={IconCopy} size="sm" />
    </button>
    <pre className="typo-para-small overflow-x-auto p-3 pr-10 font-mono text-gray-900">
      {code}
    </pre>
  </div>
);

export type ConnectSdkModalProps = {
  isOpen: boolean;
  flagId: string;
  onClose: () => void;
};

const ConnectSdkModal = ({ isOpen, flagId, onClose }: ConnectSdkModalProps) => {
  const { t } = useTranslation(['common', 'message']);
  const { notify } = useToast();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const [selectedSDKId, setSelectedSDKId] = useState(SDK_DEFINITIONS[0].id);
  const selectedSDK =
    SDK_DEFINITIONS.find(sdk => sdk.id === selectedSDKId) || SDK_DEFINITIONS[0];

  const { data: featureCollection } = useQueryFeature({
    params: {
      id: flagId,
      environmentId: currentEnvironment.id
    },
    enabled: isOpen && !!flagId
  });

  const snippetVars: SnippetVars = {
    apiKey: 'YOUR_API_KEY',
    apiEndpoint: urls.API_ENDPOINT || 'https://your-api-endpoint',
    flagId,
    featureTag: featureCollection?.feature?.tags?.[0] || 'YOUR_FEATURE_TAG'
  };

  const handleCopy = (text: string) => {
    copyToClipBoard(text);
    notify({ message: t('message:copied') });
  };

  return (
    <DialogModal
      title={t('walkthrough.connect-sdk.title')}
      isOpen={isOpen}
      onClose={onClose}
      className="w-full max-w-2xl"
    >
      <div className="flex flex-col gap-y-4 p-5 pb-6">
        <p className="typo-para-medium text-gray-700">
          {t('walkthrough.connect-sdk.description')}
        </p>
        <div className="flex flex-wrap gap-2">
          {SDK_DEFINITIONS.map(sdk => (
            <button
              key={sdk.id}
              type="button"
              onClick={() => setSelectedSDKId(sdk.id)}
              className={cn(
                'typo-para-small rounded-full border px-3 py-1.5',
                sdk.id === selectedSDK.id
                  ? 'border-primary-500 bg-primary-50 text-primary-500'
                  : 'border-gray-300 text-gray-600 hover:border-gray-500'
              )}
            >
              {sdk.label}
            </button>
          ))}
        </div>
        <div className="flex flex-col gap-y-2">
          <p className="typo-para-small text-gray-600">
            {t('walkthrough.connect-sdk.install')}
          </p>
          <CodeBlock
            code={selectedSDK.install(snippetVars)}
            onCopy={() => handleCopy(selectedSDK.install(snippetVars))}
          />
        </div>
        <div className="flex flex-col gap-y-2">
          <p className="typo-para-small text-gray-600">
            {t('walkthrough.connect-sdk.initialize')}
          </p>
          <CodeBlock
            code={selectedSDK.code(snippetVars)}
            onCopy={() => handleCopy(selectedSDK.code(snippetVars))}
          />
        </div>
        <div className="flex justify-end">
          <Button type="button" variant="primary" onClick={onClose}>
            {t('walkthrough.done')}
          </Button>
        </div>
      </div>
    </DialogModal>
  );
};

export default ConnectSdkModal;
