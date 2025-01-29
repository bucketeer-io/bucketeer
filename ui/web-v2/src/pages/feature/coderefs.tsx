import React, { FC, memo, useEffect, useRef, useState } from 'react';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';
import { useHistory } from 'react-router-dom';
import { useCurrentEnvironment } from '../../modules/me';
import { AppDispatch } from '../../store';
import { listAPIKeys } from '../../modules/apiKeys';
import { getOrganizationId } from '../../storage/organizationId';
import { APIKEY_LIST_PAGE_SIZE } from '../../constants/apiKey';
import {
  listCodeRefs,
  selectAll as selectAllCodeRefs
} from '../../modules/codeRefs';
import { APIKey } from '../../proto/account/api_key_pb';
import { ListAPIKeysResponse } from '../../proto/account/service_pb';
import { ListCodeReferencesRequest } from '../../proto/coderef/service_pb';
import { AppState } from '../../modules';
import { PAGE_PATH_APIKEYS, PAGE_PATH_ROOT } from '../../constants/routing';
import { DetailSkeleton } from '../../components/DetailSkeleton';
import { CodeReference } from '../../proto/coderef/code_reference_pb';
import { classNames } from '../../utils/css';
import { Option, Select } from '../../components/Select';
import { components } from 'react-select';
import GithubIcon from '../../assets/svg/github-icon.svg';
import GitlabIcon from '../../assets/svg/gitlab-icon.svg';
import BitbucketIcon from '../../assets/svg/bitbucket-icon.svg';
import { ChevronDownIcon, ChevronUpIcon } from '@heroicons/react/outline';
import { useIntl } from 'react-intl';
import { messages } from '../../lang/messages';
import { ListSkeleton } from '../../components/ListSkeleton';
import SyntaxHighlighter from 'react-syntax-highlighter';
import { docco } from 'react-syntax-highlighter/dist/esm/styles/hljs';

const repositoryOptions = [
  {
    label: 'All',
    value: CodeReference.RepositoryType.REPOSITORY_TYPE_UNSPECIFIED.toString()
  },
  {
    label: 'GitHub',
    value: CodeReference.RepositoryType.GITHUB.toString()
  },
  {
    label: 'GitLab',
    value: CodeReference.RepositoryType.GITLAB.toString()
  },
  {
    label: 'Bitbucket',
    value: CodeReference.RepositoryType.BITBUCKET.toString()
  }
];

const fileExtensionOptions = [
  { label: 'All', value: null },
  { label: '.js', value: '.js' },
  { label: '.jsx', value: '.jsx' },
  { label: '.ts', value: '.ts' }
];

interface FeatureCodeRefsPageProps {
  featureId: string;
}

