import {
  KeyboardEvent,
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState
} from 'react';
import {
  IconArchiveOutlined,
  IconEditOutlined
} from 'react-icons-material-design';
import { ProjectCreatorCommand, ProjectFetcherParams } from '@api/project';
import { projectUpdate } from '@api/project/project-update';
import { useQueryProjects } from '@queries/projects';
import { LIST_PAGE_SIZE } from 'constants/app';
import { FormFieldProps } from 'containers/common-form';
import CommonSlider from 'containers/common-slider';
import { onSubmitSuccess } from 'containers/pages/organizations';
import { SortingType } from 'containers/pages/projects';
import { useToast } from 'hooks';
import { ColumnType } from 'hooks/use-table';
import { useTranslation } from 'i18n';
import { debounce } from 'lodash';
import * as yup from 'yup';
import { OrderBy, Project, ProjectCollection } from '@types';
import { sortingFn } from 'utils/sort';
import { IconInfo } from '@icons';
import Icon from 'components/icon';
import { PopoverValue } from 'components/popover';
import Flag from 'components/table/table-row-items/flag';
import FilterLayout from '../filter-layout';
import { ContentDetailsProps } from '../page-content';

export const ProjectContent = ({ organizationId }: ContentDetailsProps) => {
  const defaultParams: ProjectFetcherParams = useMemo(
    () => ({
      pageSize: LIST_PAGE_SIZE,
      cursor: String(0),
      orderBy: 'DEFAULT',
      orderDirection: 'ASC',
      searchKeyword: '',
      disabled: false,
      organizationIds: organizationId ? [organizationId] : []
    }),
    [organizationId]
  );

  const initLoadedRef = useRef(true);
  const { notify } = useToast();
  const { t } = useTranslation(['common', 'form', 'table']);

  const [sortingState, setSortingState] = useState<SortingType>({
    id: 'default',
    orderBy: 'DEFAULT',
    orderDirection: 'ASC'
  });
  const [isOpenSlider, setIsOpenSlider] = useState(false);
  const [isLoadingSubmit, setIsLoadingSubmit] = useState(false);
  const [projectSelected, setProjectSelected] = useState<Project>();
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
        sortDescFirst: false,
        size: '25%',
        minSize: 160,
        sorting: true,
        headerCellType: 'title',
        renderFunc: row => (
          <Flag text={row?.name} status={row?.trial ? 'new' : undefined} />
        )
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
        header: `${t('table:flags')}`,
        size: '15%',
        minSize: 160,
        headerCellType: 'title',
        sorting: true
      },
      {
        accessorKey: 'createdAt',
        id: 'createdAt',
        header: `${t('table:created-at')}`,
        size: '15%',
        minSize: 160,
        headerCellType: 'title',
        sorting: true
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
    []
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
        isDisabled: true,
        labelIcon: (
          <div className="inline-flex items-center ml-2.5">
            <Icon icon={IconInfo} size={'fit'} />
          </div>
        )
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

  const handleClickPopover = useCallback(
    (value: PopoverValue, row?: Project) => {
      if (value === 'edit-prj') {
        setIsOpenSlider(true);
        return setProjectSelected(row);
      }
    },
    []
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

  const handleOnSubmit = useCallback((formValues: ProjectCreatorCommand) => {
    setIsLoadingSubmit(true);
    projectUpdate({
      id: formValues.id,
      changeDescriptionCommand: {
        description: formValues?.description
      },
      renameCommand: {
        name: formValues?.name
      }
    })
      .then(() =>
        onSubmitSuccess({
          name: formValues.name,
          submitType: 'updated',
          notify,
          cb: () => {
            refetch();
            setIsOpenSlider(false);
          }
        })
      )
      .finally(() => {
        setIsLoadingSubmit(false);
      });
  }, []);

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
    <>
      <FilterLayout
        isLoading={isLoading && initLoadedRef.current}
        columns={columns}
        data={
          projectData?.projects?.length ? projectData.projects : data?.projects
        }
        emptyTitle={t('table:empty.project-title')}
        emptyDescription={t('table:empty.project-org-desc')}
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
        searchValue={searchValue}
        onChangeSearchValue={handleChangeSearchValue}
        onKeyDown={handleKeyDown}
        onSortingTable={onSortingTable}
      />
      <CommonSlider
        title={t('update-project')}
        isLoading={isLoadingSubmit}
        isOpen={isOpenSlider}
        formFields={formFields}
        submitTextBtn={t('update-project')}
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
    </>
  );
};
