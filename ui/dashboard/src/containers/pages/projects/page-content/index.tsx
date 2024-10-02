import {
  KeyboardEvent,
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState
} from 'react';
import {
  IconAddOutlined,
  IconArchiveOutlined,
  IconEditOutlined
} from 'react-icons-material-design';
import { useNavigate } from 'react-router-dom';
import { ProjectFetcherParams } from '@api/project';
import {
  projectCreator,
  ProjectCreatorCommand
} from '@api/project/project-creator';
import { projectUpdate } from '@api/project/project-update';
import { useQueryProjects } from '@queries/projects';
import { getCurrentEnvironment, useAuth } from 'auth';
import { LIST_PAGE_SIZE } from 'constants/app';
import { FormFieldProps } from 'containers/common-form';
import CommonSlider from 'containers/common-slider';
import Filter from 'containers/filter';
import { onSubmitSuccess, SubmitType } from 'containers/pages/organizations';
import TableContent from 'containers/table-content';
import { commonTabs } from 'helpers/tab';
import { ColumnType } from 'hooks/use-table';
import { useToast } from 'hooks/use-toast';
import { useTranslation } from 'i18n';
import { debounce } from 'lodash';
import * as yup from 'yup';
import { OrderBy, OrderDirection, Project, ProjectCollection } from '@types';
import { sortingFn } from 'utils/sort';
import { IconInfo } from '@icons';
import Button from 'components/button';
import Icon from 'components/icon';
import { PopoverValue } from 'components/popover';
import Tab from 'components/tab';
import Flag from 'components/table/table-row-items/flag';

export type SortingType = {
  id: string;
  orderBy: OrderBy;
  orderDirection: OrderDirection;
};

