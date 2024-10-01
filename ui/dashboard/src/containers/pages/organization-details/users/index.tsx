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
import { AccountsFetcherParams } from '@api/account/accounts-fetcher';
import { EnvironmentsFetcherParams } from '@api/environment';
import { useQueryAccounts } from '@queries/accounts';
import { useQueryEnvironments } from '@queries/environments';
import { LIST_PAGE_SIZE } from 'constants/app';
import { SortingType } from 'containers/pages/projects';
import { ColumnType } from 'hooks/use-table';
import { useTranslation } from 'i18n';
import { debounce } from 'lodash';
import { Account, AccountCollection, OrderBy } from '@types';
import { sortingFn } from 'utils/sort';
import { PopoverValue } from 'components/popover';
import Text from 'components/table/table-row-items/text';
import FilterLayout from '../filter-layout';
import { ContentDetailsProps } from '../page-content';
import UserDetailsSlider from './_elements/user-details-slider';

export const User = ({ organizationId }: ContentDetailsProps) => {
  const { t } = useTranslation(['common', 'form', 'table']);
  const initLoadedRef = useRef(true);

  const defaultParams: AccountsFetcherParams = useMemo(
    () => ({
      pageSize: LIST_PAGE_SIZE,
      cursor: String(0),
      searchKeyword: '',
      orderBy: 'DEFAULT',
      orderDirection: 'ASC',
      disabled: false,
      organizationId
    }),
    [organizationId]
  );

  const [sortingState, setSortingState] = useState<SortingType>({
    id: 'default',
    orderBy: 'DEFAULT',
    orderDirection: 'ASC'
  });
  const [accountSelected, setAccountSelected] = useState<Account>();
  const [accountData, setAccountData] = useState<AccountCollection>();
  const [cursor, setCursor] = useState(0);
  const [searchValue, setSearchValue] = useState('');
  const [isOpenSlider, setIsOpenSlider] = useState(false);
  const [accountParams, setAccountParams] =
    useState<AccountsFetcherParams>(defaultParams);

  const environmentParams: EnvironmentsFetcherParams = useMemo(
    () => ({
      pageSize: LIST_PAGE_SIZE,
      cursor: String(cursor),
      orderBy: 'DEFAULT',
      orderDirection: 'ASC',
      searchKeyword: '',
      disabled: false,
      archived: false
    }),
    []
  );

  const { data, isLoading } = useQueryAccounts({
    params: accountParams
  });

  const { data: environmentData } = useQueryEnvironments({
    params: environmentParams,
    enabled: !!isOpenSlider
  });

  const columns = useMemo<ColumnType<Account>[]>(
    () => [
      {
        accessorKey: 'name',
        id: 'name',
        header: `${t('name')}`,
        sortDescFirst: false,
        size: '40%',
        minSize: 160,
        cellType: 'member',
        descriptionKey: 'email',
        sorting: true,
        sortingKey: 'EMAIL'
      },
      {
        accessorKey: 'organizationRole',
        id: 'organizationRole',
        header: `${t('role')}`,
        size: '15%',
        minSize: 160,
        cellType: 'text',
        sorting: true
      },
      {
        accessorKey: 'environmentCount',
        id: 'environmentCount',
        header: `${t('environments')}`,
        size: '15%',
        minSize: 160,
        headerCellType: 'title',
        renderFunc: row => (
          <Text text={String(row?.environmentRoles?.length || '')} />
        ),
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
        size: '10%',
        minSize: 52,
        headerCellType: 'empty',
        cellType: 'icon',
        options: [
          {
            label: `${t('table:popover.edit-user')}`,
            icon: IconEditOutlined,
            value: 'edit-user'
          },
          {
            label: `${t('table:popover.archive-user')}`,
            icon: IconArchiveOutlined,
            value: 'archive-user'
          }
        ],
        onClickPopover: (value, row) => handleClickPopover(value, row)
      }
    ],
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
    (value: PopoverValue, row?: Account) => {
      if (value === 'edit-user') {
        setIsOpenSlider(true);
        return setAccountSelected(row);
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

  useEffect(() => {
    if (data) {
      setAccountData(data);
    }
  }, [data]);

  useEffect(() => {
    if ((cursor >= 0 || sortingState) && !initLoadedRef.current) {
      setAccountParams(prev => ({
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
          accountData?.accounts?.length ? accountData.accounts : data?.accounts
        }
        emptyTitle={t('table:empty.user-title')}
        emptyDescription={t('table:empty.user-desc')}
        paginationProps={{
          cursor,
          pageSize: LIST_PAGE_SIZE,
          totalCount: accountData?.totalCount
            ? Number(accountData?.totalCount)
            : 0,
          setCursor,
          cb: () => (initLoadedRef.current = false)
        }}
        searchValue={searchValue}
        sortingState={sortingState}
        onChangeSearchValue={handleChangeSearchValue}
        onKeyDown={handleKeyDown}
        onSortingTable={onSortingTable}
      />
      <UserDetailsSlider
        isOpenSlider={isOpenSlider}
        accountSelected={accountSelected}
        environmentData={environmentData?.environments}
        setIsOpenSlider={setIsOpenSlider}
        setAccountSelected={setAccountSelected}
      />
    </>
  );
};
