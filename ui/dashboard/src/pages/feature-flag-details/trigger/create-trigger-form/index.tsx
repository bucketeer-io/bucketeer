import { useCallback, useMemo } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { triggerCreator } from '@api/trigger';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateTriggers } from '@queries/triggers';
import { useQueryClient } from '@tanstack/react-query';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import { TriggerActionType, TriggerItemType, TriggerType } from '@types';
import { IconInfo, IconWebhook } from '@icons';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import TextArea from 'components/textarea';
import { Tooltip } from 'components/tooltip';
import FeatureFlagStatus from 'elements/feature-flag-status';
import { CreateTriggerSchema, formSchema } from './form-schema';

const CreateTriggerForm = ({
  featureId,
  environmentId,
  onCancel,
  setTriggerNewlyCreated
}: {
  featureId: string;
  environmentId: string;
  onCancel: () => void;
  setTriggerNewlyCreated: (trigger: TriggerItemType) => void;
}) => {
  const { t } = useTranslation(['table', 'form', 'common', 'message']);
  const queryClient = useQueryClient();

  const { notify, errorNotify } = useToast();
  const triggerTypeOptions = [
    {
      label: t('trigger.webhook'),
      value: TriggerType.WEBHOOK,
      icon: IconWebhook
    }
  ];

  const triggerActionOptions = [
    {
      label: (
        <div className="flex items-center gap-x-2">
          <Trans
            i18nKey="table:trigger.turn-flag-status"
            components={{
              status: (
                <FeatureFlagStatus
                  status={t('form:experiments.on')}
                  enabled={true}
                />
              )
            }}
          />
        </div>
      ),
      value: TriggerActionType.ON
    },
    {
      label: (
        <div className="flex items-center gap-x-2">
          <Trans
            i18nKey="table:trigger.turn-flag-status"
            components={{
              status: (
                <FeatureFlagStatus
                  status={t('form:experiments.off')}
                  enabled={false}
                />
              )
            }}
          />
        </div>
      ),
      value: TriggerActionType.OFF
    }
  ];

  const form = useForm<CreateTriggerSchema>({
    resolver: yupResolver(formSchema),
    defaultValues: {
      type: TriggerType.WEBHOOK,
      action: undefined,
      description: ''
    }
  });

  const {
    control,
    formState: { isDirty, isValid, isSubmitting },
    watch
  } = form;

  const triggerTypeWatch = watch('type');
  const triggerActionWatch = watch('action');

  const currentTypeOption = useMemo(
    () => triggerTypeOptions.find(item => item.value === triggerTypeWatch),
    [triggerTypeWatch]
  );

  const currentActionOption = useMemo(
    () => triggerActionOptions.find(item => item.value === triggerActionWatch),
    [triggerActionWatch]
  );

  const onSubmit = useCallback(
    async (values: CreateTriggerSchema) => {
      try {
        const resp = await triggerCreator({
          ...values,
          featureId,
          environmentId
        });

        if (resp) {
          invalidateTriggers(queryClient);
          notify({
            message: t('message:trigger-created')
          });
          setTriggerNewlyCreated(resp);
          onCancel();
        }
      } catch (error) {
        errorNotify(error);
      }
    },
    [featureId, environmentId]
  );

  return (
    <div className="w-full  p-6 border border-gray-400 rounded-lg">
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
                <Form.Label required className="relative w-fit !mb-2">
                  {t('trigger.trigger-type')}
                  <Tooltip
                    content={'tooltip'}
                    trigger={
                      <div className="flex-center size-fit absolute top-0 -right-6">
                        <Icon icon={IconInfo} size="xs" color="gray-500" />
                      </div>
                    }
                    className="max-w-[400px]"
                  />
                </Form.Label>
                <Form.Control>
                  <DropdownMenu>
                    <DropdownMenuTrigger
                      trigger={
                        <div className="flex items-center w-full gap-x-2">
                          {currentTypeOption?.icon && (
                            <Icon icon={currentTypeOption?.icon} />
                          )}
                          <p className="text-gray-600 typo-para-medium">
                            {currentTypeOption?.label}
                          </p>
                        </div>
                      }
                      isExpand
                    />
                    <DropdownMenuContent align="start">
                      {triggerTypeOptions.map((item, index) => (
                        <DropdownMenuItem
                          key={index}
                          label={item.label}
                          value={item.value}
                          icon={item.icon}
                          onSelectOption={field.onChange}
                        />
                      ))}
                    </DropdownMenuContent>
                  </DropdownMenu>
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
                <Form.Label required className="relative w-fit !mb-2">
                  {t('trigger.action')}
                  <Tooltip
                    content={'tooltip'}
                    trigger={
                      <div className="flex-center size-fit absolute top-0 -right-6">
                        <Icon icon={IconInfo} size="xs" color="gray-500" />
                      </div>
                    }
                    className="max-w-[400px]"
                  />
                </Form.Label>
                <Form.Control>
                  <DropdownMenu>
                    <DropdownMenuTrigger
                      placeholder={t('trigger.select-action')}
                      label={currentActionOption?.label}
                      isExpand
                    />
                    <DropdownMenuContent align="start">
                      {triggerActionOptions.map((item, index) => (
                        <DropdownMenuItem
                          key={index}
                          label={item.label}
                          value={item.value}
                          onSelectOption={field.onChange}
                        />
                      ))}
                    </DropdownMenuContent>
                  </DropdownMenu>
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
              <Button variant="secondary" onClick={onCancel}>
                {t('common:cancel')}
              </Button>
            }
            secondaryButton={
              <Button disabled={!isDirty || !isValid} loading={isSubmitting}>
                {t('trigger.save-trigger')}
              </Button>
            }
            className="w-fit p-0 border-0"
          />
        </Form>
      </FormProvider>
    </div>
  );
};

export default CreateTriggerForm;
