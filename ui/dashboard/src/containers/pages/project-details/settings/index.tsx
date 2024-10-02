import { useMemo, useState } from 'react';
import { IconLaunchOutlined } from 'react-icons-material-design';
import { AccountsFetcherParams } from '@api/account/accounts-fetcher';
import { ProjectCreatorCommand } from '@api/project';
import { projectUpdate } from '@api/project/project-update';
import { useQueryAccounts } from '@queries/accounts';
import { LIST_PAGE_SIZE } from 'constants/app';
import CommonForm, { FormFieldProps } from 'containers/common-form';
import { onSubmitSuccess } from 'containers/pages/organizations';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { Project } from '@types';
import Button from 'components/button';
import Icon from 'components/icon';

export const Settings = ({ projectData }: { projectData?: Project }) => {
  const { t } = useTranslation(['common', 'form']);

  const [isLoadingSubmit, setIsLoadingSubmit] = useState(false);

  const accountParams: AccountsFetcherParams = {
    pageSize: LIST_PAGE_SIZE,
    cursor: String(0),
    orderBy: 'DEFAULT',
    orderDirection: 'ASC',
    searchKeyword: '',
    disabled: false,
    organizationId: projectData?.organizationId
  };
  const { notify } = useToast();
  const { data } = useQueryAccounts({
    params: accountParams,
    enabled: !!projectData?.organizationId
  });

  const formFields: FormFieldProps[] = useMemo(
    () => [
      {
        name: 'creatorEmail',
        label: `${t(`maintainer`)}`,
        placeholder: `${t(`maintainer`)}`,
        isRequired: true,
        isOptional: false,
        isExpand: true,
        fieldType: 'dropdown',
        dropdownOptions: data?.accounts?.map(account => ({
          label: `${account.name} (${account.email})`,
          value: account.email
        })),
        renderTrigger: value => {
          const account = data?.accounts?.find(
            account => account.email === value
          );
          return (
            <p>{account ? `${account?.name} (${account?.email})` : value}</p>
          );
        }
      },
      {
        name: 'name',
        label: `${t(`name`)}`,
        placeholder: `${t(`form:placeholder-name`)}`,
        isRequired: true,
        isOptional: false,
        isExpand: true,
        fieldType: 'input'
      },
      {
        name: 'description',
        label: `${t(`form:description`)}`,
        placeholder: `${t(`form:placeholder-desc`)}`,
        isRequired: false,
        isOptional: true,
        isExpand: true,
        fieldType: 'textarea'
      }
    ],
    [data, projectData]
  );

  const formSchema = yup.object().shape({
    creatorEmail: yup.string().required(),
    name: yup.string().required()
  });

  const handleOnClickSubmitBtn = () => {
    document.getElementById('common-submit-btn')?.click();
  };

  const handleOnSubmit = (formValues: ProjectCreatorCommand) => {
    setIsLoadingSubmit(true);
    projectUpdate({
      id: formValues?.id || '',
      changeDescriptionCommand: {
        description: formValues.description
      },
      renameCommand: {
        name: formValues.name
      }
    })
      .then(() => {
        onSubmitSuccess({
          name: formValues.name,
          submitType: 'updated',
          notify
        });
      })
      .finally(() => setIsLoadingSubmit(false));
  };

  return (
    <div className="flex flex-col gap-y-6 w-full">
      <div className="flex lg:items-center justify-between flex-col lg:flex-row">
        <div className="w-fit">
          <p className="typo-head-bold-small">{t(`settings`)}</p>
        </div>
        <div className="flex items-center gap-4 mt-3 lg:mt-0">
          <Button variant="text" className="flex-1 lg:flex-none">
            <Icon icon={IconLaunchOutlined} size="sm" />
            {t('documentation')}
          </Button>
          <Button
            type="submit"
            className="w-[120px]"
            loading={isLoadingSubmit}
            onClick={handleOnClickSubmitBtn}
          >
            {t('save')}
          </Button>
        </div>
      </div>
      <CommonForm
        formFields={formFields}
        formClassName="gap-6 mt-0"
        isShowSubmitBtn={false}
        formData={projectData}
        formSchema={formSchema}
        onSubmit={formValues =>
          handleOnSubmit(formValues as ProjectCreatorCommand)
        }
      />
    </div>
  );
};