export const ProjectsContent = () => {
  const { consoleAccount, myOrganizations } = useAuth();
  const { notify } = useToast();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const { t } = useTranslation(['common', 'form', 'table']);
  const navigate = useNavigate();

  const initLoadedRef = useRef(true);

  const defaultParams: ProjectFetcherParams = useMemo(
    () => ({
      pageSize: LIST_PAGE_SIZE,
      cursor: String(0),
      orderBy: 'DEFAULT',
      orderDirection: 'ASC',
      searchKeyword: '',
      disabled: false,
      archived: false,
      organizationIds: !consoleAccount?.isSystemAdmin
        ? [myOrganizations[0].id]
        : []
    }),
    []
  );
  const [targetTab, setTargetTab] = useState(commonTabs[0].value);
  const [isOpenSlider, setIsOpenSlider] = useState(false);
  const [isLoadingSubmit, setIsLoadingSubmit] = useState(false);
  const [submitType, setSubmitType] = useState<SubmitType>('created');
  const [projectSelected, setProjectSelected] = useState<Project>();
  const [sortingState, setSortingState] = useState<SortingType>({
    id: 'default',
    orderBy: 'DEFAULT',
    orderDirection: 'ASC'
  });
  const [projectData, setProjectData] = useState<ProjectCollection>();
  const [cursor, setCursor] = useState(0);
  const [searchValue, setSearchValue] = useState('');
  const [projectParams, setProjectParams] =
    useState<ProjectFetcherParams>(defaultParams);

  const { data, isLoading, refetch } = useQueryProjects({
    params: projectParams
  });

  const columns = useMemo<ColumnType<Project>[]>(
    () => [
      {
        accessorKey: 'name',
        id: 'name',
        header: `${t('name')}`,
        size: '25%',
        minSize: 160,
        enableSorting: false,
        sorting: true,
        headerCellType: 'title',
        renderFunc: row => (
          <Flag text={row?.name} status={row?.trial ? 'new' : undefined} />
        ),
        onClickCell: row =>
          navigate(`/${currentEnvironment?.id}/project/${row?.id}`)
      },
      {
        accessorKey: 'creatorEmail',
        id: 'creatorEmail',
        header: `${t('maintainer')}`,
        size: '15%',
        minSize: 160,
        headerCellType: 'title',
        sorting: true
      },
      {
        accessorKey: 'environmentCount',
        id: 'environmentCount',
        header: `${t('environments')}`,
        size: '15%',
        minSize: 160,
        headerCellType: 'title',
        sorting: true
      },
      {
        accessorKey: 'featureFlagCount',
        id: 'featureFlagCount',
        header: t('table:flags'),
        size: '15%',
        minSize: 160,
        headerCellType: 'title',
        sorting: true
      },
      {
        accessorKey: 'createdAt',
        id: 'createdAt',
        header: t('table:created-at'),
        size: '15%',
        minSize: 160,
        sorting: true,
        headerCellType: 'title'
      },
      {
        accessorKey: 'actions',
        id: 'actions',
        header: 'Actions',
        size: '5%',
        minSize: 52,
        headerCellType: 'empty',
        cellType: 'icon',
        options: [
          {
            label: `${t('table:popover.edit-project')}`,
            icon: IconEditOutlined,
            value: 'edit-prj'
          },
          {
            label: `${t('table:popover.archive-project')}`,
            icon: IconArchiveOutlined,
            value: 'archive-prj'
          }
        ],
        onClickPopover: (value, row) => handleClickPopover(value, row)
      }
    ],
    [sortingState]
  );

  const formFields: FormFieldProps[] = useMemo(
    () => [
      {
        name: 'name',
        label: `${t('name')}`,
        placeholder: `${t('form:placeholder-name')}`,
        isRequired: true,
        isExpand: true,
        fieldType: 'input'
      },
      {
        name: 'urlCode',
        label: `${t('form:url-code')}`,
        placeholder: `${t('form:placeholder-code')}`,
        isRequired: true,
        isExpand: true,
        fieldType: 'input',
        labelIcon: (
          <div className="inline-flex items-center ml-2.5">
            <Icon icon={IconInfo} size={'fit'} />
          </div>
        ),
        isDisabled: submitType === 'updated' ? true : false
      },
      {
        name: 'description',
        label: `${t('form:description')}`,
        placeholder: `${t('form:placeholder-desc')}`,
        isOptional: true,
        isExpand: true,
        fieldType: 'textarea'
      }
    ],
    [submitType]
  );

  const formSchema = useMemo(
    () =>
      yup.object().shape({
        name: yup.string().required(),
        urlCode: yup.string().required()
      }),
    []
  );

  const handleChangeSearchValue = useCallback(
    debounce((value: string) => {
      initLoadedRef.current = false;
      setSearchValue(value);
    }),
    []
  );

  const handleKeyDown = useCallback((e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') handleChangeSearchValue.cancel();
  }, []);

  const handleOnSubmit = useCallback(
    (formValues: ProjectCreatorCommand) => {
      setIsLoadingSubmit(true);
      return (
        submitType === 'created'
          ? projectCreator({
              command: {
                ...formValues,
                id: currentEnvironment.organizationId
              }
            })
          : projectUpdate({
              id: formValues.id,
              changeDescriptionCommand: {
                description: formValues?.description
              },
              renameCommand: {
                name: formValues?.name
              }
            })
      )
        .then(() =>
          onSubmitSuccess({
            name: formValues.name,
            submitType,
            notify,
            cb: () => {
              refetch();
              setIsOpenSlider(false);
              setProjectSelected(undefined);
            }
          })
        )
        .finally(() => {
          setIsLoadingSubmit(false);
        });
    },
    [submitType]
  );

  const handleClickPopover = useCallback(
    (value: PopoverValue, row?: Project) => {
      if (value === 'edit-prj') {
        if (submitType !== 'updated') setSubmitType('updated');
        setIsOpenSlider(true);
        return setProjectSelected(row);
      }
    },
    [submitType]
  );

  const onSortingTable = useCallback(
    (accessorKey: string, sortingKey?: OrderBy) => {
      initLoadedRef.current = false;
      sortingFn({
        accessorKey,
        sortingKey,
        sortingState,
        setSortingState
      });
    },
    [data, sortingState]
  );

  useEffect(() => {
    if (data) {
      setProjectData(data);
    }
  }, [data]);

  useEffect(() => {
    if ((cursor >= 0 || sortingState) && !initLoadedRef.current) {
      setProjectParams(prev => ({
        ...prev,
        searchKeyword: searchValue.trim().toLowerCase(),
        orderBy: sortingState.orderBy,
        orderDirection: sortingState.orderDirection,
        cursor: String(cursor)
      }));
    }
  }, [cursor, sortingState, searchValue]);

  return (
    <div className="py-8 px-6">
      <Filter
        additionalActions={
          <Button
            className="flex-1 lg:flex-none"
            onClick={() => {
              setIsOpenSlider(true);
              if (submitType !== 'created') setSubmitType('created');
            }}
          >
            <Icon icon={IconAddOutlined} size="sm" />
            {t(`new-project`)}
          </Button>
        }
        searchValue={searchValue}
        onChangeSearchValue={handleChangeSearchValue}
        onKeyDown={handleKeyDown}
      />
      <div className="mt-6">
        <Tab
          options={commonTabs}
          value={targetTab}
          onSelect={value => setTargetTab(value)}
        />

        <TableContent
          isLoading={isLoading && initLoadedRef.current}
          columns={columns}
          data={
            projectData?.projects?.length
              ? projectData.projects
              : data?.projects
          }
          className="mt-5"
          emptyTitle={t('table:empty.project-title')}
          emptyDescription={t('table:empty.project-desc')}
          emptyActions={
            <div className="flex-center">
              <Button className="w-fit">
                <Icon icon={IconAddOutlined} size="sm" />
                {t(`new-project`)}
              </Button>
            </div>
          }
          paginationProps={{
            cursor,
            pageSize: LIST_PAGE_SIZE,
            totalCount: projectData?.totalCount
              ? Number(projectData?.totalCount)
              : 0,
            setCursor,
            cb: () => (initLoadedRef.current = false)
          }}
          sortingState={sortingState}
          onSortingTable={onSortingTable}
        />
      </div>
      <CommonSlider
        title={t(submitType === 'created' ? 'new-project' : 'update-project')}
        isLoading={isLoadingSubmit}
        isOpen={isOpenSlider}
        formFields={formFields}
        submitTextBtn={t(
          submitType === 'created' ? 'create-project' : 'update-project'
        )}
        formSchema={formSchema}
        formData={projectSelected}
        onSubmit={formValues =>
          handleOnSubmit(formValues as ProjectCreatorCommand)
        }
        onClose={() => {
          setIsOpenSlider(false);
          setProjectSelected(undefined);
        }}
      />
    </div>
  );
};
