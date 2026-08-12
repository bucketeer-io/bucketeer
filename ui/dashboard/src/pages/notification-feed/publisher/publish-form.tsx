import { useEffect, useMemo, useRef, useState } from 'react';
import { FormProvider, useFieldArray, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { useToast } from 'hooks';
import useFormSchema from 'hooks/use-form-schema';
import { getLanguage, Language, useTranslation } from 'i18n';
import { IconEnglishFlag, IconJapanFlag } from '@icons';
import Button from 'components/button';
import Form from 'components/form';
import Input from 'components/input';
import {
  usePublishNotification,
  useSaveDraft,
  useUpdateNotification
} from '../collection-loader/use-fetch-notifications';
import {
  NotificationDetail,
  NotificationLocalizationInput,
  NotificationStatus,
  PublishNotificationInput
} from '../types';
import { formSchema, PublishFormValues } from './form-schema';
import LanguageTabs from './language-tabs';
import MarkdownEditor from './markdown-editor';
import TagSelect from './tag-select';

const LANGUAGE_META = {
  [Language.ENGLISH]: {
    label: 'English',
    englishName: 'English',
    icon: IconEnglishFlag
  },
  [Language.JAPANESE]: {
    label: '日本語',
    englishName: 'Japanese',
    icon: IconJapanFlag
  }
};

// Languages the form can author, in the order they appear in the add menu.
const FORM_LANGUAGES: Language[] = [Language.ENGLISH, Language.JAPANESE];

const emptyLocalization = (
  language: string
): NotificationLocalizationInput => ({
  language,
  title: '',
  content: '',
  tags: []
});

const PublishForm = ({
  disabled,
  initialDraft,
  onClear
}: {
  disabled?: boolean;
  initialDraft?: NotificationDetail;
  onClear?: () => void;
}) => {
  const { t, i18n } = useTranslation(['common', 'form', 'message']);
  const { notify, errorNotify } = useToast();

  const defaultLanguage = useMemo(() => {
    const lang = getLanguage();
    return FORM_LANGUAGES.includes(lang) ? lang : Language.ENGLISH;
  }, [i18n.language]);

  const buildLocalizations = (
    draft?: NotificationDetail
  ): NotificationLocalizationInput[] => {
    if (draft?.localizations?.length) {
      return draft.localizations.map(loc => ({
        language: loc.language,
        title: loc.title,
        content: loc.content,
        tags: loc.tags
      }));
    }
    return [emptyLocalization(defaultLanguage)];
  };

  const initialActiveLanguage = (
    locs: NotificationLocalizationInput[]
  ): string =>
    locs.some(l => l.language === defaultLanguage)
      ? defaultLanguage
      : (locs[0]?.language ?? defaultLanguage);

  const form = useForm<PublishFormValues>({
    resolver: yupResolver(useFormSchema(formSchema)),
    mode: 'onChange',
    defaultValues: { localizations: buildLocalizations(initialDraft) }
  });

  const { control } = form;
  const {
    fields: localizationFields,
    append: appendLocalization,
    remove: removeLocalization
  } = useFieldArray({
    control,
    name: 'localizations'
  });

  const [activeLanguage, setActiveLanguage] = useState<string>(() =>
    initialActiveLanguage(buildLocalizations(initialDraft))
  );

  useEffect(() => {
    if (initialDraft) {
      const locs = buildLocalizations(initialDraft);
      form.reset({ localizations: locs });
      setActiveLanguage(initialActiveLanguage(locs));
      return;
    }
    if (form.formState.isDirty) return;
    form.reset({ localizations: [emptyLocalization(defaultLanguage)] });
    setActiveLanguage(defaultLanguage);
  }, [initialDraft?.id, initialDraft?.updatedAt, defaultLanguage]);

  const activeIndex = Math.max(
    0,
    localizationFields.findIndex(f => f.language === activeLanguage)
  );
  const canRemoveLanguage = localizationFields.length > 1;
  const availableToAdd = FORM_LANGUAGES.filter(
    lang => !localizationFields.some(f => f.language === lang)
  );

  const publishMutation = usePublishNotification();
  const saveDraftMutation = useSaveDraft();
  const updateMutation = useUpdateNotification();

  const isEditing = !!initialDraft;
  const editingId = initialDraft?.id;

  const editingIdRef = useRef(editingId);
  editingIdRef.current = editingId;

  const {
    formState: { isValid }
  } = form;
  const canSubmit = !disabled && isValid;
  const [pendingAction, setPendingAction] = useState<
    'publish' | 'draft' | null
  >(null);
  const isPublishPending =
    publishMutation.isPending ||
    (updateMutation.isPending && pendingAction === 'publish');
  const isDraftPending =
    saveDraftMutation.isPending ||
    (updateMutation.isPending && pendingAction === 'draft');

  const addLanguage = (language: string) => {
    if (localizationFields.some(f => f.language === language)) return;
    appendLocalization(emptyLocalization(language));
    setActiveLanguage(language);
  };

  const removeLanguage = (index: number, language: string) => {
    if (!canRemoveLanguage) return;
    removeLocalization(index);
    if (activeLanguage === language) {
      const next = localizationFields.find(f => f.language !== language);
      if (next) setActiveLanguage(next.language);
    }
  };

  const toInput = (status: NotificationStatus): PublishNotificationInput => ({
    status,
    localizations: form.getValues('localizations')
  });

  const resetForm = () => {
    form.reset({ localizations: [emptyLocalization(defaultLanguage)] });
    setActiveLanguage(defaultLanguage);
    onClear?.();
  };

  const handlePublish = form.handleSubmit(() => {
    const payload = toInput(NotificationStatus.PUBLISHED);
    const submittedFor = editingId;
    setPendingAction('publish');
    const onDone = {
      onSuccess: () => {
        if (editingIdRef.current !== submittedFor) return;
        notify({ message: t('message:published-successfully') });
        resetForm();
      },
      onError: (error: Error) => {
        if (editingIdRef.current !== submittedFor) return;
        errorNotify(error);
      }
    };
    if (isEditing && editingId) {
      updateMutation.mutate({ id: editingId, input: payload }, onDone);
    } else {
      publishMutation.mutate(payload, onDone);
    }
  });

  const handleSaveDraft = form.handleSubmit(() => {
    const payload = toInput(NotificationStatus.DRAFT);
    const submittedFor = editingId;
    setPendingAction('draft');
    const onDone = {
      onSuccess: () => {
        if (editingIdRef.current !== submittedFor) return;
        notify({ message: t('message:draft-saved') });
        resetForm();
      },
      onError: (error: Error) => {
        if (editingIdRef.current !== submittedFor) return;
        errorNotify(error);
      }
    };
    if (isEditing && editingId) {
      updateMutation.mutate({ id: editingId, input: payload }, onDone);
    } else {
      saveDraftMutation.mutate(payload, onDone);
    }
  });

  const handleClear = () => resetForm();

  return (
    <FormProvider {...form}>
      <Form onSubmit={handlePublish} className="flex flex-col gap-6">
        <LanguageTabs
          fields={localizationFields}
          activeLanguage={activeLanguage}
          availableToAdd={availableToAdd}
          canRemove={canRemoveLanguage}
          languageMeta={LANGUAGE_META}
          onSelect={setActiveLanguage}
          onAdd={addLanguage}
          onRemove={removeLanguage}
        />

        <Form.Field
          control={control}
          name={`localizations.${activeIndex}.title`}
          render={({ field }) => (
            <Form.Item>
              <Form.Label required>{t('title')}</Form.Label>
              <Form.Control>
                <Input
                  {...field}
                  placeholder={t('form:notification-title-placeholder')}
                />
              </Form.Control>
              <Form.Message />
            </Form.Item>
          )}
        />

        <Form.Field
          control={control}
          name={`localizations.${activeIndex}.tags`}
          render={({ field }) => (
            <Form.Item>
              <Form.Label>{t('tags')}</Form.Label>
              <Form.Control>
                <TagSelect
                  value={field.value}
                  language={activeLanguage}
                  onChange={field.onChange}
                />
              </Form.Control>
              <Form.Message />
            </Form.Item>
          )}
        />

        <Form.Field
          control={control}
          name={`localizations.${activeIndex}.content`}
          render={({ field }) => (
            <Form.Item>
              <Form.Label required>{t('description')}</Form.Label>
              <Form.Control>
                <MarkdownEditor
                  value={field.value}
                  onChange={field.onChange}
                  placeholder={t('form:description-placeholder')}
                />
              </Form.Control>
              <Form.Message />
            </Form.Item>
          )}
        />

        <div className="flex items-center gap-4">
          <Button
            type="submit"
            disabled={!canSubmit}
            loading={isPublishPending}
          >
            {t('publish')}
          </Button>
          <Button
            type="button"
            variant="secondary"
            onClick={handleSaveDraft}
            disabled={!canSubmit}
            loading={isDraftPending}
          >
            {isEditing ? t('form:update-draft') : t('save-draft')}
          </Button>
          <Button
            type="button"
            variant="secondary"
            onClick={handleClear}
            disabled={disabled}
          >
            {t('clear')}
          </Button>
        </div>
      </Form>
    </FormProvider>
  );
};

export default PublishForm;
