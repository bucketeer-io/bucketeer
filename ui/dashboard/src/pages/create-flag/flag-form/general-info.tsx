import { useMemo } from 'react';
import { useFormContext } from 'react-hook-form';
import { useTranslation } from 'i18n';
import { Tag } from '@types';
import { onGenerateSlug } from 'utils/converts';
import { IconInfo } from '@icons';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import { Tooltip } from 'components/tooltip';
import SelectMenu from 'elements/select-menu';

const GeneralInfo = ({
  tags,
  isUpdate
}: {
  tags: Tag[];
  isUpdate?: boolean;
}) => {
  const { t } = useTranslation(['form', 'common']);
  const tagOptions = useMemo(
    () =>
      tags?.map(tag => ({
        label: tag.name,
        value: tag.name
      })),
    [tags]
  );
  const { control, getFieldState, setValue, getValues } = useFormContext();

  return (
    <div className="flex flex-col w-full p-5 gap-y-6 bg-white rounded-lg shadow-card">
      <p className="typo-para-medium text-gray-700">{t('general-info')}</p>
      <div className="flex flex-col w-full gap-y-5">
        <div className="flex flex-col sm:flex-row w-full gap-4">
          <Form.Field
            control={control}
            name="name"
            render={({ field }) => (
              <Form.Item className="w-full py-0">
                <Form.Label required className="!mb-2">
                  {t('common:name')}
                </Form.Label>
                <Form.Control>
                  <Input
                    placeholder={`${t('placeholder-name')}`}
                    {...field}
                    name="flag-name"
                    onChange={value => {
                      field.onChange(value);
                      if (!isUpdate) {
                        const isFlagIdDirty = getFieldState('flagId').isDirty;
                        const flagId = getValues('flagId');
                        setValue(
                          'flagId',
                          isFlagIdDirty ? flagId : onGenerateSlug(value)
                        );
                      }
                    }}
                  />
                </Form.Control>
                <Form.Message />
              </Form.Item>
            )}
          />
          <Form.Field
            control={control}
            name="flagId"
            render={({ field }) => (
              <Form.Item className="w-full py-0">
                <Form.Label required className="relative w-fit !mb-2">
                  {t('feature-flags.flag-id')}
                  <Tooltip
                    content={t('flag-id-tooltip')}
                    trigger={
                      <div className="flex-center size-fit absolute top-0 -right-6">
                        <Icon icon={IconInfo} size="xs" color="gray-500" />
                      </div>
                    }
                    className="max-w-[400px]"
                  />
                </Form.Label>
                <Form.Control>
                  <Input
                    placeholder={`${t('feature-flags.placeholder-flag')}`}
                    disabled={isUpdate}
                    {...field}
                    name="flag-id"
                  />
                </Form.Control>
                <Form.Message />
              </Form.Item>
            )}
          />
        </div>
        <div className="flex flex-col sm:flex-row w-full gap-4">
          <Form.Field
            control={control}
            name="description"
            render={({ field }) => (
              <Form.Item className="w-full py-0">
                <Form.Label optional className="!mb-2">
                  {t('description')}
                </Form.Label>
                <Form.Control>
                  <Input
                    placeholder={t('placeholder-desc')}
                    {...field}
                    name="flag-description"
                  />
                </Form.Control>
                <Form.Message />
              </Form.Item>
            )}
          />
          <Form.Field
            control={control}
            name={`tags`}
            render={({ field }) => (
              <Form.Item className="w-full py-0 overflow-hidden">
                <Form.Label required className="relative w-fit !mb-2">
                  {t('common:tags')}
                  <Tooltip
                    content={t('tags-tooltip')}
                    trigger={
                      <div className="flex-center size-fit absolute top-0 -right-6">
                        <Icon icon={IconInfo} size="xs" color="gray-500" />
                      </div>
                    }
                    className="max-w-[400px]"
                  />
                </Form.Label>
                <Form.Control>
                  <SelectMenu
                    fieldValues={field.value || []}
                    onChange={field.onChange}
                    options={tagOptions}
                  />
                </Form.Control>
                <Form.Message />
              </Form.Item>
            )}
          />
        </div>
      </div>
    </div>
  );
};

export default GeneralInfo;
