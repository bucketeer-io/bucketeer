import { useMemo } from 'react';
import { Trans } from 'react-i18next';
import { useTranslation } from 'i18n';
import { RepositoryType } from '@types';
import Dropdown, { DropdownOption } from 'components/dropdown';
import { CodeRefFilters } from '../types';

const FiltersBar = ({
  repositoryOptions,
  branchOptions,
  extensionOptions,
  filters,
  onChangeFilters
}: {
  repositoryOptions: DropdownOption[];
  branchOptions: DropdownOption[];
  extensionOptions: DropdownOption[];
  filters: CodeRefFilters;
  onChangeFilters: (filters: Partial<CodeRefFilters>) => void;
}) => {
  const { t } = useTranslation(['table']);
  const repositoryLabel = useMemo(
    () =>
      repositoryOptions.find(item => item.value === filters?.repositoryType)
        ?.label || '',
    [filters, repositoryOptions]
  );
  const branchLabel = useMemo(
    () =>
      branchOptions.find(item => item.value === filters?.repositoryBranch)
        ?.label || '',
    [filters, branchOptions]
  );
  const extensionLabel = useMemo(
    () =>
      extensionOptions.find(item => item.value === filters?.fileExtension)
        ?.label || '',
    [filters, extensionOptions]
  );
  const repositoryLabelText = useMemo(
    () => (
      <Trans
        i18nKey={'table:code-refs.repository-name'}
        values={{
          name: repositoryLabel ? repositoryLabel : t('table:code-refs.all')
        }}
        components={{
          text: (
            <span
              className={repositoryLabel ? 'text-gray-600' : 'text-gray-500'}
            />
          )
        }}
      />
    ),
    [repositoryLabel, t]
  );

  const branchLabelText = useMemo(
    () => (
      <Trans
        i18nKey={'table:code-refs.branch-name'}
        values={{
          name: branchLabel ? branchLabel : t('table:code-refs.all')
        }}
        components={{
          text: (
            <span className={branchLabel ? 'text-gray-600' : 'text-gray-500'} />
          )
        }}
      />
    ),
    [branchLabel, t]
  );

  const extensionLabelText = useMemo(
    () => (
      <Trans
        i18nKey={'table:code-refs.file-extension-name'}
        values={{
          name: extensionLabel ? extensionLabel : t('table:code-refs.all')
        }}
        components={{
          text: (
            <span
              className={extensionLabel ? 'text-gray-600' : 'text-gray-500'}
            />
          )
        }}
      />
    ),
    [extensionLabel, t]
  );

  return (
    <div className="flex items-center justify-between w-full gap-x-10">
      <div className="flex flex-col gap-y-4">
        <p className="typo-head-bold-small text-gray-800">
          {t('code-refs.title')}
        </p>
        <p className="typo-para-medium text-gray-500">
          {t('code-refs.description')}
        </p>
      </div>

      <div className="flex items-center gap-x-4">
        <Dropdown
          labelCustom={repositoryLabelText}
          options={repositoryOptions}
          value={filters?.repositoryType || 'all'}
          onChange={value =>
            onChangeFilters({
              repositoryType:
                value === 'all' ? undefined : (value as RepositoryType)
            })
          }
          isTruncate={false}
        />

        <Dropdown
          labelCustom={branchLabelText}
          options={branchOptions}
          value={filters?.repositoryBranch || 'all'}
          onChange={value =>
            onChangeFilters({
              repositoryBranch: value === 'all' ? undefined : (value as string)
            })
          }
          isTruncate={false}
        />

        <Dropdown
          labelCustom={extensionLabelText}
          options={extensionOptions}
          value={filters?.fileExtension || 'all'}
          isTruncate={false}
          onChange={value =>
            onChangeFilters({
              fileExtension: value === 'all' ? undefined : (value as string)
            })
          }
        />
      </div>
    </div>
  );
};

export default FiltersBar;
