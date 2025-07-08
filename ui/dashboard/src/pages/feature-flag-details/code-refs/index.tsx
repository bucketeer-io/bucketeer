import { useEffect, useRef, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useQueryAPIKeys } from '@queries/api-keys';
import { useQueryCodeRefs } from '@queries/code-refs';
import { getCurrentEnvironment, useAuth } from 'auth';
import { LIST_PAGE_SIZE } from 'constants/app';
import { PAGE_PATH_APIKEYS } from 'constants/routing';
import { usePartialState } from 'hooks';
import { pickBy } from 'lodash';
import { CodeReference, Feature, RepositoryType } from '@types';
import { isEmptyObject, isNotEmpty } from 'utils/data-type';
import { checkEnvironmentEmptyId } from 'utils/function';
import { useSearchParams } from 'utils/search-params';
import { IconBitbucket, IconGithub, IconGitLab } from '@icons';
import { DropdownOption } from 'components/dropdown';
import Pagination from 'components/pagination';
import FormLoading from 'elements/form-loading';
import PageLayout from 'elements/page-layout';
import CodeAccordion from './code-accordion';
import EnableCodeRefs from './empty-collection';
import FiltersBar from './filters-bar';
import { CodeRefFilters } from './types';

export const repositoryTypeMap = {
  [RepositoryType.GITHUB]: {
    label: 'Github',
    icon: IconGithub
  },
  [RepositoryType.GITLAB]: {
    label: 'Gitlab',
    icon: IconGitLab
  },
  [RepositoryType.BITBUCKET]: {
    label: 'Bitbucket',
    icon: IconBitbucket
  }
};

const CodeReferencesPage = ({ feature }: { feature: Feature }) => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const navigate = useNavigate();
  const { searchOptions, onChangSearchParams } = useSearchParams();
  const searchFilters: Partial<CodeRefFilters> = searchOptions;
  const initLoadedRef = useRef(true);

  const defaultFilters = {
    cursor: String(0),
    pageSize: LIST_PAGE_SIZE,
    repositoryBranch: undefined,
    repositoryType: undefined,
    fileExtension: undefined,
    ...searchFilters
  } as CodeRefFilters;

  const [filters, setFilters] = usePartialState<CodeRefFilters>(defaultFilters);
  const [repositoryOptions, setRepositoryOptions] = useState<DropdownOption[]>(
    []
  );
  const [branchOptions, setBranchOptions] = useState<DropdownOption[]>([]);
  const [fileExtensionOptions, setFileExtensionOptions] = useState<
    DropdownOption[]
  >([]);
  const [hasValidApiKey, setHasValidApiKey] = useState<boolean>(false);

  const { data: apiCollection, isLoading: isLoadingApiKey } = useQueryAPIKeys({
    params: {
      cursor: String(0),
      environmentIds: [currentEnvironment.id],
      organizationId: currentEnvironment.organizationId
    }
  });
  const apiKeys = apiCollection?.apiKeys || [];

  const { data: codeRefCollection, isLoading: isLoadingCodeRefs } =
    useQueryCodeRefs({
      params: {
        environmentId: checkEnvironmentEmptyId(currentEnvironment.id),
        featureId: feature.id,
        ...filters
      },
      enabled: hasValidApiKey,
      gcTime: 0
    });

  const codeReferences = codeRefCollection?.codeReferences || [];
  const totalCount = codeRefCollection
    ? Number(codeRefCollection?.totalCount)
    : 0;
  const { repositoryBranch, repositoryType, fileExtension } = filters || {};
  const isEmpty =
    !hasValidApiKey ||
    (!repositoryType &&
      !repositoryBranch &&
      !fileExtension &&
      !isLoadingCodeRefs &&
      !codeReferences.length);

  const onChangeFilters = (values: Partial<CodeRefFilters>) => {
    const options = pickBy({ ...filters, ...values }, v => {
      return isNotEmpty(v);
    });
    onChangSearchParams(options);
    setFilters({ ...values });
  };

  useEffect(() => {
    if (isEmptyObject(searchOptions)) {
      onChangeFilters({ ...defaultFilters });
    }
  }, [searchOptions]);

  useEffect(() => {
    if (codeReferences.length > 0) {
      initLoadedRef.current = false;
      const allOption = { label: 'All', value: 'all' };

      const getUniqueOptions = (field: keyof CodeReference) => {
        const uniqueOptions = [
          ...new Set(codeReferences.map(codeRef => codeRef[field]))
        ];
        return uniqueOptions as string[];
      };
      if (!branchOptions.length) {
        const uniqueBranches = getUniqueOptions('repositoryBranch');
        const formattedBranches = uniqueBranches.map(branch => ({
          label: branch.charAt(0).toUpperCase() + branch.slice(1),
          value: branch
        }));
        setBranchOptions([allOption, ...formattedBranches]);
      }

      if (fileExtensionOptions.length === 0) {
        const uniqueFileExtensions =
          getUniqueOptions('fileExtension').filter(Boolean);
        setFileExtensionOptions([
          allOption,
          ...uniqueFileExtensions.map(fileExtension => ({
            label: fileExtension,
            value: fileExtension
          }))
        ]);
      }

      if (repositoryOptions.length === 0) {
        const uniqueRepositoryOptions =
          getUniqueOptions('repositoryType').filter(Boolean);
        setRepositoryOptions([
          allOption,
          ...uniqueRepositoryOptions.map(repositoryType => ({
            label: repositoryTypeMap[repositoryType as RepositoryType]?.label,
            value: repositoryType
          }))
        ]);
      }
    }
  }, [codeReferences]);

  useEffect(() => {
    if (apiKeys.length) {
      const validApiKey = apiKeys.find(item =>
        ['PUBLIC_API_WRITE', 'PUBLIC_API_ADMIN'].includes(item.role)
      );
      setHasValidApiKey(!!validApiKey);
    }
  }, [apiKeys]);
  if ((isLoadingApiKey || isLoadingCodeRefs) && initLoadedRef.current)
    return <PageLayout.LoadingState />;
  if (isEmpty)
    return (
      <EnableCodeRefs
        variant={!hasValidApiKey ? 'invalid' : 'no-data'}
        onAdd={() =>
          navigate(`/${currentEnvironment.urlCode}${PAGE_PATH_APIKEYS}`)
        }
      />
    );

  return (
    <PageLayout.Content className="p-6 pt-0 gap-y-6 min-w-[900px]">
      <FiltersBar
        repositoryOptions={repositoryOptions}
        branchOptions={branchOptions}
        extensionOptions={fileExtensionOptions}
        filters={filters}
        onChangeFilters={onChangeFilters}
      />
      {isLoadingCodeRefs ? (
        <FormLoading />
      ) : (
        <>
          {codeReferences?.map((codeRef, index) => (
            <CodeAccordion
              key={index}
              featureId={feature.id}
              codeRef={codeRef}
            />
          ))}
          <Pagination
            page={+filters.cursor}
            totalCount={totalCount}
            onChange={cursor => setFilters({ cursor: `${cursor}` })}
          />
        </>
      )}
    </PageLayout.Content>
  );
};

export default CodeReferencesPage;
