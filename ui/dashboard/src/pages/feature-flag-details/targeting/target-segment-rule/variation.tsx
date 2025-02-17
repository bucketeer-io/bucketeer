import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { IconInfo } from '@icons';
import Form from 'components/form';
import Icon from 'components/icon';
import ServeDropdown from '../serve-dropdown';

const formSchema = yup.object().shape({
  variation: yup.number().required()
});

const SegmentVariation = ({ variation }: { variation: boolean }) => {
  const { t } = useTranslation(['table', 'common']);

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: { variation: !variation ? 0 : 1 }
  });

  const onSubmit: SubmitHandler<{ variation: number }> = values => {
    console.log(values);
  };

  return (
    <FormProvider {...form}>
      <Form
        onSubmit={form.handleSubmit(onSubmit)}
        className="flex items-end w-full gap-x-4"
      >
        <Form.Field
          control={form.control}
          name="variation"
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
                  serveValue={field.value}
                  onChangeServe={field.onChange}
                />
              </Form.Control>
              <Form.Message />
            </Form.Item>
          )}
        />
      </Form>
    </FormProvider>
  );
};

export default SegmentVariation;
