import { useCallback, useMemo } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { pushUpdater, TagChange } from '@api/push';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidatePushes } from '@queries/pushes';
import { useQueryClient } from '@tanstack/react-query';
import { useAuth } from 'auth';
import { useToast } from 'hooks';
import useFormSchema, { FormSchemaProps } from 'hooks/use-form-schema';
import { useTranslation } from 'i18n';
import uniqBy from 'lodash/uniqBy';
import * as yup from 'yup';
import { Push } from '@types';
import { checkEnvironmentEmptyId, onFormatEnvironments } from 'utils/function';
import { UserMessage } from 'pages/feature-flag-details/targeting/individual-rule';
import { useFetchTags } from 'pages/members/collection-loader';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import { CreatableSelect } from 'components/creatable-select';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import Input from 'components/input';
import SlideModal from 'components/modal/slide';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import FormLoading from 'elements/form-loading';

interface EditPushModalProps {
  disabled?: boolean;
  isOpen: boolean;
  isLoadingPush: boolean;
  push?: Push;
  onClose: () => void;
}

export interface EditPushForm {
  name: string;
  tags?: string[];
  environmentId: string;
}

const formSchema = ({ requiredMessage }: FormSchemaProps) =>
  yup.object().shape({
    name: yup.string().required(requiredMessage),
    tags: yup.array(),
    environmentId: yup.string().required(requiredMessage)
  });

