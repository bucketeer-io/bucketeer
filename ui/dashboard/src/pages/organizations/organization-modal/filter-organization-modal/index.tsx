import { useState } from 'react';
import { IconAddOutlined } from 'react-icons-material-design';
import { useTranslation } from 'i18n';
import { IconTrash } from '@icons';
import { OrganizationFilters } from 'pages/organizations/types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';

export type FilterProps = {
  onSubmit: (v: Partial<OrganizationFilters>) => void;
  isOpen: boolean;
  onClose: () => void;
  onClearFilters: () => void;
};

export interface Option {
  value: string;
  label: string;
}

export enum FilterTypes {
  ENABLED = 'enabled'
}

export const filterOptions: Option[] = [
  {
    value: FilterTypes.ENABLED,
    label: 'Enabled'
  }
];

export const enabledOptions: Option[] = [
  {
    value: 'yes',
    label: 'Yes'
  },
  {
    value: 'no',
    label: 'No'
  }
];

const FilterOrganizationModal = ({
  onSubmit,
  isOpen,
  onClose,
  onClearFilters
}: FilterProps) => {
  const { t } = useTranslation(['common']);
  const [selectedFilterType, setSelectedFilterType] = useState<Option>();
  const [valueOption, setValueOption] = useState<Option>();

  const onConfirmHandler = () => {
    switch (selectedFilterType?.value) {
      case FilterTypes.ENABLED:
        onSubmit({
          disabled: valueOption?.value === 'no'
        });
        return;
    }
  };

  const onClearHandler = () => {
    onClearFilters();
    setSelectedFilterType(undefined);
    setValueOption(undefined);
  };

  return (
    <DialogModal
      className="w-[665px]"
      title={t('filters')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="flex flex-col w-full items-start p-5 gap-y-4">
        <div className="flex items-center w-full h-12 gap-x-4">
          <div className="typo-para-small text-center py-[3px] px-4 rounded text-accent-pink-500 bg-accent-pink-50">
            {`If`}
          </div>
          <Divider vertical={true} className="border-primary-500" />
          <DropdownMenu>
            <DropdownMenuTrigger
              placeholder={`Select type`}
              label={selectedFilterType?.label}
              variant="secondary"
              className="w-full"
            />
            <DropdownMenuContent className="w-[235px]" align="start">
              {filterOptions.map((item, index) => (
                <DropdownMenuItem
                  key={index}
                  value={item.value}
                  label={item.label}
                  onSelectOption={() => setSelectedFilterType(item)}
                />
              ))}
            </DropdownMenuContent>
          </DropdownMenu>
          <p className="typo-para-medium text-gray-600">{`is`}</p>
          <DropdownMenu>
            <DropdownMenuTrigger
              placeholder={`Select value`}
              label={valueOption?.label}
              variant="secondary"
              className="w-full"
            />
            <DropdownMenuContent className="w-[235px]" align="start">
              {enabledOptions.map((item, index) => (
                <DropdownMenuItem
                  key={index}
                  value={item.value}
                  label={item.label}
                  onSelectOption={() => setValueOption(item)}
                />
              ))}
            </DropdownMenuContent>
          </DropdownMenu>

          <Button variant={'text'} size={'icon'} className="p-0 size-5">
            <Icon icon={IconTrash} size={'fit'} />
          </Button>
        </div>

        <Button variant={'text'} size={'sm'} className="px-0 typo-para-medium">
          <Icon icon={IconAddOutlined} size="sm" />
          {t('add-filter')}
        </Button>
      </div>

      <ButtonBar
        secondaryButton={
          <Button onClick={onConfirmHandler}>{`Confirm`}</Button>
        }
        primaryButton={
          <Button
            onClick={onClearHandler}
            variant="secondary"
          >{`Clear`}</Button>
        }
      />
    </DialogModal>
  );
};

export default FilterOrganizationModal;
