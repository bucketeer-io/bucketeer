import { useCallback, useMemo, useState } from 'react';
import { OrganizationCreatorCommand } from '@api/organization';
import { organizationUpdate } from '@api/organization/organization-update';
import { useQueryOrganizationDetails } from '@queries/organization-details';
import CommonForm, { FormFieldProps } from 'containers/common-form';
import { onSubmitSuccess } from 'containers/pages/organizations';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { ContentDetailsProps } from '../page-content';

export const Setting = ({ organizationId }: ContentDetailsProps) => {
  const { data } = useQueryOrganizationDetails({
    params: {
      id: organizationId!
    },
    enabled: !!organizationId
  });
  const { notify } = useToast();
  const { t } = useTranslation(['common', 'form']);

  const [isLoadingSubmit, setIsLoadingSubmit] = useState(false);

  const formFields: FormFieldProps[] = useMemo(
    () => [
      {
        name: 'name',
        label: `${t('name')}`,
        placeholder: `${t('form:placeholder-name')}`,
        isRequired: true,
        isOptional: false,
        isExpand: true,
        fieldType: 'input'
      },
      {
        name: 'urlCode',
        label: `${t('form:url-code')}`,
        placeholder: `${t('form:placeholder-code')}`,
        isRequired: true,
        isOptional: false,
        isExpand: false,
        fieldType: 'input',
        isDisabled: true
      },
      {
        name: 'description',
        label: `${t('form:description')}`,
        placeholder: `${t('form:placeholder-desc')}`,
        isRequired: false,
        isOptional: true,
        isExpand: true,
        fieldType: 'textarea'
      },
      {
        name: 'owner',
        label: `${t('form:owner')}`,
        placeholder: `${t('form:placeholder-email')}`,
        isRequired: true,
        isExpand: true,
        fieldType: 'input'
      }
    ],
    []
  );
  const formSchema = useMemo(
    () =>
      yup.object().shape({
        name: yup.string().required(),
        urlCode: yup.string().required()
      }),
    []
  );

  const handleOnSubmit = useCallback(
    (formValues: OrganizationCreatorCommand) => {
      setIsLoadingSubmit(true);
      organizationUpdate({
        id: organizationId as string,
        changeDescriptionCommand: {
          description: formValues.description
        },
        renameCommand: {
          name: formValues.name
        },
        changeOwnerEmailCommand: {
          ownerEmail: formValues.ownerEmail
        }
      })
        .then(() => {
          onSubmitSuccess({
            name: formValues.name,
            notify,
            submitType: 'updated'
          });
        })
        .finally(() => setIsLoadingSubmit(false));
    },
    []
  );
  return (
    <CommonForm
      title={t('form:general-info')}
      isLoading={isLoadingSubmit}
      formFields={formFields}
      formData={data?.organization}
      formSchema={formSchema}
      onSubmit={formValues =>
        handleOnSubmit(formValues as OrganizationCreatorCommand)
      }
    />
  );
};