export const FeatureCodeRefsPage: FC<FeatureCodeRefsPageProps> = memo(
  ({ featureId }) => {
    const dispatch = useDispatch<AppDispatch>();
    const currentEnvironment = useCurrentEnvironment();
    const history = useHistory();
    const { formatMessage: f } = useIntl();

    const [isLoading, setIsLoading] = React.useState<boolean>(true);
    const [hasValidApiKey, setHasValidApiKey] = React.useState<boolean>(false);
    const [selectedRepository, setSelectedRepository] = useState<Option | null>(
      null
    );
    const [selectedBranch, setSelectedBranch] = useState<Option | null>(null);
    const [selectedFileExtension, setSelectedFileExtension] =
      useState<Option | null>(null);

    const [branchOptions, setBranchOptions] = useState<Option[]>([]);

    const isLoadingCodeRefs = useSelector<AppState, boolean>(
      (state) => state.codeRefs.loading,
      shallowEqual
    );
    const codeRefs = useSelector<AppState, CodeReference.AsObject[]>(
      (state) => selectAllCodeRefs(state.codeRefs),
      shallowEqual
    );

    const ControlComponent = ({ children, ...props }) => {
      return (
        <components.Control {...props}>
          <span className="ml-2">{props.name}:</span> {children}
        </components.Control>
      );
    };

    useEffect(() => {
      const fetchApiKeysAndCodeRefs = async () => {
        try {
          const res = await dispatch(
            listAPIKeys({
              organizationId: getOrganizationId(),
              environmentIds: [currentEnvironment.id],
              pageSize: APIKEY_LIST_PAGE_SIZE,
              cursor: '0',
              searchKeyword: '',
              orderBy: ListCodeReferencesRequest.OrderBy.DEFAULT,
              orderDirection: ListCodeReferencesRequest.OrderDirection.ASC
            })
          );
          const { apiKeysList } = res.payload as ListAPIKeysResponse.AsObject;

          const validApiKey = apiKeysList.some(
            (apiKey) =>
              apiKey.role === APIKey.Role.PUBLIC_API_ADMIN ||
              apiKey.role === APIKey.Role.PUBLIC_API_WRITE
          );

          if (validApiKey) {
            setHasValidApiKey(true);
          }
        } catch (error) {
          console.error('Error fetching API keys or code references:', error);
        } finally {
          setIsLoading(false);
        }
      };

      fetchApiKeysAndCodeRefs();
    }, []);

    useEffect(() => {
      const fetchFilteredCodeRefs = async () => {
        try {
          const repositoryType = selectedRepository
            ? Number(selectedRepository.value)
            : CodeReference.RepositoryType.REPOSITORY_TYPE_UNSPECIFIED;
          const repositoryBranch = selectedBranch ? selectedBranch.value : null;
          const fileExtension = selectedFileExtension
            ? selectedFileExtension.value
            : null;

          await fetchCodeRefs({
            repositoryType:
              repositoryType as CodeReference.RepositoryTypeMap[keyof CodeReference.RepositoryTypeMap],
            repositoryBranch,
            fileExtension
          });
        } catch (error) {
          console.error('Error fetching filtered code references:', error);
        }
      };

      if (hasValidApiKey) {
        fetchFilteredCodeRefs();
      }
    }, [
      selectedRepository,
      selectedBranch,
      selectedFileExtension,
      hasValidApiKey
    ]);

    useEffect(() => {
      if (branchOptions.length === 0 && codeRefs.length > 0) {
        setBranchOptions([
          { label: 'All', value: null },
          ...codeRefs.map((codeRef) => ({
            label:
              codeRef.repositoryBranch.charAt(0).toUpperCase() +
              codeRef.repositoryBranch.slice(1),
            value: codeRef.repositoryBranch
          }))
        ]);
      }
    }, [codeRefs]);

    const fetchCodeRefs = async ({
      fileExtension = null,
      repositoryBranch = null,
      repositoryType = null
    }: {
      fileExtension?: string;
      repositoryBranch?: string;
      repositoryType?: CodeReference.RepositoryTypeMap[keyof CodeReference.RepositoryTypeMap];
    } = {}) => {
      return await dispatch(
        listCodeRefs({
          environmentId: currentEnvironment.id,
          featureId: featureId,
          pageSize: 0,
          fileExtension,
          repositoryBranch,
          repositoryType
        })
      );
    };

    if (isLoading) {
      return (
        <div className="p-9 bg-gray-100">
          <DetailSkeleton />
        </div>
      );
    }

    if (!hasValidApiKey) {
      return (
        <div className="my-20 flex justify-center">
          <div className="w-[600px] text-gray-700 text-center">
            <div className="space-y-1">
              <h1 className="text-lg font-medium">Enable code references</h1>
              <p className="text-sm text-gray-500">
                with direct links from Bucketeer to the platform of your choice.
              </p>
              <p className="text-sm text-gray-500">
                Quickly see instances of feature flags being leveraged in your
                codebase,
              </p>
            </div>
            <button
              type="button"
              className="btn-submit mt-4"
              onClick={() => {
                history.push(
                  `${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_APIKEYS}`
                );
              }}
            >
              Create API Key
            </button>
          </div>
        </div>
      );
    }

    return (
      <div className="p-10">
        <div>
          <div className="flex justify-between">
            <div className="h-full">
              <p className="font-semibold text-gray-900">Code Refs</p>
              <p className="text-sm text-gray-500">
                References to this feature flag found in your codebase
              </p>
            </div>
            <div className="flex space-x-4">
              <Select
                placeholder={f(messages.all)}
                options={repositoryOptions}
                className={classNames('flex-none w-[200px]')}
                value={selectedRepository}
                onChange={setSelectedRepository}
                customControl={(props) => (
                  <ControlComponent {...props} name="Repository" />
                )}
              />
              <Select
                placeholder={f(messages.all)}
                options={branchOptions}
                className={classNames('flex-none w-[200px]')}
                value={selectedBranch}
                onChange={setSelectedBranch}
                customControl={(props) => (
                  <ControlComponent {...props} name="Branch" />
                )}
              />
              <Select
                placeholder={f(messages.all)}
                options={fileExtensionOptions}
                className={classNames('flex-none w-[200px]')}
                value={selectedFileExtension}
                onChange={setSelectedFileExtension}
                customControl={(props) => (
                  <ControlComponent {...props} name="File Extensions" />
                )}
              />
            </div>
          </div>
          <div className="mt-10">
            {isLoadingCodeRefs ? (
              <ListSkeleton />
            ) : codeRefs.length === 0 ? (
              <div className="my-10 flex justify-center">
                <div className="w-[600px] text-gray-700 text-center space-y-1">
                  <h1 className="text-lg font-medium">
                    No registered code references
                  </h1>
                  <p className="text-sm">
                    There are no code references in your codebase yet.
                  </p>
                  <p className="text-sm">
                    When a reference is added, it will appear here.
                  </p>
                </div>
              </div>
            ) : (
              <div className="space-y-6">
                {codeRefs.map((codeRef) => (
                  <div
                    key={codeRef.id}
                    className="rounded-md bg-white shadow p-4 border border-gray-200"
                  >
                    <div className="flex py-1">
                      <RepositoryType codeRef={codeRef} />

                      <p className="text-gray-500">
                        {countOccurrences(codeRef.codeSnippet, featureId)}{' '}
                        reference
                        {countOccurrences(codeRef.codeSnippet, featureId) > 1
                          ? '(s)'
                          : ''}{' '}
                        found in{' '}
                        <a
                          href={codeRef.branchUrl}
                          className="text-primary underline"
                          target="_blank"
                        >
                          {codeRef.repositoryBranch}
                        </a>{' '}
                        on the <a href="#">default</a> branch
                      </p>
                    </div>
                    <div className="mt-4">
                      <CodeAccordion codeRef={codeRef} featureId={featureId} />
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>
    );
  }
);

const RepositoryType = ({ codeRef }: { codeRef: CodeReference.AsObject }) => {
  let icon = null;
  let type = null;
  if (codeRef.repositoryType === CodeReference.RepositoryType.GITHUB) {
    icon = <GithubIcon className="w-5 h-5" />;
    type = 'Github';
  } else if (codeRef.repositoryType === CodeReference.RepositoryType.GITLAB) {
    icon = <GitlabIcon className="w-5 h-5" />;
    type = 'GitLab';
  } else if (
    codeRef.repositoryType === CodeReference.RepositoryType.BITBUCKET
  ) {
    icon = <BitbucketIcon className="w-5 h-5" />;
    type = 'Bitbucket';
  } else {
    return null;
  }

  return (
    <div className="flex">
      <div className="flex space-x-3 items-center">
        {icon}
        <span>{type}</span>
      </div>
      <div className="h-6 mx-4 border-l"></div>
    </div>
  );
};

interface CodeAccordionProps {
  codeRef: CodeReference.AsObject;
  featureId: string;
}

const CodeAccordion = ({ codeRef, featureId }: CodeAccordionProps) => {
  const [isOpen, setIsOpen] = useState(true);

  const ref = useRef(null);

  const lineHighlights = highlightLine({
    codeSnippet: codeRef.codeSnippet,
    text: featureId,
    lineNumber: codeRef.lineNumber
  });

  useEffect(() => {
    if (ref.current) {
      const html = ref.current.innerHTML;

      // Use a regex to replace the featureId with a styled span
      const highlightedHtml = html.replace(
        new RegExp(`\\b${featureId}\\b`, 'g'),
        `<span style="background-color: #5d3597; color: white; font-weight: bold;">${featureId}</span>`
      );

      ref.current.innerHTML = highlightedHtml;
    }
  }, [featureId]);

  return (
    <div className="rounded-md bg-[#F8FAFC]">
      <button
        className="w-full flex justify-between items-center px-5 py-4 text-gray-700 cursor-pointer"
        onClick={() => setIsOpen(!isOpen)}
      >
        <span>{codeRef.filePath}</span>
        <div className="flex items-center space-x-5">
          <a
            href={codeRef.sourceUrl}
            className="text-primary underline"
            target="_blank"
          >
            View in source
          </a>
          {isOpen ? (
            <ChevronUpIcon width={16} />
          ) : (
            <ChevronDownIcon width={16} />
          )}
        </div>
      </button>

      <div ref={ref}>
        <SyntaxHighlighter
          showLineNumbers
          startingLineNumber={codeRef.lineNumber}
          wrapLines
          lineProps={(lineNumber) => {
            const line = lineHighlights.find(
              (item) => item.lineNumber === lineNumber
            );
            const lineStyle = line ? line.style : {};
            return {
              style: {
                ...lineStyle,
                display: 'block',
                width: '100%'
              }
            };
          }}
          codeTagProps={{
            style: {
              display: 'block',
              width: '100%',
              padding: '0.5rem 0'
            }
          }}
          style={docco}
          customStyle={{
            backgroundColor: '#F8FAFC',
            padding: 0
          }}
          lineNumberStyle={{
            paddingLeft: '1.5rem'
          }}
        >
          {codeRef.codeSnippet}
        </SyntaxHighlighter>
      </div>
      {/* <div
        className={classNames(
          'transition-all duration-300 ease-in-out overflow-hidden overflow-y-scroll',
          isOpen ? 'max-h-60 border-t' : 'max-h-0'
        )}
      >
        <div
          style={{
            backgroundColor: '#F8FAFC',
            borderRadius: '8px',
            position: 'relative',
            overflowX: 'auto' // Enable horizontal scrolling
          }}
        >
          <pre
            style={{
              margin: 0,
              color: '#666666',
              fontFamily: 'monospace',
              whiteSpace: 'pre', // Prevent wrapping
              display: 'block',
              padding: '1rem 0',
              minWidth: '100%',
              width: 'max-content'
            }}
          >
            {codeRef.codeSnippet.split('\n').map((line, index) => (
              <div
                key={index}
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  paddingLeft: '1.25rem',
                  backgroundColor: line.includes(featureId)
                    ? '#E8E4F1'
                    : 'transparent'
                }}
              >
                <span
                  style={{
                    width: '2rem', // Ensure enough space for line numbers
                    textAlign: 'left',
                    flexShrink: 0 // Prevent shrinking of line numbers
                  }}
                >
                  {codeRef.lineNumber + index}
                </span>
                <code>{line}</code>
              </div>
            ))}
          </pre>
        </div>
      </div> */}
    </div>
  );
};

const highlightLine = ({
  codeSnippet,
  text,
  lineNumber
}: {
  codeSnippet: string;
  text: string;
  lineNumber: number;
}) => {
  const lines = codeSnippet.split('\n');
  return lines.map((line, index) => ({
    lineNumber: index + lineNumber,
    style: line.includes(text) ? { backgroundColor: '#E8E4F1' } : {} // Style the line containing the text
  }));
};

const countOccurrences = (str: string, subStr: string) => {
  return (str.match(new RegExp(subStr, 'g')) || []).length;
};
