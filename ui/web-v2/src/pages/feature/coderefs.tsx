import React, { FC, memo, useEffect, useState } from 'react';
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
import { Highlight, themes, Prism } from 'prism-react-renderer';
import OpenInNewSvg from '../../assets/svg/open-new-tab.svg';

(typeof global !== 'undefined' ? global : window).Prism = Prism;
require('prismjs/components/prism-dart');

const repositoryTypeMap = {
  [CodeReference.RepositoryType.GITHUB]: {
    label: 'Github',
    icon: <GithubIcon className="w-5 h-5" />
  },
  [CodeReference.RepositoryType.GITLAB]: {
    label: 'Gitlab',
    icon: <GitlabIcon className="w-5 h-5" />
  },
  [CodeReference.RepositoryType.BITBUCKET]: {
    label: 'Bitbucket',
    icon: <BitbucketIcon className="w-5 h-5" />
  }
};

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

    const [repositoryOptions, setRepositoryOptions] = useState<Option[]>([]);
    const [branchOptions, setBranchOptions] = useState<Option[]>([]);
    const [fileExtensionOptions, setFileExtensionOptions] = useState<Option[]>(
      []
    );

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
      if (codeRefs.length > 0) {
        if (branchOptions.length === 0) {
          const uniqueBranches = [
            ...new Set(codeRefs.map((codeRef) => codeRef.repositoryBranch))
          ];
          const formattedBranches = uniqueBranches.map((branch) => ({
            label: branch.charAt(0).toUpperCase() + branch.slice(1),
            value: branch
          }));
          setBranchOptions([
            { label: 'All', value: null },
            ...formattedBranches
          ]);
        }

        if (fileExtensionOptions.length === 0) {
          const uniqueFileExtensions = [
            ...new Set(
              codeRefs.map((codeRef) => codeRef.fileExtension).filter(Boolean)
            )
          ];
          setFileExtensionOptions([
            { label: 'All', value: null },
            ...uniqueFileExtensions.map((fileExtension) => ({
              label: fileExtension,
              value: fileExtension
            }))
          ]);
        }

        if (repositoryOptions.length === 0) {
          const uniqueRepositoryOptions = [
            ...new Set(
              codeRefs.map((codeRef) => codeRef.repositoryType).filter(Boolean)
            )
          ];
          setRepositoryOptions([
            {
              label: 'All',
              value: null
            },
            ...uniqueRepositoryOptions.map((repositoryType) => ({
              label: repositoryTypeMap[repositoryType]?.label,
              value: repositoryType.toString()
            }))
          ]);
        }
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
              <h1 className="text-lg font-medium">
                {f(messages.codeRefs.enableCodeRefs)}
              </h1>
              <p className="text-sm text-gray-500">
                {f(messages.codeRefs.enableCodeRefsDescription)}
              </p>
              <DocumentationLink />
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
              {f(messages.apiKey.add.header.title)}
            </button>
          </div>
        </div>
      );
    }

    if (
      !selectedBranch &&
      !selectedFileExtension &&
      !selectedRepository &&
      codeRefs.length === 0
    ) {
      return (
        <div className="my-10 flex justify-center">
          <div className="w-[600px] text-gray-700 text-center space-y-1">
            <h1 className="text-lg font-medium">
              {f(messages.codeRefs.noRegisteredRefs)}
            </h1>
            <p className="text-sm">{f(messages.codeRefs.noRefsInCodebase)}</p>
            <DocumentationLink />
          </div>
        </div>
      );
    }

    return (
      <div className="p-10 bg-white">
        <div className="flex justify-between">
          <div className="h-full">
            <p className="font-semibold text-gray-900">
              {f(messages.feature.tab.codeRefs)}
            </p>
            <p className="text-sm text-gray-500">
              {f(messages.codeRefs.description)}
            </p>
          </div>
          <div className="flex space-x-4">
            <DocumentationLink />
            <Select
              placeholder={f(messages.all)}
              options={repositoryOptions}
              className={classNames('flex-none w-[200px]')}
              value={selectedRepository}
              onChange={setSelectedRepository}
              customControl={(props) => (
                <ControlComponent
                  {...props}
                  name={f(messages.codeRefs.repository)}
                />
              )}
            />
            <Select
              placeholder={f(messages.all)}
              options={branchOptions}
              className={classNames('flex-none w-[200px]')}
              value={selectedBranch}
              onChange={setSelectedBranch}
              customControl={(props) => (
                <ControlComponent
                  {...props}
                  name={f(messages.codeRefs.branch)}
                />
              )}
            />
            <Select
              placeholder={f(messages.all)}
              options={fileExtensionOptions}
              className={classNames('flex-none w-[210px]')}
              value={selectedFileExtension}
              onChange={setSelectedFileExtension}
              customControl={(props) => (
                <ControlComponent
                  {...props}
                  name={f(messages.codeRefs.fileExtensions)}
                />
              )}
            />
          </div>
        </div>
        <div className="mt-10">
          {isLoadingCodeRefs ? (
            <ListSkeleton />
          ) : (
            <div className="space-y-6">
              {codeRefs.map((codeRef) => {
                const occurrenceCount = countOccurrences(
                  codeRef.codeSnippet,
                  featureId
                );
                const branchLink = (
                  <a
                    href={codeRef.branchUrl}
                    className="text-primary underline"
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    {codeRef.repositoryBranch}
                  </a>
                );

                return (
                  <div
                    key={codeRef.id}
                    className="rounded bg-white shadow p-4 border border-gray-100"
                  >
                    <div className="flex py-1">
                      <div className="flex">
                        <div className="flex space-x-3 items-center">
                          {repositoryTypeMap[codeRef.repositoryType]?.icon}
                          <span>
                            {repositoryTypeMap[codeRef.repositoryType]?.label}
                          </span>
                        </div>
                        <div className="h-6 mx-4 border-l"></div>
                      </div>
                      <p className="text-gray-500">
                        {occurrenceCount > 1
                          ? f(messages.codeRefs.multipleReferenceFound, {
                              value: occurrenceCount,
                              branchLink
                            })
                          : f(messages.codeRefs.referenceFound, {
                              value: occurrenceCount,
                              branchLink
                            })}
                      </p>
                    </div>
                    <div className="mt-4">
                      <CodeAccordion codeRef={codeRef} featureId={featureId} />
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </div>
      </div>
    );
  }
);

const DocumentationLink = () => {
  const { formatMessage: f } = useIntl();
  return (
    <div className="text-sm flex text-center justify-center pt-1">
      <a
        href="https://docs.bucketeer.io/feature-flags/code-reference"
        target="_blank"
        rel="noreferrer"
        className="underline text-primary flex items-center space-x-1 ml-1"
      >
        <span>{f(messages.feature.documentation)}</span>
        <OpenInNewSvg />
      </a>
    </div>
  );
};

interface CodeAccordionProps {
  codeRef: CodeReference.AsObject;
  featureId: string;
}

const supportedExtensions = ['kt', 'swift', 'go', 'dart', 'js', 'ts'];

const CodeAccordion = ({ codeRef, featureId }: CodeAccordionProps) => {
  const { formatMessage: f } = useIntl();
  const [isOpen, setIsOpen] = useState(true);

  let language = codeRef.fileExtension.replace('.', '');

  if (!supportedExtensions.includes(language)) {
    language = 'js';
  }

  return (
    <div className="rounded-md bg-[#F8FAFC]">
      <button
        className="w-full flex justify-between items-center px-5 h-14 text-gray-700 cursor-pointer"
        onClick={() => setIsOpen(!isOpen)}
      >
        <span>{codeRef.filePath}</span>
        <div className="flex items-center space-x-4">
          <a
            href={codeRef.sourceUrl}
            className="text-primary underline px-2 py-1"
            target="_blank"
            onClick={(e) => e.stopPropagation()}
          >
            {f(messages.codeRefs.viewInSource)}
          </a>
          {isOpen ? (
            <ChevronUpIcon width={16} />
          ) : (
            <ChevronDownIcon width={16} />
          )}
        </div>
      </button>
      <div
        className={classNames(
          'overflow-x-auto transition-all duration-300 ease-in-out overflow-hidden overflow-y-auto rounded-b-md',
          isOpen ? 'max-h-60 border-t border-gray-300' : 'max-h-0'
        )}
      >
        <Highlight
          theme={themes.vsLight}
          code={codeRef.codeSnippet}
          language={language}
        >
          {({ style, tokens, getLineProps, getTokenProps }) => (
            <pre
              className="w-max min-w-full"
              style={{
                ...style,
                backgroundColor: '#F8FAFC'
              }}
            >
              {tokens.map((line, i) => {
                const lineProps = getLineProps({ line, key: i });

                return (
                  <div
                    {...lineProps}
                    style={{
                      backgroundColor: line.some((token) =>
                        token.content.includes(featureId)
                      )
                        ? '#E8E4F1'
                        : 'transparent'
                    }}
                  >
                    <span
                      className={classNames(
                        'inline-block w-16 text-right pr-4 select-none bg-primary text-white text-opacity-90',
                        i === 0 && 'pt-3',
                        i === tokens.length - 1 && 'pb-3'
                      )}
                    >
                      {codeRef.lineNumber + i}
                    </span>
                    {line.map((token, key) => {
                      const tokenProps = getTokenProps({ token, key });

                      const tokenContent = token.content;

                      // Split token content to highlight matched substring
                      const parts = tokenContent.split(
                        new RegExp(`(${featureId})`, 'gi')
                      );

                      return (
                        <span {...tokenProps} style={{ ...tokenProps.style }}>
                          {parts.map((part, index) =>
                            part.toLowerCase() === featureId.toLowerCase() ? (
                              <span
                                key={index}
                                style={{
                                  backgroundColor: '#D6CFE5'
                                }}
                              >
                                {part}
                              </span>
                            ) : (
                              <span key={index} className="">
                                {part}
                              </span>
                            )
                          )}
                        </span>
                      );
                    })}
                  </div>
                );
              })}
            </pre>
          )}
        </Highlight>
      </div>
    </div>
  );
};

const countOccurrences = (str: string, subStr: string) => {
  return (str.match(new RegExp(subStr, 'g')) || []).length;
};
