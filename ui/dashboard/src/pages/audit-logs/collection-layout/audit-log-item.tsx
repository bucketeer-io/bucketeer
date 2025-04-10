import React, { useMemo, useState } from 'react';
import DiffViewer from 'react-diff-viewer-continued';
import primaryAvatar from 'assets/avatars/primary.svg';
import { useTranslation } from 'i18n';
import { AuditLog } from '@types';
import { formatLongDateTime } from 'utils/date-time';
import { IconChevronDown, IconCopy, IconWatch } from '@icons';
import { AvatarImage } from 'components/avatar';
import Button from 'components/button';
import Icon from 'components/icon';
import { Tabs, TabsContent, TabsList, TabsTrigger } from 'components/tabs';
import { AuditLogTab } from '../types';

const AuditLogItem = ({ auditLog }: { auditLog: AuditLog }) => {
  const {
    editor,
    localizedMessage,
    timestamp,
    options,
    entityData,
    previousEntityData
  } = auditLog;
  const { t } = useTranslation(['common']);
  const [currentTab, setCurrentTab] = useState<AuditLogTab>(
    AuditLogTab.CHANGES
  );

  const time = useMemo(
    () =>
      formatLongDateTime({
        value: timestamp,
        overrideOptions: {
          day: undefined,
          year: undefined,
          hour: '2-digit',
          minute: '2-digit',
          hour12: false
        }
      }),
    [timestamp]
  );

  return (
    <div className="flex flex-col w-full p-3 gap-y-5 bg-white shadow-card rounded-lg">
      <div className="flex items-center w-full gap-x-3">
        <AvatarImage
          image={primaryAvatar}
          alt="member-avatar"
          className="size-10"
        />
        <div className="flex flex-col flex-1 gap-y-1 overflow-hidden">
          <div className="flex items-center gap-x-1.5 max-w-full">
            <p className="typo-para-medium font-bold text-gray-700">
              {editor.name || editor.email}
            </p>
            <p className="typo-para-medium text-gray-700 truncate">
              {localizedMessage.message}
            </p>
          </div>
          <div className="flex items-center gap-x-1">
            <Icon icon={IconWatch} size={'sm'} />
            <p className="typo-para-small text-gray-500">{time}</p>
          </div>
        </div>
        <div className="flex items-center min-w-fit gap-x-5">
          <Button variant={'text'}>
            <Icon icon={IconCopy} color="primary-500" />
            Copy Link
          </Button>
          <div className="flex-center cursor-pointer">
            <Icon icon={IconChevronDown} />
          </div>
        </div>
      </div>
      {options.comment && (
        <div className="flex items-center w-full p-3 bg-gray-100 rounded typo-para-medium text-gray-600 break-all">
          {options.comment}
        </div>
      )}
      <Tabs
        className="flex-1 flex h-full flex-col"
        value={currentTab}
        onValueChange={value => setCurrentTab(value as AuditLogTab)}
      >
        <TabsList>
          <TabsTrigger value={AuditLogTab.CHANGES}>{t(`changes`)}</TabsTrigger>
          <TabsTrigger value={AuditLogTab.SNAPSHOT}>
            {t(`snapshot`)}
          </TabsTrigger>
        </TabsList>
        <TabsContent value={currentTab} className="flex flex-col gap-y-4">
          <p className="typo-para-small text-gray-500 uppercase">
            {t('patch')}
          </p>
          <DiffViewer
            oldValue={previousEntityData}
            newValue={entityData}
            splitView={false}
            showDiffOnly={true}
            hideMarkers
            hideLineNumbers
            useDarkTheme={false}
            renderGutter={({ lineNumber }) => {
              return (
                <p className="px-3 text-right typo-para-small text-gray-600 min-w-[38px]">
                  {lineNumber}
                </p>
              );
            }}
            renderContent={source => {
              if (!source.includes(':'))
                return (
                  <div className="typo-para-small text-gray-600">{source}</div>
                );
              const splitText = source.split(':');

              return (
                <div className="typo-para-small w-fit">
                  <span className="text-accent-pink-500 w-fit">
                    {splitText[0]?.trim() || ''}
                  </span>
                  <span className="text-gray-600">:</span>
                  <span className="text-accent-green-500">
                    {splitText[1] || ''}
                  </span>
                </div>
              );
            }}
            styles={{
              variables: {
                light: {
                  codeFoldGutterBackground: 'transparent',
                  codeFoldBackground: 'transparent',
                  gutterBackground: '#64748B1F',
                  gutterColor: '#64748B',
                  diffViewerBackground: '#F8FAFC',
                  addedBackground: '#DCF4DE',
                  removedBackground: '#FCE3F3',
                  wordAddedBackground: 'transparent',
                  wordRemovedBackground: 'transparent'
                }
              },
              codeFold: {
                display: 'none'
              },
              line: {
                display: 'flex',
                width: 'fit-content'
              },
              contentText: {
                display: 'flex',
                padding: '0 1px'
              },
              wordDiff: {
                padding: 0
              },
              diffContainer: {
                display: 'flex',
                padding: '12px',
                paddingLeft: 0,
                tbody: {
                  display: 'flex',
                  flexDirection: 'column',
                  rowGap: '1px',
                  tr: {
                    'td > .react-diff-1i45b3o-content-text > span > div > span:last-child':
                      {
                        color: '#64748B'
                      }
                  }
                }
              }
            }}
            extraLinesSurroundingDiff={2}
            codeFoldMessageRenderer={() => <></>}
          />

          {/* <CodeBlock>{entityData}</CodeBlock> */}
          {/* <div className="bg-gray-100 rounded-lg">
            {Diff.createTwoFilesPatch(
              'old version',
              'new version',
              previousEntityData,
              entityData
            )
              .split('\n')
              .filter(line => {
                if (
                  line.startsWith('\\') ||
                  line.startsWith('=') ||
                  line.startsWith('---') ||
                  line.startsWith('+++')
                ) {
                  return false;
                }
                return true;
              })
              .map(line => {
                if (line.startsWith('@@')) {
                  return '...';
                }
                return line;
              })
              .map((line, i) => {
                switch (line[0]) {
                  case '+':
                    return (
                      <div
                        key={i}
                        className={'bg-accent-green-50 text-accent-green-800'}
                      >
                        <span key={i}>{line}</span>
                      </div>
                    );
                  case '-':
                    return (
                      <div
                        key={i}
                        className={'bg-accent-red-50 text-accent-red-800'}
                        data-testid="deleted-line"
                      >
                        <span>{line}</span>
                      </div>
                    );
                  default:
                    return <div key={i}>{line}</div>;
                }
              })}
          </div> */}
        </TabsContent>
      </Tabs>
    </div>
  );
};

export default AuditLogItem;
