import { FilterTypes } from 'hooks/use-options';
import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import { IconPlus, IconTrash } from '@icons';
import { APIKeysFilters } from 'pages/api-keys/types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
import Dropdown, { DropdownOption, DropdownValue } from 'components/dropdown';
import Icon from 'components/icon';
import SlideModal from 'components/modal/slide';
import useFilterAPIKeyLogic from './use-filter-api-key-logic';

export type FilterProps = {
  onSubmit: (v: Partial<APIKeysFilters>) => void;
  isOpen: boolean;
  onClose: () => void;
  onClearFilters: () => void;
  filters?: Partial<APIKeysFilters>;
};

export interface Option {
  value: string;
  label: string;
}

const FilterAPIKeySlide = ({
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
    handleChangeOption,
    isLoadingEnvironments,
    onConfirmHandler,
    handleGetLabelFilterValue,
    handleChangeFilterValue,
    setSelectedFilters,
    isDisabledAddButton,
    isDisabledSubmitButton
  } = useFilterAPIKeyLogic(onSubmit, filters);
  return (
    <SlideModal title={t('filters')} isOpen={isOpen} onClose={onClose}>
      <div className="w-full h-full min-w-0 flex flex-col justify-between relative">
        <div className="flex flex-col min-w-0 w-full items-start p-5 gap-y-4">
          {selectedFilters.map((filterOption, filterIndex) => {
            const { label, value: filterType } = filterOption;
            const isEnvironmentFilter =
              filterType === FilterTypes.ENVIRONMENT_IDs;
            const valueOptions = getValueOptions(filterOption);

            return (
              <div
                className="flex min-w-0 items-start w-full h-[100px] gap-x-3"
                key={filterIndex}
              >
                <div className="h-full flex flex-col gap-y-4 items-center justify-center">
                  <div
                    className={cn(
                      'mt-2 typo-para-small text-center py-[3px] w-[42px] min-w-[42px] rounded text-accent-pink-500 bg-accent-pink-50',
                      {
                        'bg-gray-200 text-gray-600': filterIndex !== 0
                      }
                    )}
                  >
                    {t(filterIndex === 0 ? `if` : 'and')}
                  </div>
                  <Divider vertical={true} className="border-primary-500" />
                </div>
                <div className="flex flex-col w-full min-w-0">
                  <Dropdown
                    placeholder={t(`select-filter`)}
                    labelCustom={label}
                    className="w-full truncate py-2"
                    options={remainingFilterOptions.map(item => ({
                      ...item,
                      value: item.value || '',
                      label: item.label
                    }))}
                    value={filterType}
                    onChange={val =>
                      handleChangeOption(val as DropdownValue, filterIndex)
                    }
                    contentClassName="w-[270px]"
                  />
                  <div className="flex min-w-0 items-center gap-3 mt-3 pl-3">
                    <p className="typo-para-medium text-gray-600">{t('is')}</p>
                    <Dropdown
                      isSearchable={isEnvironmentFilter}
                      disabled={
                        (isEnvironmentFilter && isLoadingEnvironments) ||
                        !filterType
                      }
                      loading={isEnvironmentFilter && isLoadingEnvironments}
                      placeholder={t(`select-value`)}
                      labelCustom={handleGetLabelFilterValue(filterOption)}
                      className="w-full truncate py-2"
                      options={valueOptions as DropdownOption[]}
                      multiselect={isEnvironmentFilter}
                      value={
                        isEnvironmentFilter
                          ? (filterOption.filterValue as string[])
                          : (filterOption.filterValue as string)
                      }
                      onChange={val => {
                        handleChangeFilterValue(
                          val as string | number,
                          filterIndex
                        );
                      }}
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
                          selectedFilters.filter(
                            (_, index) => filterIndex !== index
                          )
                        )
                      }
                    >
                      <Icon icon={IconTrash} size={'sm'} />
                    </Button>
                  </div>
                </div>
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
          className="sticky bottom-0 left-0 bg-white"
          secondaryButton={
            <Button
              disabled={isDisabledSubmitButton}
              onClick={onConfirmHandler}
            >
              {t(`confirm`)}
            </Button>
          }
          primaryButton={
            <Button onClick={onClearFilters} variant="secondary">
              {t(`clear`)}
            </Button>
          }
        />
      </div>
    </SlideModal>
  );
};

export default FilterAPIKeySlide;
