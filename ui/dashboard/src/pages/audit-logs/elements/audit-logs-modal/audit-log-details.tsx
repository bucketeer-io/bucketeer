import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { Trans } from 'react-i18next';
import { useParams } from 'react-router-dom';
import { useQueryAuditLogDetails } from '@queries/audit-log-details';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast } from 'hooks';
import { useAuditLogDataPatterns } from 'hooks/use-audit-log-data-patterns';
import { getLanguage, useTranslation } from 'i18n';
import { formatLongDateTime } from 'utils/date-time';
import { cn } from 'utils/style';
import AuditLogAvatar from 'pages/audit-logs/collection-layout/audit-log-avatar';
import AuditLogTitle from 'pages/audit-logs/collection-layout/audit-log-title';
import ReactDiffViewer from 'pages/audit-logs/collection-layout/json-compare';
import { AuditLogTab } from 'pages/audit-logs/types';
import { getActionText } from 'pages/audit-logs/utils';
import SlideModal from 'components/modal/slide';
import { Tabs, TabsContent, TabsList, TabsTrigger } from 'components/tabs';
import DateTooltip from 'elements/date-tooltip';
import FormLoading from 'elements/form-loading';

const AuditLogDetailsModal = ({
  auditLogId,
  isOpen,
  title,
  onClose
}: {
  auditLogId: string;
  isOpen: boolean;
  title?: string;
  onClose: () => void;
}) => {
  const { t } = useTranslation(['common', 'table']);

  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const params = useParams();

  const { errorNotify } = useToast();
  const lineNumberRef = useRef(0);
  const isLanguageJapanese = getLanguage() === 'ja';

  const [currentTab, setCurrentTab] = useState<AuditLogTab>(
    AuditLogTab.CHANGES
  );

  const {
    data: collection,
    isLoading,
    error
  } = useQueryAuditLogDetails({
    params: {
      environmentId: currentEnvironment.id,
      id: auditLogId
    },
    enabled: params?.envUrlCode === currentEnvironment?.urlCode
  });

  const auditLog = collection?.auditLog;

  // Use centralized audit log data patterns logic
  const {
    isHaveEntityData,
    parsedEntityData,
    effectiveIsSameData,
    shouldShowChanges,
    displayMode,
    entityData,
    previousEntityData
  } = useAuditLogDataPatterns(auditLog);

  const dateTime = useMemo(
    () =>
      auditLog
        ? formatLongDateTime({
            value: auditLog?.timestamp,
            overrideOptions: {
              year: 'numeric',
              month: '2-digit',
              day: '2-digit',
              hour12: false,
              hour: '2-digit',
              minute: '2-digit'
            },
            locale: 'ja-JP'
          })?.replace(' ', ' - ')
        : '',
    [auditLog]
  );

  const handleChangeTab = useCallback((value: AuditLogTab) => {
    lineNumberRef.current = 0;
    setCurrentTab(value);
  }, []);

  useEffect(() => {
    return () => {
      lineNumberRef.current = 0;
    };
  }, []);

  useEffect(() => {
    if (error) {
      errorNotify(error);
    }
  }, [error]);

  return (
    <SlideModal
      title={title || t('audit-log-details')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="w-full p-5 pr-2">
        {isLoading ? (
          <FormLoading />
        ) : error ? (
          <div className="typo-para-medium text-gray-500">
            <Trans
              i18nKey={'form:not-found-entity'}
              values={{
                entity: t('audit-log')
              }}
            />
          </div>
        ) : (
          <div className="flex flex-col w-full gap-y-5">
            <div className="flex flex-col flex-1 gap-y-3 truncate">
              <div className="flex items-center w-full gap-x-3">
                <AuditLogAvatar editor={auditLog?.editor} className="size-8" />
                <div
                  className={cn(
                    'flex items-center gap-x-1.5 max-w-full typo-para-medium font-normal text-gray-700 truncate',
                    {
                      'gap-x-0': isLanguageJapanese
                    }
                  )}
                >
                  {auditLog && (
                    <AuditLogTitle
                      isHaveEntityData={isHaveEntityData}
                      entityId={parsedEntityData?.id}
                      action={getActionText(auditLog.type, isLanguageJapanese)}
                      entityName={
                        parsedEntityData?.name ||
                        parsedEntityData?.feature_name ||
                        ''
                      }
                      entityType={auditLog.entityType}
                      urlCode={currentEnvironment.urlCode}
                      username={auditLog.editor.name || auditLog.editor.email}
                      additionalText={
                        parsedEntityData?.name && isLanguageJapanese && 'の'
                      }
                    />
                  )}
                </div>
              </div>
              <DateTooltip
                align="start"
                trigger={
                  <div className="typo-para-small text-gray-500 w-fit">
                    {dateTime}
                  </div>
                }
                date={auditLog?.timestamp || null}
              />
            </div>
            {auditLog?.options?.comment && (
              <div className="flex items-center w-full p-3 bg-gray-100 rounded typo-para-small text-gray-600 break-all border-l-4 border-gray-500">
                {auditLog?.options?.comment}
              </div>
            )}

            <Tabs
              className="flex w-full flex-col gap-y-4"
              value={currentTab}
              onValueChange={value => handleChangeTab(value as AuditLogTab)}
            >
              {shouldShowChanges && (
                <TabsList>
                  <TabsTrigger value={AuditLogTab.CHANGES}>
                    {t(`changes`)}
                  </TabsTrigger>
                  <TabsTrigger value={AuditLogTab.SNAPSHOT}>
                    {t(`snapshot`)}
                  </TabsTrigger>
                </TabsList>
              )}

              <p className="typo-para-small text-gray-500 uppercase">
                {t(
                  currentTab === AuditLogTab.SNAPSHOT
                    ? 'current-version'
                    : displayMode
                )}
              </p>

              <TabsContent
                value={currentTab}
                className="flex flex-col gap-y-4 mt-0"
              >
                {auditLog && (
                  <ReactDiffViewer
                    isSameData={effectiveIsSameData}
                    prefix="line-00"
                    currentTab={currentTab}
                    lineNumber={lineNumberRef.current}
                    entityData={entityData}
                    previousEntityData={previousEntityData}
                  />
                )}
              </TabsContent>
            </Tabs>
          </div>
        )}
      </div>
    </SlideModal>
  );
};

export default AuditLogDetailsModal;
