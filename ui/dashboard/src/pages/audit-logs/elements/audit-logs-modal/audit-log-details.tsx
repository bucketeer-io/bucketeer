import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { Trans } from 'react-i18next';
import { Link, useNavigate } from 'react-router-dom';
import { useQueryAuditLogDetails } from '@queries/audit-log-details';
import primaryAvatar from 'assets/avatars/primary.svg';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast } from 'hooks';
import { getLanguage, useTranslation } from 'i18n';
import { formatLongDateTime } from 'utils/date-time';
import { cn } from 'utils/style';
import ReactDiffViewer from 'pages/audit-logs/collection-layout/diff-viewer';
import { AuditLogTab } from 'pages/audit-logs/types';
import {
  getActionText,
  getEntityTypeText,
  getPathName
} from 'pages/audit-logs/utils';
import { AvatarImage } from 'components/avatar';
import SlideModal from 'components/modal/slide';
import { Tabs, TabsContent, TabsList, TabsTrigger } from 'components/tabs';
import FormLoading from 'elements/form-loading';

const AuditLogDetailsModal = ({
  auditLogId,
  isOpen,
  onClose
}: {
  auditLogId: string;
  isOpen: boolean;
  onClose: () => void;
}) => {
  const { t } = useTranslation(['common', 'table']);

  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const { errorNotify } = useToast();
  const navigate = useNavigate();
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
    }
  });

  const auditLog = collection?.auditLog;

  const parsedEntityData = useMemo(
    () => (auditLog?.entityData ? JSON.parse(auditLog?.entityData) : {}),
    [auditLog?.entityData]
  );

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
      title={t('audit-log-details')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="w-full p-5 pb-28">
        {isLoading ? (
          <FormLoading />
        ) : (
          <div className="flex flex-col w-full gap-y-5">
            <div className="flex flex-col flex-1 gap-y-3 truncate">
              <div className="flex items-center w-full gap-x-3">
                <AvatarImage image={primaryAvatar} alt="member-avatar" />
                <div
                  className={cn(
                    'flex items-center gap-x-1.5 max-w-full typo-para-medium font-normal text-gray-700 truncate',
                    {
                      'gap-x-0': isLanguageJapanese
                    }
                  )}
                >
                  {auditLog && (
                    <Trans
                      i18nKey="table:audit-log-title"
                      values={{
                        username:
                          auditLog?.editor.name || auditLog?.editor.email,
                        action: getActionText(
                          auditLog?.type,
                          isLanguageJapanese
                        ),
                        entityType: getEntityTypeText(auditLog?.entityType),
                        entityName: parsedEntityData?.name,
                        additionalText:
                          parsedEntityData?.name && isLanguageJapanese && 'の'
                      }}
                      components={{
                        b: <span className="font-bold text-gray-700 -mt-0.5" />,
                        highlight: (
                          <Link
                            to={
                              getPathName(
                                auditLog.id,
                                auditLog.entityType
                              ) as string
                            }
                            onClick={e => {
                              e.preventDefault();
                              const pathName = getPathName(
                                auditLog.id,
                                auditLog.entityType
                              );
                              if (pathName)
                                navigate(
                                  `/${currentEnvironment.urlCode}${pathName}`,
                                  {
                                    replace: false
                                  }
                                );
                            }}
                            className="text-primary-500 underline truncate"
                          />
                        )
                      }}
                    />
                  )}
                </div>
              </div>
              <div className="typo-para-small text-gray-500">{dateTime}</div>
            </div>
            {auditLog?.options?.comment && (
              <div className="flex items-center w-full p-3 bg-gray-100 rounded typo-para-medium text-gray-600 break-all">
                {auditLog?.options?.comment}
              </div>
            )}

            <Tabs
              className="flex w-full flex-col"
              value={currentTab}
              onValueChange={value => handleChangeTab(value as AuditLogTab)}
            >
              <TabsList>
                <TabsTrigger value={AuditLogTab.CHANGES}>
                  {t(`changes`)}
                </TabsTrigger>
                <TabsTrigger value={AuditLogTab.SNAPSHOT}>
                  {t(`snapshot`)}
                </TabsTrigger>
              </TabsList>
              <TabsContent value={currentTab} className="flex flex-col gap-y-4">
                {auditLog && (
                  <ReactDiffViewer
                    prefix="line-00"
                    currentTab={currentTab}
                    lineNumber={lineNumberRef.current}
                    entityData={auditLog?.entityData}
                    previousEntityData={auditLog.previousEntityData}
                    type={auditLog.type}
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
