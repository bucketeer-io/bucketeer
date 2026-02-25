import { forwardRef, Ref, useCallback, useMemo } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { triggerCreator } from '@api/trigger';
import { triggerUpdate } from '@api/trigger/triggers-update';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateTriggers } from '@queries/triggers';
import { useQueryClient } from '@tanstack/react-query';
import { useToast } from 'hooks';
import useFormSchema from 'hooks/use-form-schema';
import { useUnsavedLeavePage } from 'hooks/use-unsaved-leave-page';
import { useTranslation } from 'i18n';
import { TriggerActionType, TriggerItemType, TriggerType } from '@types';
import { IconWebhook } from '@icons';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Dropdown from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import TextArea from 'components/textarea';
import FeatureFlagStatus from 'elements/feature-flag-status';
import { CreateTriggerSchema, formSchema } from './form-schema';

const CreateTriggerForm = forwardRef(
  (
    {
      disabled,
      featureId,
      environmentId,
      selectedTrigger,
      onCancel,
      setTriggerNewlyCreated
    }: {
      disabled: boolean;
      featureId: string;
      environmentId: string;
      selectedTrigger?: TriggerItemType;
      onCancel: () => void;
      setTriggerNewlyCreated: (trigger: TriggerItemType | undefined) => void;
    },
    ref: Ref<HTMLDivElement>
  ) => {
    const { t } = useTranslation(['table', 'form', 'common', 'message']);
    const queryClient = useQueryClient();

    const { notify, errorNotify } = useToast();

    const triggerTypeOptions = useMemo(
      () => [
        {
          label: (
            <div className="flex items-center w-full gap-x-2">
              <Icon icon={IconWebhook} />
              <p className="text-gray-600 typo-para-medium truncate">{`${t('trigger.dropdown-desc')}`}</p>
            </div>
          ),
          value: TriggerType.WEBHOOK
        }
      ],
      []
    );

    const createLabel = useCallback(
      (enabled: boolean) => (
        <div className="flex items-center gap-x-2">
          <Trans
            i18nKey="table:trigger.turn-flag-status"
            components={{
              status: (
                <FeatureFlagStatus
                  status={t(`form:experiments.${enabled ? 'on' : 'off'}`)}
                  enabled={enabled}
                />
              )
            }}
          />
        </div>
      ),
      []
    );

    const triggerActionOptions = useMemo(
      () => [
        {
          label: createLabel(true),
          value: TriggerActionType.ON
        },
        {
          label: createLabel(false),
          value: TriggerActionType.OFF
        }
      ],
      []
    );

    const form = useForm<CreateTriggerSchema>({
      resolver: yupResolver(useFormSchema(formSchema)),
      defaultValues: {
        type: selectedTrigger?.flagTrigger?.type || TriggerType.WEBHOOK,
        action: selectedTrigger?.flagTrigger?.action || undefined,
        description: selectedTrigger?.flagTrigger?.description || ''
      }
    });

    const {
      control,
      formState: { isDirty, isValid, isSubmitting }
    } = form;
    const isEdit = useMemo(() => !!selectedTrigger, [selectedTrigger]);

    const onSubmit = useCallback(
      async (values: CreateTriggerSchema) => {
        if (!disabled) {
          try {
            let resp: TriggerItemType;
            if (isEdit && selectedTrigger) {
              resp = await triggerUpdate({
                id: selectedTrigger.flagTrigger.id,
                environmentId,
                description: values.description
              });
            } else {
              resp = await triggerCreator({
                ...values,
                featureId,
                environmentId
              });
            }

            if (resp) {
              invalidateTriggers(queryClient);
              notify({
                message: t('message:collection-action-success', {
                  collection: t('feature-flags.trigger'),
                  action: t(isEdit ? 'common:updated' : 'common:created')
                })
              });
              setTriggerNewlyCreated(isEdit ? undefined : resp);
              onCancel();
            }
          } catch (error) {
            errorNotify(error);
          }
        }
      },
      [featureId, environmentId, selectedTrigger, isEdit, disabled]
    );

    useUnsavedLeavePage({ isShow: isDirty && !isSubmitting });
    return (
      <div ref={ref} className="w-full  p-6 border border-gray-400 rounded-lg">
        <FormProvider {...form}>
          <Form
            onSubmit={form.handleSubmit(onSubmit)}
            className="flex flex-col w-full gap-y-6"
          >
            <Form.Field
              control={control}
              name="type"
              render={({ field }) => (
                <Form.Item className="py-0 w-full">
                  <Form.Label required className="!mb-2">
                    {t('trigger.trigger-type')}
                  </Form.Label>
                  <Form.Control>
                    <Dropdown
                      disabled={isEdit || disabled}
                      options={triggerTypeOptions}
                      value={field.value}
                      onChange={field.onChange}
                      contentClassName="max-w-[300px] sm:w-full"
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />

            <Form.Field
              control={control}
              name="action"
              render={({ field }) => (
                <Form.Item className="py-0 w-full">
                  <Form.Label required className="!mb-2">
                    {t('trigger.action')}
                  </Form.Label>
                  <Form.Control>
                    <Dropdown
                      options={triggerActionOptions}
                      value={field.value}
                      onChange={field.onChange}
                      disabled={isEdit || disabled}
                      placeholder={t('trigger.select-action')}
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <Form.Field
              control={control}
              name="description"
              render={({ field }) => (
                <Form.Item className="py-0 w-full">
                  <Form.Label optional className="!mb-2">
                    {t('form:description')}
                  </Form.Label>
                  <Form.Control>
                    <TextArea
                      {...field}
                      placeholder={t('form:placeholder-desc')}
                      rows={2}
                      disabled={disabled}
                      style={{
                        resize: 'vertical',
                        maxHeight: 98
                      }}
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />

            <ButtonBar
              primaryButton={
                <Button type="button" variant="secondary" onClick={onCancel}>
                  {t('common:cancel')}
                </Button>
              }
              secondaryButton={
                <Button
                  disabled={!isDirty || !isValid || disabled}
                  loading={isSubmitting}
                >
                  {t('trigger.save-trigger')}
                </Button>
              }
              className="w-fit p-0 border-0"
            />
          </Form>
        </FormProvider>
      </div>
    );
  }
);

export default CreateTriggerForm;
