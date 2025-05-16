import { useState } from 'react';
import { Trans } from 'react-i18next';
import { Link } from 'react-router-dom';
import { useTranslation } from 'i18n';
import { CodeReference } from '@types';
import { cn } from 'utils/style';
import { IconChevronDown } from '@icons';
import Button from 'components/button';
import Icon from 'components/icon';
import Card from 'elements/card';
import { repositoryTypeMap } from '..';
import CodeHighlighter from './highlighter';

const countOccurrences = (str: string, subStr: string) => {
  return (str.match(new RegExp(subStr, 'g')) || []).length;
};

const CodeAccordion = ({
  featureId,
  codeRef
}: {
  featureId: string;
  codeRef: CodeReference;
}) => {
  const { t } = useTranslation(['table']);
  const repository = repositoryTypeMap[codeRef.repositoryType];
  const occurrenceCount = countOccurrences(codeRef.codeSnippet, featureId);

  const [isOpen, setIsOpen] = useState(true);

  return (
    <Card className="p-5 gap-y-5">
      <div className="flex items-center divide-x divide-gray-400">
        <div className="flex items-center gap-x-3 pr-3">
          {repository ? <Icon icon={repository.icon} size="sm" /> : <></>}
          <p className="typo-para-medium text-gray-700">
            {repository?.label || codeRef.repositoryType}
          </p>
        </div>
        <div className="pl-3 typo-para-medium text-gray-500">
          <Trans
            i18nKey={`table:code-refs.${occurrenceCount > 1 ? 'multiple-refs' : 'single-ref'}`}
            components={{
              highlight: (
                <Link
                  to={codeRef.branchUrl}
                  className="underline text-primary-500"
                  target="_blank"
                >
                  {codeRef.repositoryBranch}
                </Link>
              )
            }}
            values={{
              count: occurrenceCount,
              branchLink: codeRef.repositoryBranch
            }}
          />
        </div>
      </div>
      <div className="flex flex-col w-full rounded-lg overflow-hidden">
        <div className="flex items-center justify-between w-full p-4 bg-gray-100">
          <p className="typo-para-medium text-gray-600">{codeRef.filePath}</p>
          <div className="flex items-center gap-x-4">
            <Link
              target="_blank"
              to={codeRef.sourceUrl}
              className="typo-para-medium text-primary-500 underline"
            >
              {t('code-refs.view-in-src')}
            </Link>
            <Button
              variant="grey"
              className={cn(
                'flex-center size-5 p-0 transition-all duration-300 rotate-0',
                {
                  'rotate-180': isOpen
                }
              )}
              onClick={() => setIsOpen(!isOpen)}
            >
              <Icon icon={IconChevronDown} size="sm" />
            </Button>
          </div>
        </div>
        <div
          className={cn(
            'overflow-x-auto transition-all duration-300 ease-in-out overflow-y-auto rounded-b-md max-h-0',
            { 'max-h-60 border-t border-gray-300': isOpen }
          )}
        >
          <CodeHighlighter featureId={featureId} codeRef={codeRef} />
        </div>
      </div>
    </Card>
  );
};

export default CodeAccordion;
