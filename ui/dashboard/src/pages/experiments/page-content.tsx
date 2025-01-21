import { useEffect } from 'react';
import { IconAddOutlined } from 'react-icons-material-design';
import { usePartialState, useToggleOpen } from 'hooks';
import { useTranslation } from 'i18n';
import { pickBy } from 'lodash';
import { isEmptyObject, isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import Button from 'components/button';
import Icon from 'components/icon';
import Filter from 'elements/filter';
import PageLayout from 'elements/page-layout';
import { ExperimentFilters } from './types';

const PageContent = ({
  onAdd,
  onHandleActions
}: {
  onAdd: () => void;
  onHandleActions: () => void;
}) => {
  const { t } = useTranslation(['common']);

  const { searchOptions, onChangSearchParams } = useSearchParams();
  const searchFilters: Partial<ExperimentFilters> = searchOptions;
  const defaultFilters = {
    page: 1,
    orderBy: 'CREATED_AT',
    orderDirection: 'DESC',
    ...searchFilters
  } as ExperimentFilters;

  const [filters, setFilters] =
    usePartialState<ExperimentFilters>(defaultFilters);

  const [openFilterModal, onOpenFilterModal, onCloseFilterModal] =
    useToggleOpen(false);

  const onChangeFilters = (values: Partial<ExperimentFilters>) => {
    const options = pickBy({ ...filters, ...values }, v => isNotEmpty(v));
    onChangSearchParams(options);
    setFilters({ ...values });
  };

  useEffect(() => {
    if (isEmptyObject(searchOptions)) {
      setFilters({ ...defaultFilters });
    }
  }, [searchOptions]);

  return (
    <PageLayout.Content>
      <Filter
        onOpenFilter={onOpenFilterModal}
        action={
          <Button className="flex-1 lg:flex-none" onClick={onAdd}>
            <Icon icon={IconAddOutlined} size="sm" />
            {t(`new-push`)}
          </Button>
        }
        searchValue={filters.searchQuery}
        filterCount={isNotEmpty(filters.disabled) ? 1 : undefined}
        onSearchChange={searchQuery => onChangeFilters({ searchQuery })}
      />
    </PageLayout.Content>
  );
};

export default PageContent;
