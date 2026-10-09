import { useCallback, useEffect } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { useNavigate } from 'react-router';
import { featureBulkClone } from '@api/features';
import { yupResolver } from '@hookform/resolvers/yup';
import { useQueryFeature } from '@queries/feature-details';
import {
  getCurrentEnvironment,
  getEditorEnvironments,
  hasEditable,
  useAuth
} from 'auth';
import {
  PAGE_PATH_FEATURE_TARGETING,
  PAGE_PATH_FEATURES
} from 'constants/routing';
import { useToast } from 'hooks';
import useFormSchema, { FormSchemaProps } from 'hooks/use-form-schema';
import { useTranslation } from 'i18n';
import { setCurrentEnvIdStorage } from 'storage/environment';
import { setCurrentProjectEnvironmentStorage } from 'storage/project-environment';
import * as yup from 'yup';
import { checkEnvironmentEmptyId, onFormatEnvironments } from 'utils/function';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Dropdown, { DropdownValue } from 'components/dropdown';
import Form from 'components/form';
import Input from 'components/input';
import SlideModal from 'components/modal/slide';
import EnvironmentEditorList from 'elements/environment-editor-list';
import FormLoading from 'elements/form-loading';

interface CloneFlagModalProps {
  flagId: string;
  isOpen: boolean;
  onClose: () => void;
}

export interface CloneFlagForm {
  id: string;
  name: string;
  originEnvironmentId: string;
  destinationEnvironmentIds: string[];
}

const formSchema = ({ requiredMessage }: FormSchemaProps) =>
  yup.object().shape({
    id: yup.string().required(requiredMessage),
    name: yup.string().required(requiredMessage),
    originEnvironmentId: yup.string().required(requiredMessage),
    destinationEnvironmentIds: yup
      .array()
      .of(yup.string().required())
      .min(1, requiredMessage)
      .required(requiredMessage)
  });

