import { useState } from 'react';
import { useNavigate } from 'react-router';
import { useQueryFeature } from '@queries/feature-details';
import { getCurrentEnvironment, useAuth } from 'auth';
import { urls } from 'configs';
import {
  PAGE_PATH_FEATURE_TARGETING,
  PAGE_PATH_FEATURES
} from 'constants/routing';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import { Highlight, themes } from 'prism-react-renderer';
import { copyToClipBoard } from 'utils/function';
import { cn } from 'utils/style';
import { IconCopy, IconSwitch } from '@icons';
import Button from 'components/button';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';
import { renderSnippet, SDK_DEFINITIONS, SnippetVars } from './snippets';

// Languages not bundled with prism-react-renderer fall back to a close one.
const PRISM_LANGUAGE_FALLBACK: Record<string, string> = {
  dart: 'clike',
  groovy: 'clike'
};

const CodeBlock = ({
  code,
  language,
  onCopy
}: {
  code: string;
  language: string;
  onCopy: () => void;
}) => (
  <div className="relative rounded-lg border border-gray-200 overflow-hidden">
    <button
      type="button"
      onClick={onCopy}
      className="absolute right-2 top-2 flex-center rounded-md p-1.5 text-gray-500 hover:bg-gray-200 hover:text-gray-700"
    >
      <Icon icon={IconCopy} size="sm" />
    </button>
    <Highlight
      theme={themes.github}
      code={code}
      language={PRISM_LANGUAGE_FALLBACK[language] ?? language}
    >
      {({ style, tokens, getLineProps, getTokenProps }) => (
        <pre
          className="typo-para-small overflow-x-auto p-3 pr-10 font-fira-code"
          style={style}
        >
          {tokens.map((line, lineIndex) => (
            <div key={lineIndex} {...getLineProps({ line })}>
              {line.map((token, tokenIndex) => (
                <span key={tokenIndex} {...getTokenProps({ token })} />
              ))}
            </div>
          ))}
        </pre>
      )}
    </Highlight>
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
  const navigate = useNavigate();
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
  const installSnippet = renderSnippet(selectedSDK.install, snippetVars);
  const codeSnippet = renderSnippet(selectedSDK.code, snippetVars);

  const handleCopy = (text: string) => {
    copyToClipBoard(text);
    notify({ message: t('message:copied') });
  };

  const handleGoToFlag = () => {
    onClose();
    navigate(
      `/${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${flagId}${PAGE_PATH_FEATURE_TARGETING}`
    );
  };

  return (
    <DialogModal
      title={
        <div className="flex items-center gap-x-2">
          <Icon icon={IconSwitch} size="sm" />
          {t('walkthrough.connect-sdk.title')}
        </div>
      }
      isOpen={isOpen}
      onClose={onClose}
      closeOnClickOutside={false}
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
            code={installSnippet}
            language={selectedSDK.installLanguage}
            onCopy={() => handleCopy(installSnippet)}
          />
        </div>
        <div className="flex flex-col gap-y-2">
          <p className="typo-para-small text-gray-600">
            {t('walkthrough.connect-sdk.initialize')}
          </p>
          <CodeBlock
            code={codeSnippet}
            language={selectedSDK.codeLanguage}
            onCopy={() => handleCopy(codeSnippet)}
          />
        </div>
        <p
          className="typo-para-small rounded-lg border border-amber-200 bg-amber-50 px-4 py-3 text-amber-900"
          role="alert"
        >
          {t('walkthrough.connect-sdk.enable-note')}
        </p>
        <div className="flex justify-end">
          <Button type="button" variant="primary" onClick={handleGoToFlag}>
            {t('walkthrough.connect-sdk.go-to-flag')}
          </Button>
        </div>
      </div>
    </DialogModal>
  );
};

export default ConnectSdkModal;
