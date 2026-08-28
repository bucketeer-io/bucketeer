import { ReactNode, useEffect, useMemo, useState } from 'react';
import { useToast } from 'hooks';
import { getLanguage, useTranslation } from 'i18n';
import { Pencil, Send } from 'lucide-react';
import { formatLongDateTime, useFormatDateTime } from 'utils/date-time';
import { cn } from 'utils/style';
import Button from 'components/button';
import SlideModal from 'components/modal/slide';
import Spinner from 'components/spinner';
import { usePublishDraft } from '../collection-loader/use-fetch-notifications';
import { LANGUAGE_META } from '../language-meta';
import LanguageTabs from '../publisher/language-tabs';
import {
  NotificationDetail,
  NotificationLocalizationInput,
  NotificationStatus
} from '../types';
import { MarkdownContent } from './markdown-content';
import TagChip from './tag-chip';

interface NotificationDetailModalProps {
  notification?: NotificationDetail;
  isLoading?: boolean;
  isError?: boolean;
  isOpen: boolean;
  onClose: () => void;
  onEditDraft?: (notification: NotificationDetail) => void;
}

const formatTimestamp = (epochSeconds: string) =>
  formatLongDateTime({ value: epochSeconds });

const PersonBadge = ({ email }: { email: string }) => {
  const initial = (email.trim()[0] ?? '?').toUpperCase();
  const name = email.split('@')[0];
  return (
    <div className="flex min-w-0 items-center gap-2">
      <span className="truncate typo-para-small text-gray-900">{name}</span>
      <span className="flex-center shrink-0 size-6 rounded-full bg-primary-100 typo-para-tiny font-medium text-primary-600">
        {initial}
      </span>
    </div>
  );
};

const DetailRow = ({
  label,
  children
}: {
  label: string;
  children: ReactNode;
}) => (
  <div className="flex min-h-8 flex-wrap items-center justify-between gap-x-4 gap-y-1">
    <span className="shrink-0 typo-para-small text-gray-500">{label}</span>
    <div className="min-w-0 typo-para-small text-gray-900">{children}</div>
  </div>
);

const Section = ({
  title,
  className,
  children
}: {
  title: string;
  className?: string;
  children: ReactNode;
}) => (
  <section className={cn('flex flex-col gap-3', className)}>
    <h3 className="typo-para-medium font-bold text-gray-900">{title}</h3>
    {children}
  </section>
);

const NotificationDetailModal = ({
  notification,
  isLoading,
  isError,
  isOpen,
  onClose,
  onEditDraft
}: NotificationDetailModalProps) => {
  const { t } = useTranslation(['common', 'table', 'form', 'message']);
  const formatDateTime = useFormatDateTime();
  const { notify, errorNotify } = useToast();
  const publishMutation = usePublishDraft();

  const localizations: NotificationLocalizationInput[] = useMemo(() => {
    if (notification?.localizations?.length) return notification.localizations;
    if (!notification) return [];
    return [
      {
        language: getLanguage(),
        title: notification.title,
        content: notification.content,
        tags: notification.tags
      }
    ];
  }, [notification]);

  const [activeLanguage, setActiveLanguage] = useState<string>(() =>
    getLanguage()
  );

  useEffect(() => {
    if (!localizations.length) return;
    setActiveLanguage(current => {
      if (localizations.some(l => l.language === current)) return current;
      const preferred = getLanguage();
      return localizations.some(l => l.language === preferred)
        ? preferred
        : localizations[0].language;
    });
  }, [localizations]);

  if (!isOpen) return null;

  if (isError) {
    return (
      <SlideModal title="" isOpen={isOpen} onClose={onClose}>
        <div className="flex h-full w-full flex-col items-center justify-center gap-4">
          <p className="typo-para-medium text-gray-500">
            {t('message:notification-load-failed')}
          </p>
          <Button variant="secondary" onClick={onClose}>
            {t('close')}
          </Button>
        </div>
      </SlideModal>
    );
  }

  if (!notification || isLoading) {
    return (
      <SlideModal title="" isOpen={isOpen} onClose={onClose}>
        <div className="flex h-full w-full items-center justify-center">
          <Spinner />
        </div>
      </SlideModal>
    );
  }

  const isDraft = notification.status === NotificationStatus.DRAFT;
  const timestamp = isDraft ? notification.updatedAt : notification.publishedAt;
  const activeLocalization =
    localizations.find(l => l.language === activeLanguage) ?? localizations[0];

  const onPublish = () => {
    publishMutation.mutate(notification.id, {
      onSuccess: () => {
        notify({ message: t('message:published-successfully') });
        onClose();
      },
      onError: error => errorNotify(error)
    });
  };

  return (
    <SlideModal
      title={activeLocalization.title}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="w-full h-full flex flex-col">
        <div className="flex flex-1 flex-col gap-6 overflow-auto p-4 sm:p-6">
          {localizations.length > 1 && (
            <LanguageTabs
              fields={localizations.map(l => ({
                id: l.language,
                language: l.language
              }))}
              activeLanguage={activeLanguage}
              availableToAdd={[]}
              canRemove={false}
              readOnly
              languageMeta={LANGUAGE_META}
              tooltipContent={t('form:languages-info-view')}
              onSelect={setActiveLanguage}
              onAdd={() => {}}
              onRemove={() => {}}
            />
          )}

          <div className="flex items-center gap-2">
            <TagChip
              tag={{
                name: isDraft ? t('draft-status') : t('published-status'),
                color: isDraft ? '#6B7280' : '#5D5FEF'
              }}
              hideDot
            />
            <span className="text-gray-300">•</span>
            <span className="typo-para-small text-gray-500">
              {formatDateTime(timestamp)}
            </span>
          </div>

          {activeLocalization.tags.length > 0 && (
            <div className="flex flex-wrap items-center gap-2">
              {activeLocalization.tags.map(tag => (
                <TagChip key={tag.name} tag={tag} />
              ))}
            </div>
          )}

          <Section title={t('form:content')}>
            <MarkdownContent source={activeLocalization.content} />
          </Section>

          <Section
            title={t('details')}
            className="border-t border-gray-200 pt-6"
          >
            <DetailRow label={t('form:created-by')}>
              <PersonBadge email={notification.createdBy} />
            </DetailRow>
            <DetailRow label={t('table:created-at')}>
              {formatTimestamp(notification.createdAt)}
            </DetailRow>
            <DetailRow label={t('table:updated-at')}>
              {formatTimestamp(notification.updatedAt)}
            </DetailRow>
            <DetailRow label={t('form:last-edited-by')}>
              <PersonBadge email={notification.lastEditedBy} />
            </DetailRow>
          </Section>
        </div>

        <div className="flex flex-wrap items-center justify-end gap-3 border-t border-gray-200 px-4 py-4 sm:px-6">
          <Button variant="secondary" onClick={onClose}>
            {t('close')}
          </Button>
          {isDraft && onEditDraft && (
            <Button
              variant="secondary"
              onClick={() => onEditDraft(notification)}
            >
              <Pencil size={16} />
              {t('form:edit-draft')}
            </Button>
          )}
          {isDraft && (
            <Button onClick={onPublish} loading={publishMutation.isPending}>
              <Send size={16} />
              {t('publish')}
            </Button>
          )}
        </div>
      </div>
    </SlideModal>
  );
};

export default NotificationDetailModal;
