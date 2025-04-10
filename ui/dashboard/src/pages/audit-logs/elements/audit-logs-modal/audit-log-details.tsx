import { useCallback, useEffect, useRef, useState } from 'react';
// import primaryAvatar from 'assets/avatars/primary.svg';
import { useTranslation } from 'i18n';
// import ReactDiffViewer from 'pages/audit-logs/collection-layout/diff-viewer';
import { AuditLogTab } from 'pages/audit-logs/types';
// import { AvatarImage } from 'components/avatar';
import SlideModal from 'components/modal/slide';
import { Tabs, TabsContent, TabsList, TabsTrigger } from 'components/tabs';

const AuditLogDetailsModal = ({
  isOpen,
  onClose
}: {
  auditLogId: string;
  isOpen: boolean;
  onClose: () => void;
}) => {
  const { t } = useTranslation(['common']);

  const lineNumberRef = useRef(0);

  const [currentTab, setCurrentTab] = useState<AuditLogTab>(
    AuditLogTab.CHANGES
  );

  // const dateTime = useMemo(
  //   () =>
  //     formatLongDateTime({
  //       value:
  //         currentAuditLog?.timestamp ||
  //         Math.trunc(new Date().getTime() / 1000).toString(),
  //       overrideOptions: {
  //         year: 'numeric',
  //         month: '2-digit',
  //         day: '2-digit',
  //         hour12: false,
  //         hour: '2-digit',
  //         minute: '2-digit'
  //       },
  //       locale: 'ja-JP'
  //     })?.replace(' ', ' - '),
  //   [currentAuditLog]
  // );

  // const parsedEntityData = useMemo(
  //   () =>
  //     currentAuditLog?.entityData
  //       ? JSON.parse(currentAuditLog.entityData) || {}
  //       : {},
  //   [currentAuditLog]
  // );

  const handleChangeTab = useCallback((value: AuditLogTab) => {
    lineNumberRef.current = 0;
    setCurrentTab(value);
  }, []);

  useEffect(() => {
    return () => {
      lineNumberRef.current = 0;
    };
  }, []);

  return (
    <SlideModal title={t('new-api-key')} isOpen={isOpen} onClose={onClose}>
      <div className="w-full p-5 pb-28">
        <div className="flex flex-col w-full gap-y-5">
          {/* <div className="flex flex-col flex-1 gap-y-3 truncate">
              <div className="flex items-center w-full gap-x-3">
                <AvatarImage image={primaryAvatar} alt="member-avatar" />
                <div className="flex items-center gap-x-1.5 max-w-full">
                  <p className="typo-para-medium font-bold text-gray-700">
                    {currentAuditLog?.editor?.name ||
                      currentAuditLog?.editor?.email}
                  </p>
                  <p className="typo-para-medium text-gray-700 min-w-fit">
                    {currentAuditLog?.localizedMessage.message}
                  </p>
                  {parsedEntityData?.name && (
                    <p className="typo-para-medium text-primary-500 underline truncate">
                      {parsedEntityData?.name}
                    </p>
                  )}
                </div>
              </div>
              <div className="typo-para-small text-gray-500">{dateTime}</div>
            </div> */}
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
              {/* <ReactDiffViewer
                  currentTab={currentTab}
                  lineNumber={lineNumberRef.current}
                  entityData={currentAuditLog?.entityData}
                  previousEntityData={currentAuditLog?.previousEntityData}
                  type={currentAuditLog?.type}
                /> */}
            </TabsContent>
          </Tabs>
        </div>
      </div>
    </SlideModal>
  );
};

export default AuditLogDetailsModal;