const CloneFlagModal = ({ flagId, isOpen, onClose }: CloneFlagModalProps) => {
  const { consoleAccount } = useAuth();
  const { editorEnvironments, projects } = getEditorEnvironments(
    consoleAccount!
  );
  const { emptyEnvironmentId, formattedEnvironments } =
    onFormatEnvironments(editorEnvironments);
  const navigate = useNavigate();
  const { t } = useTranslation(['common', 'form', 'message']);
  const { notify, errorNotify } = useToast();

  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const editable = hasEditable(consoleAccount!);

  const {
    data: featureCollection,
    isLoading: isLoadingFeature,
    error: featureError
  } = useQueryFeature({
    params: {
      id: flagId as string,
      environmentId: currentEnvironment?.id
    },
    enabled: !!flagId
  });

  const feature = featureCollection?.feature;

  const form = useForm({
    resolver: yupResolver(useFormSchema(formSchema)),
    values: {
      id: feature?.id || '',
      name: feature?.name || '',
      originEnvironmentId: currentEnvironment?.id || emptyEnvironmentId || '',
      destinationEnvironmentIds: [] as string[]
    }
  });

  const getCurrentLabelEnv = useCallback(
    (environmentId: string) => {
      const currentEnv = formattedEnvironments.find(
        item => item.id === environmentId
      );

      return `${currentEnv?.name} (${t('common:source-type.project')}: ${projects.find(project => project.id === currentEnv?.projectId)?.name})`;
    },
    [formattedEnvironments, projects]
  );

  const getEnvName = useCallback(
    (environmentId: string) => {
      const env = formattedEnvironments.find(item => item.id === environmentId);
      return env?.name || environmentId;
    },
    [formattedEnvironments]
  );

  const onSubmit: SubmitHandler<CloneFlagForm> = useCallback(
    async values => {
      try {
        const { id, destinationEnvironmentIds, originEnvironmentId } = values;
        const resp = await featureBulkClone({
          id,
          environmentId: checkEnvironmentEmptyId(originEnvironmentId),
          targetEnvironmentIds: destinationEnvironmentIds.map(envId =>
            checkEnvironmentEmptyId(envId)
          )
        });

        if (resp) {
          const results = resp.results || [];
          const successes = results.filter(
            r => r.status === 'BULK_CLONE_FEATURE_STATUS_SUCCESS'
          );
          const alreadyExists = results.filter(
            r => r.status === 'BULK_CLONE_FEATURE_STATUS_ALREADY_EXISTS'
          );
          const failures = results.filter(
            r => r.status === 'BULK_CLONE_FEATURE_STATUS_FAILED'
          );

          if (successes.length > 0) {
            notify({
              message: t('message:collection-action-success', {
                collection: t('common:source-type.feature-flag'),
                action: t('cloned')
              })
            });
          }

          alreadyExists.forEach(r => {
            notify({
              messageType: 'warning',
              message: t('message:bulk-clone-already-exists', {
                env: getEnvName(r.environmentId)
              })
            });
          });

          failures.forEach(r => {
            notify({
              messageType: 'error',
              message: t('message:bulk-clone-failed', {
                env: getEnvName(r.environmentId)
              })
            });
          });

          onClose();
          if (
            successes.length === 1 &&
            destinationEnvironmentIds.length === 1
          ) {
            const targetEnvironment = formattedEnvironments.find(
              item => item.id === destinationEnvironmentIds[0]
            );
            if (targetEnvironment) {
              setCurrentEnvIdStorage(targetEnvironment.id);
              setCurrentProjectEnvironmentStorage({
                environmentId: targetEnvironment.id,
                projectId: targetEnvironment.projectId
              });
              navigate(
                `/${targetEnvironment.urlCode}${PAGE_PATH_FEATURES}/${id}${PAGE_PATH_FEATURE_TARGETING}`,
                { replace: true }
              );
            }
          }
        }
      } catch (error) {
        errorNotify(error);
      }
    },
    [formattedEnvironments, getEnvName]
  );

  useEffect(() => {
    if (featureError) {
      errorNotify(featureError);
    }
  }, [featureError]);

  return (
    <SlideModal
      title={t('form:feature-flags.clone-title')}
      isOpen={isOpen}
      onClose={onClose}
    >
      {isLoadingFeature ? (
        <FormLoading />
      ) : (
        <div className="w-full p-5">
          <p className="text-gray-600 dark:text-dark-gray-300 typo-para-small">
            {t('form:feature-flags.clone-desc')}
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
                        {...field}
                        placeholder={`${t('form:placeholder-name')}`}
                        disabled
                      />
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />
              <Form.Field
                control={form.control}
                name={`originEnvironmentId`}
                render={({ field }) => (
                  <Form.Item className="py-2">
                    <Form.Label required>{t('form:origin-env')}</Form.Label>
                    <Form.Control>
                      <Dropdown
                        options={formattedEnvironments.map(item => ({
                          label: item.name,
                          value: item.id
                        }))}
                        labelCustom={getCurrentLabelEnv(
                          currentEnvironment?.id || emptyEnvironmentId
                        )}
                        value={field.value}
                        onChange={field.onChange}
                        placeholder={t(`form:select-environment`)}
                        disabled
                        className="w-full"
                        contentClassName="min-w-[502px]"
                      />
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />

              <Form.Field
                control={form.control}
                name={`destinationEnvironmentIds`}
                render={({ field }) => (
                  <Form.Item className="py-2">
                    <Form.Label required>
                      {t('form:destination-envs')}
                    </Form.Label>
                    <Form.Control>
                      <EnvironmentEditorList
                        value={field.value}
                        disabled={!editable}
                        currentEnvironmentId={
                          currentEnvironment?.id || emptyEnvironmentId
                        }
                        onSelectOption={(selectedValue: DropdownValue) => {
                          const envId = String(selectedValue);
                          const current = field.value;
                          if (current.includes(envId)) {
                            field.onChange(current.filter(v => v !== envId));
                          } else {
                            field.onChange([...current, envId]);
                          }
                        }}
                      />
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />

              <div className="absolute left-0 bottom-0 bg-gray-50 dark:bg-dark-black-800 w-full rounded-b-lg">
                <ButtonBar
                  primaryButton={
                    <Button variant="secondary" onClick={onClose}>
                      {t(`cancel`)}
                    </Button>
                  }
                  secondaryButton={
                    <Button
                      type="submit"
                      disabled={!form.formState.isDirty || !editable}
                      loading={form.formState.isSubmitting}
                    >
                      {t(`clone-flag`)}
                    </Button>
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

export default CloneFlagModal;
