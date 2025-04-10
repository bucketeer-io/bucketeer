import { memo, useCallback, useEffect, useMemo, useRef, useState } from 'react';
import primaryAvatar from 'assets/avatars/primary.svg';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import { pickBy } from 'lodash';
import { AuditLog, DomainEventType } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { formatLongDateTime } from 'utils/date-time';
import { copyToClipBoard } from 'utils/function';
import { stringifyParams, useSearchParams } from 'utils/search-params';
import { cn } from 'utils/style';
import { IconChevronDown, IconLink, IconWatch } from '@icons';
import { AvatarImage } from 'components/avatar';
import Button from 'components/button';
import Icon from 'components/icon';
import { AuditLogTab } from '../types';
import ReactDiffViewer from './diff-viewer';

const AuditLogItem = memo(
  ({
    isExpanded,
    auditLog,
    type,
    onClick
  }: {
    isExpanded: boolean;
    auditLog: AuditLog;
    type: DomainEventType;
    onClick: () => void;
  }) => {
    const {
      editor,
      localizedMessage,
      timestamp,
      options,
      entityData,
      previousEntityData
    } = auditLog;

    const { t } = useTranslation(['common']);
    const { notify } = useToast();

    const { consoleAccount } = useAuth();
    const currentEnvironment = getCurrentEnvironment(consoleAccount!);
    const { searchOptions } = useSearchParams();

    const [currentTab, setCurrentTab] = useState<AuditLogTab>(
      AuditLogTab.CHANGES
    );

    const parsedEntityData = useMemo(
      () => JSON.parse(entityData) || {},
      [entityData]
    );

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

    const handleChangeTab = useCallback((value: AuditLogTab) => {
      lineNumberRef.current = 0;
      setCurrentTab(value);
    }, []);

    const handleCopyId = useCallback(
      (id: string) => {
        const requestParams = stringifyParams(
          pickBy(searchOptions, v => isNotEmpty(v as string))
        );
        copyToClipBoard(
          `${import.meta.env.VITE_AUTH_REDIRECT_ENDPOINT}${import.meta.env.BASE_URL}${currentEnvironment.urlCode}/audit-logs/${id}${requestParams ? `?${requestParams}` : ''}`
        );
        notify({
          toastType: 'toast',
          messageType: 'success',
          message: <span>{`Copied!`}</span>
        });
      },
      [currentEnvironment, searchOptions]
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
            'h-fit gap-y-5 min-h-[179px]': isExpanded
          }
        )}
      >
        <div className="flex items-center w-full gap-x-3">
          <AvatarImage
            image={primaryAvatar}
            alt="member-avatar"
            className="size-10"
          />
          <div className="flex flex-col flex-1 gap-y-1 truncate">
            <div className="flex items-center gap-x-1.5 max-w-full">
              <p className="typo-para-medium font-bold text-gray-700">
                {editor.name || editor.email}
              </p>
              <p className="typo-para-medium text-gray-700 min-w-fit">
                {localizedMessage.message}
              </p>
              {parsedEntityData?.name && (
                <p className="typo-para-medium text-primary-500 underline truncate">
                  {parsedEntityData?.name}
                </p>
              )}
            </div>
            <div className="flex items-center gap-x-1">
              <Icon icon={IconWatch} size={'sm'} />
              <p className="typo-para-small text-gray-500">{time}</p>
            </div>
          </div>
          <div className="flex items-center min-w-fit divide-x divide-gray-200">
            <Button
              variant={'text'}
              className="pr-5 active:scale-75"
              onClick={() => handleCopyId(auditLog.id)}
            >
              <Icon icon={IconLink} color="primary-500" />
            </Button>
            <div
              className={cn('flex-center cursor-pointer pl-5')}
              onClick={onClick}
            >
              <Icon
                icon={IconChevronDown}
                className={cn('rotate-0 transition-all duration-100', {
                  'rotate-180': isExpanded
                })}
              />
            </div>
          </div>
        </div>
        {options.comment && isExpanded && (
          <div className="flex items-center w-full p-3 bg-gray-100 rounded typo-para-medium text-gray-600 break-all">
            {options.comment}
          </div>
        )}

        <div
          className={cn(
            'flex items-center w-full justify-between h-0 opacity-0 z-[-1]',
            {
              'opacity-100 z-[0] h-10': isExpanded
            }
          )}
        >
          <p className="typo-para-small text-gray-500 uppercase">
            {t('patch')}
          </p>
          <div className="flex items-center">
            <Button
              variant={'secondary-2'}
              size={'sm'}
              className={cn(
                'typo-para-medium rounded-r-none !shadow-none border border-gray-200 hover:border-gray-400',
                {
                  'text-accent-pink-500 border-accent-pink-500 hover:text-accent-pink-500 hover:border-accent-pink-500':
                    currentTab === AuditLogTab.CHANGES
                }
              )}
              onClick={() => handleChangeTab(AuditLogTab.CHANGES)}
            >
              {t(`changes`)}
            </Button>
            <Button
              variant={'secondary-2'}
              size={'sm'}
              className={cn(
                'typo-para-medium rounded-l-none !shadow-none border border-gray-200 hover:border-gray-400',
                {
                  'text-accent-pink-500 border-accent-pink-500 hover:text-accent-pink-500 hover:border-accent-pink-500':
                    currentTab === AuditLogTab.SNAPSHOT
                }
              )}
              onClick={() => handleChangeTab(AuditLogTab.SNAPSHOT)}
            >
              {t(`snapshot`)}
            </Button>
          </div>
        </div>
        <div
          className={cn('z-[-1] h-0 opacity-0', {
            'z-[0] h-fit opacity-100': isExpanded
          })}
        >
          <ReactDiffViewer
            type={type}
            lineNumber={lineNumberRef.current}
            currentTab={currentTab}
            previousEntityData={previousEntityData}
            entityData={entityData}
          />
        </div>
      </div>
    );
  }
);

export default AuditLogItem;
