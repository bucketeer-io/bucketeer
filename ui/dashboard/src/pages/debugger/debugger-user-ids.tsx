import { useFormContext } from 'react-hook-form';
import { useTranslation } from 'i18n';
import { CreatableSelect } from 'components/creatable-select';
import Form from 'components/form';
import { AddDebuggerFormType } from './form-schema';

const DebuggerUserIds = () => {
  const { t } = useTranslation(['form']);
  const { control } = useFormContext<AddDebuggerFormType>();

  return (
    <Form.Field
      name={'userIds'}
      control={control}
      render={({ field }) => (
        <Form.Item className="py-0">
          <Form.Label required>{t('user-id')}</Form.Label>
          <Form.Control>
            <div className="flex flex-col w-full gap-y-2">
              <CreatableSelect
                defaultValues={field.value.map((item: string) => ({
                  label: item,
                  value: item
                }))}
                placeholder={t(`enter-user-ids`)}
                options={[]}
                onChange={ids => field.onChange(ids.map(item => item.value))}
              />
              <p className="typo-para-small text-gray-600">
                {t('enter-to-add-multiple')}
              </p>
            </div>
          </Form.Control>
          <Form.Message />
        </Form.Item>
      )}
    />
  );
};

export default DebuggerUserIds;
