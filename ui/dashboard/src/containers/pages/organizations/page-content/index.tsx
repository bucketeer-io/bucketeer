import {
  KeyboardEvent,
  ReactNode,
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
import {
  organizationArchive,
  OrganizationsFetcherParams
} from '@api/organization';
import {
  organizationCreator,
  OrganizationCreatorCommand
} from '@api/organization/organization-creator';
import { organizationUpdate } from '@api/organization/organization-update';
import { useQueryOrganizations } from '@queries/organizations';
import { LIST_PAGE_SIZE } from 'constants/app';
import { FormFieldProps } from 'containers/common-form';
import CommonSlider from 'containers/common-slider';
import Filter from 'containers/filter';
import { SortingType } from 'containers/pages/projects';
import TableContent from 'containers/table-content';
import { commonTabs } from 'helpers/tab';
import { ColumnType, NotifyProps, useToast } from 'hooks';
import { useTranslation } from 'i18n';
import { debounce } from 'lodash';
import * as yup from 'yup';
import { OrderBy, Organization, OrganizationCollection } from '@types';
import { sortingFn } from 'utils/sort';
import { IconInfo } from '@icons';
import { Button } from 'components/button';
import Icon from 'components/icon';
import { PopoverValue } from 'components/popover';
import Tab from 'components/tab';

export type SubmitType = 'archived' | 'created' | 'updated';

type SubmitSuccessProps = {
  name: string;
  message?: ReactNode;
  submitType?: SubmitType;
  cb?: () => void;
  notify: (props: NotifyProps) => void;
};

export const onSubmitSuccess = ({
  name,
  submitType,
  message,
  cb,
  notify
}: SubmitSuccessProps) => {
  if (cb) cb();
  notify({
    toastType: 'toast',
    messageType: 'success',
    message: message || (
      <p>
        <strong>{name}</strong> has been successfully {`${submitType}!`}
      </p>
    )
  });
};

export const OrganizationsContent = () => {
  const navigate = useNavigate();
  const { notify } = useToast();
  const { t } = useTranslation(['common', 'form', 'table']);

  const initLoadedRef = useRef(true);

  const defaultParams: OrganizationsFetcherParams = useMemo(
    () => ({
      pageSize: LIST_PAGE_SIZE,
      cursor: String(0),
      orderBy: 'DEFAULT',
      orderDirection: 'ASC',
      searchKeyword: '',
      disabled: false,
      archived: false
    }),
    []
  );

  const [targetTab, setTargetTab] = useState(commonTabs[0].value);
  const [isOpenSlider, setIsOpenSlider] = useState(false);
  const [isLoadingCreator, setIsLoadingCreator] = useState(false);
  const [submitType, setSubmitType] = useState<SubmitType>('created');
  const [orgSelected, setOrgSelected] = useState<Organization>();
  const [sortingState, setSortingState] = useState<SortingType>({
    id: 'default',
    orderBy: 'DEFAULT',
    orderDirection: 'ASC'
  });
  const [organizationData, setOrganizationData] =
    useState<OrganizationCollection>();
  const [cursor, setCursor] = useState(0);
  const [searchValue, setSearchValue] = useState('');
  const [organizationParams, setOrganizationParams] =
    useState<OrganizationsFetcherParams>(defaultParams);

  const { data, isLoading, refetch } = useQueryOrganizations({
    params: organizationParams
  });

  const columns = useMemo<ColumnType<Organization>[]>(
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
        cellType: 'title',
        onClickCell: row => navigate(`/organization/${row?.id}`)
      },
      {
        accessorKey: 'projectCount',
        id: 'projectCount',
        header: `${t('projects')}`,
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
        accessorKey: 'userCount',
        id: 'userCount',
        header: `${t('users')}`,
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
            label: `${t('table:popover.edit-org')}`,
            icon: IconEditOutlined,
            value: 'edit-org'
          },
          {
            label: `${t('table:popover.archive-org')}`,
            icon: IconArchiveOutlined,
            value: 'archive-org'
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
      },
      {
        name: 'ownerEmail',
        label: `${t('form:owner')}`,
        placeholder: `${t('form:placeholder-email')}`,
        isRequired: true,
        isExpand: true,
        fieldType: 'input'
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

  const handleOnClickSubmit = useCallback(
    (formValues: OrganizationCreatorCommand) => {
      setIsLoadingCreator(true);

      return (
        submitType === 'created'
          ? organizationCreator({
              command: { ...formValues, isSystemAdmin: true, isTrial: true }
            })
          : organizationUpdate({
              id: orgSelected?.id as string,
              changeDescriptionCommand: {
                description: formValues.description
              },
              renameCommand: {
                name: formValues.name
              }
            })
      )
        .then(() => {
          onSubmitSuccess({
            name: formValues.name,
            submitType,
            notify,
            cb: () => {
              refetch();
              setIsOpenSlider(false);
              setOrgSelected(undefined);
            }
          });
        })
        .finally(() => {
          setIsLoadingCreator(false);
        });
    },
    [submitType]
  );

  const handleClickPopover = useCallback(
    (value: PopoverValue, row?: Organization) => {
      if (value === 'edit-org') {
        if (submitType !== 'updated') setSubmitType('updated');
        setIsOpenSlider(true);
        return setOrgSelected(row);
      }
      handleArchiveEnv(row);
    },
    [submitType]
  );

  const handleArchiveEnv = (row?: Organization) => {
    if (row) {
      organizationArchive({
        id: row.id,
        command: {}
      }).then(() => {
        onSubmitSuccess({
          name: row.name,
          notify,
          submitType: 'archived',
          cb: () => {
            refetch();
          }
        });
      });
    }
  };

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
    if (data) setOrganizationData(data);
  }, [data]);

  useEffect(() => {
    if ((cursor >= 0 || sortingState) && !initLoadedRef.current) {
      setOrganizationParams(prev => ({
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
              if (submitType === 'updated') setSubmitType('created');
            }}
          >
            <Icon icon={IconAddOutlined} size="sm" />
            {t(`new-org`)}
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
            organizationData?.Organizations?.length
              ? organizationData.Organizations
              : data?.Organizations
          }
          paginationProps={{
            cursor,
            pageSize: LIST_PAGE_SIZE,
            totalCount: organizationData?.totalCount
              ? Number(organizationData?.totalCount)
              : 0,
            setCursor,
            cb: () => (initLoadedRef.current = false)
          }}
          emptyTitle={t(`table:empty.org-title`)}
          emptyDescription={t(`table:empty.org-desc`)}
          emptyActions={
            <div className="flex justify-center">
              <Button className="w-fit">
                <Icon icon={IconAddOutlined} size="sm" />
                {t(`new-org`)}
              </Button>
            </div>
          }
          sortingState={sortingState}
          onSortingTable={onSortingTable}
        />
        <CommonSlider
          title={t(submitType === 'created' ? `new-org` : `update-org`)}
          isLoading={isLoadingCreator}
          isOpen={isOpenSlider}
          formFields={formFields}
          submitTextBtn={t(
            submitType === 'created' ? `create-org` : `update-org`
          )}
          formSchema={formSchema}
          formData={orgSelected}
          onSubmit={formValues =>
            handleOnClickSubmit(formValues as OrganizationCreatorCommand)
          }
          onClose={() => {
            setIsOpenSlider(false);
            setOrgSelected(undefined);
          }}
        />
      </div>
    </div>
  );
};
