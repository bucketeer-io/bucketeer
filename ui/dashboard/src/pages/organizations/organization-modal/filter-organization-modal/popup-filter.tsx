import { useTranslation } from 'i18n';
import { OrganizationFilters } from 'pages/organizations/types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
import Dropdown, { DropdownOption } from 'components/dropdown';
import DialogModal from 'components/modal/dialog';
import useFilterOrganizationLogic from './use-filter-organization-logic';

export type FilterProps = {
  onSubmit: (v: Partial<OrganizationFilters>) => void;
  isOpen: boolean;
  onClose: () => void;
  onClearFilters: () => void;
  filters?: Partial<OrganizationFilters>;
};

const FilterOrganizationPopup = ({
  onSubmit,
  isOpen,
  onClose,
  onClearFilters,
  filters
}: FilterProps) => {
  const { t } = useTranslation(['common']);
  const {
    selectedFilterType,
    filterEnabledOptions,
    isDisabledSubmitBtn,
    enabledOptions,
    selectedValue,
    setSelectedFilterType,
    setSelectedValue,
    onConfirmHandler
  } = useFilterOrganizationLogic(onSubmit, filters);

  return (
    <DialogModal
      className="w-[750px]"
      title={t('filters')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="flex flex-col w-full items-start p-5 gap-y-4">
        <div className="flex items-center w-full h-12 gap-x-4">
          <div className="typo-para-small text-center py-[3px] px-4 rounded text-accent-pink-500 bg-accent-pink-50">
            {t(`if`)}
          </div>
          <Divider vertical={true} className="border-primary-500" />
          <Dropdown
            value={selectedFilterType?.value}
            options={filterEnabledOptions as DropdownOption[]}
            onChange={value => {
              const selected = filterEnabledOptions.find(
                item => item.value === value
              );
              setSelectedFilterType(selected);
              setSelectedValue(undefined);
            }}
            placeholder={t(`select-filter`)}
            className="w-full"
            disabled={false}
            menuContentSide="bottom"
            itemClassName="w-[235px]"
          />

          <p className="typo-para-medium text-gray-600">is</p>
          <Dropdown
            options={enabledOptions as DropdownOption[]}
            value={selectedValue?.value}
            placeholder={t(`select-value`)}
            disabled={!selectedFilterType}
            onChange={value => {
              const selected = enabledOptions.find(
                item => item.value === value
              );
              setSelectedValue(selected);
            }}
            className="w-full"
            contentClassName="w-[235px]"
          />
        </div>
      </div>

      <ButtonBar
        secondaryButton={
          <Button disabled={isDisabledSubmitBtn} onClick={onConfirmHandler}>
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

export default FilterOrganizationPopup;
