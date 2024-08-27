import { Button } from 'components/button';
import Select from 'components/select';
import AuthWrapper from './elements/auth-wrapper';

const SelectOrganization = () => {
  return (
    <AuthWrapper>
      <h1 className="text-gray-900 typo-head-bold-huge">
        {`Select your Organization`}
      </h1>
      <p className="text-gray-600 typo-para-medium mt-4">
        {`Select the organization you want to work for.`}
      </p>
      <div className="mt-10">
        <Select
          label="Organization"
          options={[
            {
              value: '1',
              label: 'default'
            }
          ]}
          required
          placeholder="Select your Organization"
        />
      </div>
      <Button className="w-full mt-10">{`Continue`}</Button>
    </AuthWrapper>
  );
};

export default SelectOrganization;
