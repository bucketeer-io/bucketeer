import { useCallback } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { useNavigate } from 'react-router-dom';
import { featureCreator } from '@api/features';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateFeatures } from '@queries/features';
import { invalidateTags, useQueryTags } from '@queries/tags';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_FEATURES } from 'constants/routing';
import { useToast } from 'hooks';
import useFormSchema from 'hooks/use-form-schema';
import { useTranslation } from 'i18n';
import { v4 as uuid } from 'uuid';
import { FeatureVariation } from '@types';
import { checkEnvironmentEmptyId } from 'utils/function';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Form from 'components/form';
import { createFlagFormSchema, FlagFormSchema } from '../form-schema';
import { FlagSwitchVariationType } from '../types';
import FlagVariations from './flag-variations';
import GeneralInfo from './general-info';

const defaultVariations: FeatureVariation[] = [
  {
    id: uuid(),
    value: 'true',
    name: '',
    description: ''
  },
  {
    id: uuid(),
    value: 'false',
    name: '',
    description: ''
  }
];

const FlagForm = () => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const queryClient = useQueryClient();
  const { t } = useTranslation(['common', 'form', 'message']);
  const { notify, errorNotify } = useToast();
  const navigate = useNavigate();
  const formSchema = useFormSchema(createFlagFormSchema);
  const { data: collection } = useQueryTags({
    params: {
      cursor: String(0),
      environmentId: checkEnvironmentEmptyId(currentEnvironment?.id),
      entityType: 'FEATURE_FLAG'
    }
  });

  const tags = collection?.tags || [];

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      name: '',
      flagId: '',
      description: '',
      switchVariationType: FlagSwitchVariationType.CUSTOM,
      variationType: 'BOOLEAN',
      tags: [],
      variations: defaultVariations,
      defaultOnVariation: defaultVariations[0].id,
      defaultOffVariation: defaultVariations[1].id
    },
    mode: 'onChange'
  });

  const {
    formState: { isDirty, isValid, isSubmitting }
  } = form;

  const onBack = useCallback(
    () => navigate(`/${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}`),
    [currentEnvironment]
  );

  const onSubmit: SubmitHandler<FlagFormSchema> = async values => {
    try {
      const {
        flagId,
        name,
        tags,
        variationType,
        defaultOffVariation,
        defaultOnVariation,
        description,
        variations
      } = values;
      const resp = await featureCreator({
        environmentId: checkEnvironmentEmptyId(currentEnvironment.id),
        id: flagId,
        name,
        tags,
        defaultOnVariationIndex: variations.findIndex(
          item => item.id === defaultOnVariation
        ),
        defaultOffVariationIndex: variations.findIndex(
          item => item.id === defaultOffVariation
        ),
        variations,
        variationType,
        description
      });
      if (resp) {
        notify({
          message: t('message:collection-action-success', {
            collection: t('source-type.feature-flag'),
            action: t('created')
          })
        });
        invalidateFeatures(queryClient);
        invalidateTags(queryClient);
        onBack();
      }
    } catch (error) {
      errorNotify(error);
    }
  };

  return (
    <FormProvider {...form}>
      <Form
        onSubmit={form.handleSubmit(onSubmit)}
        className="flex flex-col size-full gap-y-6"
      >
        <div className="flex flex-col w-full gap-y-4 relative">
          <GeneralInfo tags={tags} />
          <FlagVariations />
        </div>
        <ButtonBar
          className="!border-0 p-0"
          primaryButton={
            <Button variant="secondary" onClick={onBack}>
              {t(`cancel`)}
            </Button>
          }
          secondaryButton={
            <Button
              type="submit"
              disabled={!isDirty || !isValid}
              loading={isSubmitting}
            >
              {t(`create-flag`)}
            </Button>
          }
        />
      </Form>
    </FormProvider>
  );
};

export default FlagForm;
