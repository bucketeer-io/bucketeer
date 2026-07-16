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

// Per-language display metadata: native name, English name (for the add menu),
// and flag icon. Mirrors the shared `languageList` used elsewhere.
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
  environmentId,
  initialDraft,
  onClear
}: {
  disabled?: boolean;
  environmentId: string;
  initialDraft?: NotificationDetail;
  onClear?: () => void;
}) => {
  const { t } = useTranslation(['common', 'form', 'message']);
  const { notify, errorNotify } = useToast();

  // The language the form starts with (the console language). It is only the
  // initial default — the author may remove it and author another language —
  // so it must stay fixed for the component's lifetime and not react to the
  // console language changing while the form is open, or in-progress input
  // would be reset out from under the author.
  const defaultLanguage = useMemo(() => {
    const lang = getLanguage();
    return FORM_LANGUAGES.includes(lang) ? lang : Language.ENGLISH;
  }, []);

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

  // `activeLanguage` (which tab is shown) is UI-only — not part of the
  // submitted payload.
  const [activeLanguage, setActiveLanguage] = useState<string>(() =>
    initialActiveLanguage(buildLocalizations(initialDraft))
  );

  // Reload the form when a different draft is opened for editing, showing that
  // draft's first language tab.
  useEffect(() => {
    const locs = buildLocalizations(initialDraft);
    form.reset({ localizations: locs });
    setActiveLanguage(initialActiveLanguage(locs));
  }, [initialDraft]);

  const activeIndex = Math.max(
    0,
    localizationFields.findIndex(f => f.language === activeLanguage)
  );
  const canRemoveLanguage = localizationFields.length > 1;
  const availableToAdd = FORM_LANGUAGES.filter(
    lang => !localizationFields.some(f => f.language === lang)
  );

  const publishMutation = usePublishNotification(environmentId);
  const saveDraftMutation = useSaveDraft(environmentId);
  const updateMutation = useUpdateNotification(environmentId);

  // Editing an existing draft updates it in place; otherwise a new one is made.
  // Edit mode is driven by whether `initialDraft` is set.
  const isEditing = !!initialDraft;
  const editingId = initialDraft?.id;

  // Mirrors `editingId` in a ref so pending mutation callbacks (defined in a
  // closure over the `editingId` from whenever they were submitted) can read
  // the *current* editing target when they resolve, not the one they closed
  // over. If the user switches to editing a different draft before an
  // in-flight mutation resolves, its `onSuccess`/`onError` must not touch the
  // form or report a result for content that's no longer on screen.
  const editingIdRef = useRef(editingId);
  editingIdRef.current = editingId;

  const {
    formState: { isValid }
  } = form;
  // Publishing requires every localization to be complete; saving a draft
  // does not, since a draft is meant to hold work in progress.
  const canPublish = !disabled && isValid;
  const canSaveDraft = !disabled;
  // When editing a draft, both "Publish" and "Update draft" submit through the
  // same `updateMutation`, so its `isPending` alone can't tell which button
  // triggered it. Track the in-flight action to light up only that button.
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
    // Leave edit mode: parent drops the draft being edited so "Update draft"
    // reverts to "Save draft".
    onClear?.();
  };

  const handlePublish = form.handleSubmit(() => {
    const payload = toInput(NotificationStatus.PUBLISHED);
    // Captured now, compared against `editingIdRef.current` when the mutation
    // resolves, so switching drafts mid-flight is detected.
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
    // Publishing an edited draft promotes the same record; otherwise create.
    if (isEditing && editingId) {
      updateMutation.mutate({ id: editingId, input: payload }, onDone);
    } else {
      publishMutation.mutate(payload, onDone);
    }
  });

  const handleSaveDraft = () => {
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
    // Updating an existing draft keeps its id; otherwise create a new one.
    if (isEditing && editingId) {
      updateMutation.mutate({ id: editingId, input: payload }, onDone);
    } else {
      saveDraftMutation.mutate(payload, onDone);
    }
  };

  // Clears the form back to a fresh "new notification" state, discarding any
  // draft that was open for editing (so "Update draft" reverts to "Save draft").
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

        {/* Fields for the active language, addressed by its field-array index. */}
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
            disabled={!canPublish}
            loading={isPublishPending}
          >
            {t('publish')}
          </Button>
          <Button
            type="button"
            variant="secondary"
            onClick={handleSaveDraft}
            disabled={!canSaveDraft}
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
