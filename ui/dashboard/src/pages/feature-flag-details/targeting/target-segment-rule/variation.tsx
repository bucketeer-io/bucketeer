import { useFormContext } from 'react-hook-form';
import { useTranslation } from 'i18n';
import { IconInfo } from '@icons';
import Form from 'components/form';
import Icon from 'components/icon';
import ServeDropdown from '../serve-dropdown';

const SegmentVariation = ({
  segmentIndex,
  ruleIndex
}: {
  segmentIndex: number;
  ruleIndex: number;
}) => {
  const { t } = useTranslation(['table', 'common']);

  const methods = useFormContext();
  const { control } = methods;

  return (
    <Form.Field
      control={control}
      name={`targetSegmentRules.${segmentIndex}.rules.${ruleIndex}.variation`}
      render={({ field }) => (
        <Form.Item className="flex flex-col flex-1 py-0 min-w-[170px]">
          <Form.Label required className="relative w-fit mb-5">
            {t('feature-flags.variation')}
            <Icon
              icon={IconInfo}
              size="xs"
              color="gray-500"
              className="absolute -right-6"
            />
          </Form.Label>
          <Form.Control>
            <ServeDropdown
              isExpand
              serveValue={field.value ? 1 : 0}
              onChangeServe={(value: number) => field.onChange(!!value)}
            />
          </Form.Control>
          <Form.Message />
        </Form.Item>
      )}
    />
  );
};

export default SegmentVariation;
