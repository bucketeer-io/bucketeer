import React, { memo, FC } from 'react';
import { Option, Select } from '../../components/Select';

import AuthWrapper from './authWrapper';

interface SelectOrganizationProps {
  options: Option[];
  onChange: (option: Option) => void;
  onSubmit: () => void;
  isSubmitBtnDisabled: boolean;
}

const SelectOrganization: FC<SelectOrganizationProps> = memo(
  ({ options, onChange, onSubmit, isSubmitBtnDisabled }) => {
    return (
      <AuthWrapper>
        <h2 className="font-bold text-xl mt-8">Organization</h2>
        <p className="mt-3 text-[#64738B]">
          Select the organization you want to access.
        </p>
        <div className="mt-8">
          <Select
            placeholder="Select organization"
            options={options}
            onChange={onChange}
          />
        </div>
        <button
          type="button"
          className="btn-submit btn mt-8 w-full"
          disabled={isSubmitBtnDisabled}
          onClick={onSubmit}
        >
          Submit
        </button>
      </AuthWrapper>
    );
  }
);

export default SelectOrganization;
