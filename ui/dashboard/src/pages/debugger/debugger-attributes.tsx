import { useFieldArray, useFormContext } from 'react-hook-form';
import { useTranslation } from 'i18n';
import { IconPlus, IconTrash } from '@icons';
import Button from 'components/button';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import { AddDebuggerFormType } from './form-schema';

const DebuggerAttributes = () => {
  const { control } = useFormContext<AddDebuggerFormType>();
  const { t } = useTranslation(['common', 'form']);

  const {
    fields: attributes,
    append,
    remove
  } = useFieldArray({
    name: 'attributes',
    control,
    keyName: 'debuggerAttribute'
  });

  return (
    <div className="flex flex-col w-full gap-y-5">
      <div className="flex flex-col w-full gap-y-3">
        <p className="typo-para-medium text-gray-700">
          {t('form:user-attributes')}
        </p>
        <p className="typo-para-small text-gray-600">
          {t('form:user-attributes-desc')}
        </p>
      </div>
      <div className="flex flex-col w-full gap-y-6">
        {attributes.map((_, index) => (
          <div key={index} className="flex items-end gap-x-4">
            <Form.Field
              name={`attributes.${index}.key`}
              control={control}
              render={({ field }) => (
                <Form.Item className="py-0 flex-1">
                  <Form.Label>{t('form:key')}</Form.Label>
                  <Form.Control>
                    <Input
                      {...field}
                      name="debugger-attribute-key"
                      onKeyDown={e => {
                        if (e.key === 'Enter') e.preventDefault();
                      }}
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <Form.Field
              name={`attributes.${index}.value`}
              control={control}
              render={({ field }) => (
                <Form.Item className="py-0 flex-1">
                  <Form.Label>{t('form:feature-flags.value')}</Form.Label>
                  <Form.Control>
                    <Input
                      {...field}
                      name="debugger-attribute-value"
                      onKeyDown={e => {
                        if (e.key === 'Enter') e.preventDefault();
                      }}
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <Button
              type="button"
              disabled={attributes.length <= 1}
              variant="grey"
              className="size-5 mb-4"
              onClick={() => remove(index)}
            >
              {<Icon icon={IconTrash} size="sm" />}
            </Button>
          </div>
        ))}
        <Button
          type="button"
          variant="text"
          className="w-fit px-0 h-6"
          onClick={() =>
            append({
              key: '',
              value: ''
            })
          }
        >
          <Icon icon={IconPlus} size="md" color="primary-500" />
          {t('add-attribute')}
        </Button>
      </div>
    </div>
  );
};

export default DebuggerAttributes;
