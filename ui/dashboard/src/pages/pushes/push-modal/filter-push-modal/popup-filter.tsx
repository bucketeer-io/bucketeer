import { FilterTypes } from 'hooks/use-options';
import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import { IconPlus, IconTrash } from '@icons';
import { NotificationFilters } from 'pages/notifications/types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
import Dropdown, { DropdownOption, DropdownValue } from 'components/dropdown';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';
import useFilterPushLogic from './use-filter-push-logic';

export type FilterProps = {
  onSubmit: (v: Partial<NotificationFilters>) => void;
  isOpen: boolean;
  onClose: () => void;
  onClearFilters: () => void;
  filters?: Partial<NotificationFilters>;
};

export interface Option {
  value: string;
  label: string;
}

const FilterPushPopup = ({
  onSubmit,
  isOpen,
  onClose,
  onClearFilters,
  filters
}: FilterProps) => {
  const { t } = useTranslation(['common']);
  const {
    selectedFilters,
    getValueOptions,
    remainingFilterOptions,
    setSelectedFilters,
    handleChangeFilterValue,
    isDisabledAddButton,
    isDisabledSubmitButton,
    handleChangeOption,
    onConfirmHandler,
    handleGetLabelFilterValue
  } = useFilterPushLogic(onSubmit, filters);
  return (
    <DialogModal
      className="max-w-[550px] lg:max-w-[750px]"
      title={t('filters')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="flex flex-col w-full items-start p-5 gap-y-4">
        {selectedFilters.map((filterOption, filterIndex) => {
          const { label, value: filterType } = filterOption;
          const isEnvironmentFilter =
            filterType === FilterTypes.ENVIRONMENT_IDs;
          const valueOptions = getValueOptions(filterOption);
          return (
            <div
              className="flex items-center w-full h-12 gap-x-4"
              key={filterIndex}
            >
              <div
                className={cn(
                  'typo-para-small text-center py-[3px] w-[42px] min-w-[42px] rounded text-accent-pink-500 bg-accent-pink-50',
                  {
                    'bg-gray-200 text-gray-600': filterIndex !== 0
                  }
                )}
              >
                {t(filterIndex === 0 ? `if` : 'and')}
              </div>
              <Divider vertical={true} className="border-primary-500" />
              <Dropdown
                labelCustom={label}
                placeholder={t(`select-filter`)}
                options={remainingFilterOptions.map(item => ({
                  ...item,
                  label: item.label,
                  value: item.value || ''
                }))}
                value={filterType || ''}
                onChange={value =>
                  handleChangeOption(value as DropdownValue, filterIndex)
                }
                className="w-full truncate"
                contentClassName="w-[235px]"
              />

              <p className="typo-para-medium text-gray-600">is</p>
              <Dropdown
                isSearchable={isEnvironmentFilter}
                disabled={!filterType}
                placeholder={t(`select-value`)}
                multiselect={isEnvironmentFilter}
                labelCustom={handleGetLabelFilterValue(filterOption)}
                options={valueOptions as DropdownOption[]}
                value={
                  isEnvironmentFilter
                    ? (filterOption.filterValue as string[])
                    : (filterOption.filterValue as string)
                }
                onChange={val => {
                  handleChangeFilterValue(val as string | number, filterIndex);
                }}
                className="w-full truncate"
                wrapTriggerStyle="truncate"
                contentClassName={cn('w-[235px]', {
                  'pt-0 w-[300px]': isEnvironmentFilter,
                  'hidden-scroll': valueOptions?.length > 15
                })}
              />
              <Button
                variant={'grey'}
                className="px-0 w-fit"
                disabled={selectedFilters.length <= 1}
                onClick={() =>
                  setSelectedFilters(
                    selectedFilters.filter((_, index) => filterIndex !== index)
                  )
                }
              >
                <Icon icon={IconTrash} size={'sm'} />
              </Button>
            </div>
          );
        })}
        <Button
          disabled={isDisabledAddButton}
          variant={'text'}
          className="h-6 px-0"
          onClick={() => {
            setSelectedFilters([
              ...selectedFilters,
              {
                label: '',
                value: undefined,
                filterValue: ''
              }
            ]);
          }}
        >
          <Icon icon={IconPlus} />
          {t('add-filter')}
        </Button>
      </div>

      <ButtonBar
        secondaryButton={
          <Button disabled={isDisabledSubmitButton} onClick={onConfirmHandler}>
            {t(`confirm`)}
          </Button>
        }
        primaryButton={
          <Button onClick={onClearFilters} variant="secondary">
            {t(`clear`)}
          </Button>
        }
      />
    </DialogModal>
  );
};

export default FilterPushPopup;