const EditPushModal = ({
  disabled,
  isOpen,
  isLoadingPush,
  push,
  onClose
}: EditPushModalProps) => {
  const { consoleAccount } = useAuth();
  const queryClient = useQueryClient();
  const { t } = useTranslation(['common', 'form', 'message']);
  const { notify, errorNotify } = useToast();

  const { data: tagCollection, isLoading: isLoadingTags } = useFetchTags({
    environmentId: push?.environmentId || '',
    entityType: 'FEATURE_FLAG',
    options: {
      enabled: !!push
    }
  });

  const tagOptions = (uniqBy(tagCollection?.tags || [], 'name') || [])?.map(
    tag => ({
      label: tag.name,
      value: tag.name
    })
  );
  const editorEnvironments = useMemo(
    () =>
      consoleAccount?.environmentRoles
        .filter(item => item.role === 'Environment_EDITOR')
        ?.map(item => item.environment) || [],
    [consoleAccount]
  );

  const { emptyEnvironmentId, formattedEnvironments } =
    onFormatEnvironments(editorEnvironments);

  const form = useForm({
    resolver: yupResolver(useFormSchema(formSchema)),
    values: {
      name: push?.name || '',
      tags: push?.tags || [],
      environmentId: push?.environmentId || emptyEnvironmentId
    }
  });

  const {
    getValues,
    formState: { isValid, isSubmitting, isDirty }
  } = form;

  const handleCheckTags = useCallback(
    (tagValues: string[]) => {
      if (push?.tags) {
        const tagChanges: TagChange[] = [];
        const { tags } = push;
        tags?.forEach(item => {
          if (!tagValues.find(tag => tag === item)) {
            tagChanges.push({
              changeType: 'DELETE',
              tag: item
            });
          }
        });
        tagValues.forEach(item => {
          const currentTag = tags.find(tag => tag === item);
          if (!currentTag) {
            tagChanges.push({
              changeType: 'CREATE',
              tag: item
            });
          }
        });

        return tagChanges;
      }
      return [];
    },
    [push]
  );

  const onSubmit: SubmitHandler<EditPushForm> = async values => {
    const { name, tags, environmentId } = values;
    await pushUpdater({
      name,
      tagChanges: handleCheckTags(tags || []),
      id: push?.id || '',
      environmentId: checkEnvironmentEmptyId(environmentId)
    })
      .then(() => {
        notify({
          message: t('message:collection-action-success', {
            collection: t('push-notification'),
            action: t('updated')
          })
        });
        invalidatePushes(queryClient);
        onClose();
      })
      .catch(error => errorNotify(error));
  };

  return (
    <SlideModal title={t('edit-push')} isOpen={isOpen} onClose={onClose}>
      {isLoadingPush ? (
        <FormLoading />
      ) : (
        <div className="w-full p-5 pb-28">
          <div className="typo-para-small text-gray-600 mb-3">
            {t('new-push-subtitle')}
          </div>
          <p className="text-gray-800 typo-head-bold-small">
            {t('form:general-info')}
          </p>
          <FormProvider {...form}>
            <Form onSubmit={form.handleSubmit(onSubmit)}>
              <Form.Field
                control={form.control}
                name="name"
                render={({ field }) => (
                  <Form.Item>
                    <Form.Label required>{t('name')}</Form.Label>
                    <Form.Control>
                      <Input
                        placeholder={`${t('form:placeholder-name')}`}
                        disabled={disabled}
                        {...field}
                      />
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />
              <Form.Field
                control={form.control}
                name={`environmentId`}
                render={({ field }) => (
                  <Form.Item className="py-2">
                    <Form.Label required>{t('environment')}</Form.Label>
                    <Form.Control>
                      <DropdownMenu>
                        <DropdownMenuTrigger
                          placeholder={t(`form:select-environment`)}
                          label={
                            formattedEnvironments.find(
                              item => item.id === getValues('environmentId')
                            )?.name
                          }
                          disabled
                          variant="secondary"
                          className="w-full"
                        />
                        <DropdownMenuContent
                          className="w-[502px]"
                          align="start"
                          {...field}
                        >
                          {formattedEnvironments.map((item, index) => (
                            <DropdownMenuItem
                              {...field}
                              key={index}
                              value={item.id}
                              label={item.name}
                              onSelectOption={value => {
                                field.onChange(value);
                              }}
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
                control={form.control}
                name={`tags`}
                render={({ field }) => (
                  <Form.Item className="py-2">
                    <Form.Label>{t('form:feature-flag-tags')}</Form.Label>
                    <Form.Control>
                      <CreatableSelect
                        value={field.value?.map(tag => {
                          const tagItem = tagOptions.find(
                            item => item.value === tag
                          );
                          return {
                            label: tagItem?.label || tag,
                            value: tagItem?.value || tag
                          };
                        })}
                        disabled={
                          isLoadingTags || !tagOptions.length || disabled
                        }
                        loading={isLoadingTags}
                        allowCreateWhileLoading={false}
                        isValidNewOption={() => false}
                        isClearable
                        onKeyDown={e => {
                          const { value } = e.target as HTMLInputElement;
                          const isExists = tagOptions.find(
                            item =>
                              item.label
                                .toLowerCase()
                                .includes(value.toLowerCase()) &&
                              !field.value?.includes(item.label)
                          );
                          if (e.key === 'Enter' && (!isExists || !value)) {
                            e.preventDefault();
                          }
                        }}
                        placeholder={t(`form:placeholder-tags`)}
                        options={tagOptions}
                        onChange={value =>
                          field.onChange(value.map(tag => tag.value))
                        }
                        noOptionsMessage={() => (
                          <UserMessage message={t('no-options-found')} />
                        )}
                      />
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />

              <div className="absolute left-0 bottom-0 bg-gray-50 w-full rounded-b-lg">
                <ButtonBar
                  primaryButton={
                    <Button variant="secondary" onClick={onClose}>
                      {t(`cancel`)}
                    </Button>
                  }
                  secondaryButton={
                    <DisabledButtonTooltip
                      align="center"
                      hidden={!disabled}
                      trigger={
                        <Button
                          type="submit"
                          disabled={!isValid || !isDirty || disabled}
                          loading={isSubmitting}
                        >
                          {t(`submit`)}
                        </Button>
                      }
                    />
                  }
                />
              </div>
            </Form>
          </FormProvider>
        </div>
      )}
    </SlideModal>
  );
};

export default EditPushModal;
