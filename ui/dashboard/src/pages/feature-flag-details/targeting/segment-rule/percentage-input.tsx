import { useFormContext } from 'react-hook-form';
import { Feature } from '@types';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import Form from 'components/form';
import Input from 'components/input';
import InputGroup from 'components/input-group';
import { VariationOption } from './variation';

interface Props {
  feature: Feature;
  variationOptions: VariationOption[];
  variationId: string;
  name: string;
  handleChangeRolloutWeight: (value: number) => void;
}

const PercentageInput = ({
  feature,
  variationOptions,
  variationId,
  name,
  handleChangeRolloutWeight
}: Props) => {
  const { control } = useFormContext();
  const currentOption = variationOptions.find(
    item => item.value === variationId
  );

  return (
    <Form.Field
      control={control}
      name={name}
      render={({ field }) => {
        let value = String(field.value);
        value =
          value.startsWith('0') && value.length > 1
            ? value.toString().slice(1)
            : value;
        return (
          <Form.Item className="flex flex-col w-full gap-y-2 py-0">
            <Form.Control>
              <div className="flex items-center gap-x-2">
                <InputGroup
                  addon={'%'}
                  addonSlot="right"
                  className="w-[82px] overflow-hidden"
                  addonClassName="top-[1px] bottom-[1px] right-[1px] translate-x-0 translate-y-0 !flex-center rounded-r-lg bg-gray-200 w-[29px] typo-para-medium text-gray-700"
                >
                  <Input
                    {...field}
                    value={value}
                    onChange={value => {
                      field.onChange(+value);
                      handleChangeRolloutWeight(+value);
                    }}
                    onWheel={e => e.currentTarget.blur()}
                    type="number"
                    className="text-right pl-[5px]"
                  />
                </InputGroup>
                <div className="flex items-center gap-x-2 typo-para-small text-gray-600">
                  {feature.variationType === 'BOOLEAN' &&
                    currentOption?.variationValue && (
                      <FlagVariationPolygon
                        index={currentOption.variationValue === 'true' ? 0 : 1}
                      />
                    )}
                  {currentOption?.label}
                </div>
              </div>
            </Form.Control>
            <Form.Message />
          </Form.Item>
        );
      }}
    />
  );
};

export default PercentageInput;
