import {
  memo,
  MouseEvent,
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState
} from 'react';
import { useParams } from 'react-router-dom';
import { getCurrentEnvironment, useAuth } from 'auth';
import { urls } from 'configs';
import {
  PAGE_PATH_FEATURE_HISTORY,
  PAGE_PATH_FEATURES
} from 'constants/routing';
import { useToast } from 'hooks';
import { useTranslation, getLanguage } from 'i18n';
import { AuditLog } from '@types';
import { formatLongDateTime } from 'utils/date-time';
import { copyToClipBoard } from 'utils/function';
import { cn } from 'utils/style';
import { IconChevronDown, IconLink, IconWatch } from '@icons';
import Button from 'components/button';
import Icon from 'components/icon';
import { Tooltip } from 'components/tooltip';
import DateTooltip from 'elements/date-tooltip';
import { useAuditLogDataPatterns } from '../hooks/use-audit-log-data-patterns';
import { AuditLogTab } from '../types';
import { getActionText } from '../utils';
import AuditLogAvatar from './audit-log-avatar';
import AuditLogTitle from './audit-log-title';
import AuditLogJSONCompare from './json-compare';

const AuditLogItem = memo(
  ({
    isExpanded,
    auditLog,
    prefix,
    onClick
  }: {
    isExpanded: boolean;
    auditLog: AuditLog;
    prefix: string;
    onClick: () => void;
  }) => {
    const { editor, timestamp, options, entityType, type } = auditLog;
    const isLanguageJapanese = getLanguage() === 'ja';
    const { t } = useTranslation(['common', 'table', 'message']);
    const { notify } = useToast();

    const { consoleAccount } = useAuth();
    const currentEnvironment = getCurrentEnvironment(consoleAccount!);
    const params = useParams();

    const [currentTab, setCurrentTab] = useState<AuditLogTab>(
      AuditLogTab.CHANGES
    );

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

    const lineNumberRef = useRef(0);
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

    const buttonCls = useMemo(
      () =>
        'typo-para-medium !text-gray-600 !shadow-none border border-gray-200 hover:border-gray-400',
      []
    );
    const buttonActiveCls = useMemo(
      () =>
        '!text-accent-pink-500 border-accent-pink-500 hover:!text-accent-pink-500 hover:border-accent-pink-500',
      []
    );

    const handleChangeTab = useCallback((value: AuditLogTab) => {
      lineNumberRef.current = 0;
      setCurrentTab(value);
    }, []);

    const handleCopyId = useCallback(
      (id: string) => {
        const flagId = params?.flagId;

        copyToClipBoard(
          `${urls.ORIGIN_URL}${currentEnvironment.urlCode}${flagId ? `${PAGE_PATH_FEATURES}/${flagId}${PAGE_PATH_FEATURE_HISTORY}` : '/audit-logs'}/${id}`
        );

        notify({
          message: t('message:copied')
        });
      },
      [currentEnvironment, params]
    );

    const handleOnExpandAuditLog = useCallback(
      (e: MouseEvent) => {
        if (e.target === e.currentTarget) {
          onClick();
        }
      },
      [onClick]
    );

    useEffect(() => {
      return () => {
        lineNumberRef.current = 0;
      };
    }, []);

    return (
      <div
        className={cn(
          'flex flex-col w-full p-3 bg-white shadow-card rounded-lg h-[73px] min-h-[73px] transition-all duration-100',
          {
            'h-fit min-h-[179px]': isExpanded && isHaveEntityData,
            'h-fit': !!options.comment
          }
        )}
      >
        <div className="flex items-center w-full">
          <AuditLogAvatar editor={editor} />
          <div
            className="flex flex-col flex-1 gap-y-1 truncate px-3 cursor-pointer"
            onClick={handleOnExpandAuditLog}
          >
            <div
              className={cn(
                'flex items-center gap-x-1.5 w-fit max-w-full typo-para-medium font-normal text-gray-700 truncate cursor-default',
                {
                  'gap-x-0': isLanguageJapanese
                }
              )}
            >
              <AuditLogTitle
                isHaveEntityData={isHaveEntityData}
                entityId={parsedEntityData?.id}
                action={getActionText(type, isLanguageJapanese)}
                entityName={
                  parsedEntityData?.name || parsedEntityData?.feature_name || ''
                }
                entityType={entityType}
                urlCode={currentEnvironment.urlCode}
                username={editor.name || editor.email}
                additionalText={
                  parsedEntityData?.name && isLanguageJapanese && 'の'
                }
              />
            </div>
            <DateTooltip
              align="start"
              alignOffset={5}
              trigger={
                <div className="flex items-center gap-x-1 w-fit cursor-default">
                  <Icon icon={IconWatch} size={'sm'} />
                  <p className="typo-para-small text-gray-500">{time}</p>
                </div>
              }
              date={timestamp}
            />
          </div>
          {isHaveEntityData && (
            <div className="flex items-center min-w-fit divide-x divide-gray-200">
              <Tooltip
                align="end"
                sideOffset={-6}
                content={t('table:copy-entry-link')}
                trigger={
                  <Button
                    variant={'text'}
                    className="pr-5 active:scale-75"
                    onClick={() => handleCopyId(auditLog.id)}
                  >
                    <Icon icon={IconLink} color="primary-500" size={'sm'} />
                  </Button>
                }
              />
              <Tooltip
                align="end"
                sideOffset={-6}
                content={t('table:show-entry')}
                trigger={
                  <Button variant={'text'} className="pl-5" onClick={onClick}>
                    <Icon
                      icon={IconChevronDown}
                      className={cn('rotate-0 transition-all duration-100', {
                        'rotate-180': isExpanded
                      })}
                      size={'sm'}
                    />
                  </Button>
                }
              />
            </div>
          )}
        </div>
        {options.comment && (
          <div className="pt-3 cursor-pointer" onClick={onClick}>
            <div className="flex items-center w-full p-3 bg-gray-100 rounded typo-para-small text-gray-600 break-all border-l-4 border-gray-500">
              {options.comment}
            </div>
          </div>
        )}
        {isHaveEntityData && (
          <>
            <div
              className={cn('', {
                'py-3 cursor-pointer': isExpanded
              })}
              onClick={handleOnExpandAuditLog}
            >
              <div
                className={cn(
                  'flex items-center w-full justify-between h-0 opacity-0 z-[-1]',
                  {
                    'opacity-100 z-[0] h-10': isExpanded
                  }
                )}
                onClick={handleOnExpandAuditLog}
              >
                <p className="typo-para-small text-gray-500 uppercase cursor-default">
                  {t(
                    currentTab === AuditLogTab.SNAPSHOT
                      ? 'current-version'
                      : displayMode
                  )}
                </p>
                {shouldShowChanges && (
                  <div className="flex items-center">
                    <Button
                      variant={'secondary-2'}
                      size={'sm'}
                      className={cn(
                        'rounded-r-none',
                        buttonCls,
                        currentTab === AuditLogTab.CHANGES && buttonActiveCls,
                        { 'pointer-events-none': !isExpanded }
                      )}
                      onClick={() => handleChangeTab(AuditLogTab.CHANGES)}
                    >
                      {t(`changes`)}
                    </Button>
                    <Button
                      variant={'secondary-2'}
                      size={'sm'}
                      className={cn(
                        'rounded-l-none',
                        buttonCls,
                        currentTab === AuditLogTab.SNAPSHOT && buttonActiveCls,
                        { 'pointer-events-none': !isExpanded }
                      )}
                      onClick={() => handleChangeTab(AuditLogTab.SNAPSHOT)}
                    >
                      {t(`snapshot`)}
                    </Button>
                  </div>
                )}
              </div>
            </div>
            {isExpanded && (
              <div
                className={cn('z-[-1] h-0 opacity-0', {
                  'z-[0] h-fit opacity-100': isExpanded
                })}
              >
                <AuditLogJSONCompare
                  isSameData={effectiveIsSameData}
                  prefix={prefix}
                  lineNumber={lineNumberRef.current}
                  currentTab={currentTab}
                  previousEntityData={previousEntityData}
                  entityData={entityData}
                />
              </div>
            )}
          </>
        )}
      </div>
    );
  }
);

export default AuditLogItem;
