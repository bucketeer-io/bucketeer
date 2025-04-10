import { memo, useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { Trans } from 'react-i18next';
import { Link, useNavigate } from 'react-router-dom';
import primaryAvatar from 'assets/avatars/primary.svg';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast } from 'hooks';
import { useTranslation, getLanguage } from 'i18n';
import { AuditLog } from '@types';
import { formatLongDateTime } from 'utils/date-time';
import { copyToClipBoard } from 'utils/function';
import { cn } from 'utils/style';
import { IconChevronDown, IconLink, IconWatch } from '@icons';
import { AvatarImage } from 'components/avatar';
import Button from 'components/button';
import Icon from 'components/icon';
import { AuditLogTab } from '../types';
import { getActionText, getEntityTypeText, getPathName } from '../utils';
import ReactDiffViewer from './diff-viewer';

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
    const {
      editor,
      timestamp,
      options,
      entityData,
      previousEntityData,
      type,
      entityType
    } = auditLog;
    const isLanguageJapanese = getLanguage() === 'ja';
    const { t } = useTranslation(['common', 'table']);
    const { notify } = useToast();
    const navigate = useNavigate();

    const { consoleAccount } = useAuth();
    const currentEnvironment = getCurrentEnvironment(consoleAccount!);

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
        copyToClipBoard(
          `${import.meta.env.VITE_AUTH_REDIRECT_ENDPOINT}${import.meta.env.BASE_URL}${currentEnvironment.urlCode}/audit-logs/${id}`
        );
        notify({
          toastType: 'toast',
          messageType: 'success',
          message: <span>{`Copied!`}</span>
        });
      },
      [currentEnvironment]
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
            <div
              className={cn(
                'flex items-center gap-x-1.5 max-w-full typo-para-medium font-normal text-gray-700 truncate',
                {
                  'gap-x-0': isLanguageJapanese
                }
              )}
            >
              <Trans
                i18nKey="table:audit-log-title"
                values={{
                  username: editor.name || editor.email,
                  action: getActionText(type, isLanguageJapanese),
                  entityType: getEntityTypeText(entityType),
                  entityName: parsedEntityData?.name,
                  additionalText:
                    parsedEntityData?.name && isLanguageJapanese && 'の'
                }}
                components={{
                  b: <span className="font-bold text-gray-700 -mt-0.5" />,
                  highlight: (
                    <Link
                      to={getPathName(auditLog.id, entityType) as string}
                      onClick={e => {
                        e.preventDefault();
                        const pathName = getPathName(auditLog.id, entityType);
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
                'rounded-r-none',
                buttonCls,
                currentTab === AuditLogTab.CHANGES && buttonActiveCls
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
                currentTab === AuditLogTab.SNAPSHOT && buttonActiveCls
              )}
              onClick={() => handleChangeTab(AuditLogTab.SNAPSHOT)}
            >
              {t(`snapshot`)}
            </Button>
          </div>
        </div>
        {isExpanded && (
          <div
            className={cn('z-[-1] h-0 opacity-0', {
              'z-[0] h-fit opacity-100': isExpanded
            })}
          >
            <ReactDiffViewer
              prefix={prefix}
              type={type}
              lineNumber={lineNumberRef.current}
              currentTab={currentTab}
              previousEntityData={previousEntityData}
              entityData={entityData}
            />
          </div>
        )}
      </div>
    );
  }
);

export default AuditLogItem;
