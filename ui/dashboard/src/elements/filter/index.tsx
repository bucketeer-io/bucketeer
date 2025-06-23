import { ReactNode } from 'react';
import {
  IconLaunchOutlined,
  IconFilterListOutlined
} from 'react-icons-material-design';
import { Link } from 'react-router-dom';
import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import Button from 'components/button';
import Icon from 'components/icon';
import SearchInput from 'components/search-input';

interface FilterProps {
  action?: ReactNode;
  searchValue?: string;
  filterCount?: number;
  isShowDocumentation?: boolean;
  className?: string;
  link?: string;
  onSearchChange?: (value: string) => void;
  onOpenFilter?: () => void;
}

const Filter = ({
  action,
  searchValue = '',
  filterCount,
  isShowDocumentation = true,
  className,
  link = '',
  onSearchChange,
  onOpenFilter
}: FilterProps) => {
  const { t } = useTranslation(['common', 'form']);
  return (
    <div
      className={cn(
        'flex w-full lg:items-center justify-between flex-col lg:flex-row px-6 gap-x-6',
        { '!flex-row !justify-end': !onSearchChange },
        className
      )}
    >
      {onSearchChange && (
        <div className="w-full lg:w-[365px]">
          <SearchInput
            placeholder={`${t('form:placeholder-search-input')}`}
            value={searchValue}
            onChange={onSearchChange}
          />
        </div>
      )}
      <div className="flex items-center gap-4 mt-3 lg:mt-0 flex-wrap">
        {isShowDocumentation && (
          <Link
            target="_blank"
            to={link}
            onClick={e => {
              if (!link) return e.preventDefault();
            }}
          >
            <Button variant="text" className="flex-1 lg:flex-none">
              <Icon icon={IconLaunchOutlined} size="sm" />
              {t('documentation')}
            </Button>
          </Link>
        )}
        {onOpenFilter && (
          <Button
            variant="secondary"
            onClick={onOpenFilter}
            className="text-gray-700 shadow-border-gray-400 flex-1 lg:flex-none"
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
