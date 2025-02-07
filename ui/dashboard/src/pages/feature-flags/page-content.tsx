import { FunctionComponent, useCallback, useEffect, useState } from 'react';
import { IconAddOutlined } from 'react-icons-material-design';
import { usePartialState } from 'hooks';
import { useTranslation } from 'i18n';
import { pickBy } from 'lodash';
import { AnyObject } from 'yup';
import { isEmptyObject, isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import { cn } from 'utils/style';
import { IconGridView, IconListView } from '@icons';
import Button from 'components/button';
import Icon from 'components/icon';
import { Tabs, TabsContent, TabsList, TabsTrigger } from 'components/tabs';
import Filter from 'elements/filter';
import PageLayout from 'elements/page-layout';
import CollectionLoader from './collection-loader';
import { FlagActionType, FlagTabType, FlagsTemp, FlagsViewType } from './types';

const GridSwitchButton = ({
  icon,
  isActive,
  className,
  onClick
}: {
  icon: FunctionComponent;
  isActive: boolean;
  className: string;
  onClick: () => void;
}) => {
  return (
    <Button
      size={'icon'}
      variant={'secondary-2'}
      className={cn(
        'w-[52px] h-12 px-4 py-[14px] rounded-xl transition-colors',
        className
      )}
      onClick={onClick}
    >
      <Icon
        icon={icon}
        size={'sm'}
        color={isActive ? 'primary-500' : 'gray-500'}
      />
    </Button>
  );
};

const PageContent = ({
  onAdd,
  onHandleActions
}: {
  onAdd: () => void;
  onHandleActions: (item: FlagsTemp, type: FlagActionType) => void;
}) => {
  const { t } = useTranslation(['common']);
  const { searchOptions, onChangSearchParams } = useSearchParams();
  const searchFilters: Partial<AnyObject> = searchOptions;

  const defaultFilters = {
    page: 1,
    orderBy: 'CREATED_AT',
    orderDirection: 'DESC',
    status: 'ACTIVE',
    ...searchFilters
  } as AnyObject;

  const [filters, setFilters] = usePartialState<AnyObject>(defaultFilters);
  const [viewType, setViewType] = useState<FlagsViewType>('LIST_VIEW');

  const onChangeFilters = (values: Partial<AnyObject>) => {
    const options = pickBy({ ...filters, ...values }, v => isNotEmpty(v));
    onChangSearchParams(options);
    setFilters({ ...values });
  };

  useEffect(() => {
    if (isEmptyObject(searchOptions)) {
      setFilters({ ...defaultFilters });
    }
  }, [searchOptions]);

  const handleOnChangeViewType = useCallback((type: FlagsViewType) => {
    setViewType(type);
  }, []);

  return (
    <PageLayout.Content>
      <Filter
        action={
          <div className="flex items-center gap-x-4">
            <div className="flex items-center">
              <GridSwitchButton
                icon={IconListView}
                className="border-r-0 rounded-r-none"
                isActive={viewType === 'LIST_VIEW'}
                onClick={() => handleOnChangeViewType('LIST_VIEW')}
              />
              <GridSwitchButton
                icon={IconGridView}
                className="border-l-0 rounded-l-none"
                isActive={viewType === 'GRID_VIEW'}
                onClick={() => handleOnChangeViewType('GRID_VIEW')}
              />
            </div>
            <Button className="flex-1 lg:flex-none" onClick={onAdd}>
              <Icon icon={IconAddOutlined} size="sm" />
              {t(`new-flag`)}
            </Button>
          </div>
        }
        searchValue={filters.searchQuery}
        onSearchChange={searchQuery => onChangeFilters({ searchQuery })}
      />
      <Tabs
        className="flex-1 flex h-full flex-col mt-6"
        value={filters.status}
        onValueChange={value => {
          const status = value as FlagTabType;
          onChangeFilters({
            searchQuery: '',
            status
          });
        }}
      >
        <TabsList>
          <TabsTrigger value="ACTIVE">{t(`active`)}</TabsTrigger>
          <TabsTrigger value="FAVORITES">{t(`favorites`)}</TabsTrigger>
          <TabsTrigger value="ARCHIVED">{t(`archived`)}</TabsTrigger>
        </TabsList>

        <TabsContent value={filters.status}>
          <CollectionLoader
            viewType={viewType}
            onHandleActions={onHandleActions}
          />
        </TabsContent>
      </Tabs>
    </PageLayout.Content>
  );
};

export default PageContent;
