import { useMemo } from 'react';
import { Trans } from 'react-i18next';
import { useTranslation } from 'i18n';
import { RepositoryType } from '@types';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
  DropdownOption
} from 'components/dropdown';
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
        <DropdownMenu>
          <DropdownMenuTrigger
            label={
              <Trans
                i18nKey={'table:code-refs.repository-name'}
                values={{
                  name: repositoryLabel
                    ? repositoryLabel
                    : t('table:code-refs.all')
                }}
                components={{
                  text: (
                    <span
                      className={
                        repositoryLabel ? 'text-gray-600' : 'text-gray-500'
                      }
                    />
                  )
                }}
              />
            }
          />
          <DropdownMenuContent>
            {repositoryOptions.map((item, index) => (
              <DropdownMenuItem
                key={index}
                label={item.label}
                value={item.value}
                onSelectOption={value =>
                  onChangeFilters({
                    repositoryType:
                      value === 'all' ? undefined : (value as RepositoryType)
                  })
                }
              />
            ))}
          </DropdownMenuContent>
        </DropdownMenu>
        <DropdownMenu>
          <DropdownMenuTrigger
            label={
              <Trans
                i18nKey={'table:code-refs.branch-name'}
                values={{
                  name: branchLabel ? branchLabel : t('table:code-refs.all')
                }}
                components={{
                  text: (
                    <span
                      className={
                        branchLabel ? 'text-gray-600' : 'text-gray-500'
                      }
                    />
                  )
                }}
              />
            }
          />
          <DropdownMenuContent>
            {branchOptions.map((item, index) => (
              <DropdownMenuItem
                key={index}
                label={item.label}
                value={item.value}
                onSelectOption={value =>
                  onChangeFilters({
                    repositoryBranch:
                      value === 'all' ? undefined : (value as string)
                  })
                }
              />
            ))}
          </DropdownMenuContent>
        </DropdownMenu>
        <DropdownMenu>
          <DropdownMenuTrigger
            label={
              <Trans
                i18nKey={'table:code-refs.file-extension-name'}
                values={{
                  name: extensionLabel
                    ? extensionLabel
                    : t('table:code-refs.all')
                }}
                components={{
                  text: (
                    <span
                      className={
                        extensionLabel ? 'text-gray-600' : 'text-gray-500'
                      }
                    />
                  )
                }}
              />
            }
          />
          <DropdownMenuContent>
            {extensionOptions.map((item, index) => (
              <DropdownMenuItem
                key={index}
                label={item.label}
                value={item.value}
                onSelectOption={value =>
                  onChangeFilters({
                    fileExtension:
                      value === 'all' ? undefined : (value as string)
                  })
                }
              />
            ))}
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </div>
  );
};

export default FiltersBar;
