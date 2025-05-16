import { ReactNode } from 'react';
import {
  IconLaunchOutlined,
  IconFilterListOutlined
} from 'react-icons-material-design';
import { useTranslation } from 'i18n';
import Button from 'components/button';
import Icon from 'components/icon';
import SearchInput from 'components/search-input';

interface FilterProps {
  action?: ReactNode;
  searchValue: string;
  filterCount?: number;
  isShowDocumentation?: boolean;
  onSearchChange: (value: string) => void;
  onOpenFilter?: () => void;
}

const Filter = ({
  action,
  searchValue,
  filterCount,
  isShowDocumentation = true,
  onSearchChange,
  onOpenFilter
}: FilterProps) => {
  const { t } = useTranslation(['common', 'form']);
  return (
    <div className="flex lg:items-center justify-between flex-col lg:flex-row px-6 gap-x-6">
      <div className="w-full lg:w-[365px]">
        <SearchInput
          placeholder={`${t('form:placeholder-search-input')}`}
          value={searchValue}
          onChange={onSearchChange}
        />
      </div>
      <div className="flex items-center gap-4 mt-3 lg:mt-0 flex-wrap">
        {isShowDocumentation && (
          <Button variant="text" className="flex-1 lg:flex-none">
            <Icon icon={IconLaunchOutlined} size="sm" />
            {t('documentation')}
          </Button>
        )}
        {onOpenFilter && (
          <Button
            variant="secondary"
            onClick={onOpenFilter}
            className="text-gray-600 shadow-border-gray-400 flex-1 lg:flex-none"
          >
            <Icon icon={IconFilterListOutlined} size="sm" />
            {t('filter')}
            {filterCount && (
              <div className="size-5 flex-center rounded-full bg-gray-200 text-[11px] text-gray-700">
                {filterCount}
              </div>
            )}
          </Button>
        )}
        {action}
      </div>
    </div>
  );
};

export default Filter;
