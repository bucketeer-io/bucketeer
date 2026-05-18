import { useFormContext } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { getLanguage, Language, useTranslation } from 'i18n';
import { FeatureVariation } from '@types';
import { FlagFormSchema } from 'pages/create-flag/form-schema';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import Dropdown from 'components/dropdown';
import Form from 'components/form';

const makeVariationTrigger = (
  variations: FeatureVariation[],
  selectedId: string,
  fallbackLabel: (index: number) => string
) => {
  const idx = variations?.findIndex(v => v.id === selectedId) ?? -1;
  if (idx < 0) return undefined;
  const variation = variations[idx];

  return (
    <div className="flex items-center gap-x-2 w-0 flex-1 text-gray-700 typo-para-medium">
      <FlagVariationPolygon index={idx} />
      <p className="truncate">{variation.name || fallbackLabel(idx)}</p>
    </div>
  );
};

const DefaultVariations = () => {
  const { t } = useTranslation(['form', 'common', 'table']);
  const { control, watch } = useFormContext<FlagFormSchema>();
  const isJapaneseLanguage = getLanguage() === Language.JAPANESE;
  const currentVariations = watch('variations') as FeatureVariation[];

  const options =
    currentVariations?.map((item, index) => ({
      value: item.id,
      label: (
        <div className="flex items-center gap-x-2 text-gray-700 typo-para-medium">
          <FlagVariationPolygon index={index} />
          {item.name || t('feature-flags.variation', { index: index + 1 })}
        </div>
      )
    })) ?? [];

  const fallbackLabel = (index: number) =>
    t('feature-flags.variation', { index: index + 1 });

  return (
    <div className="flex items-center w-full gap-x-4 pb-6">
      <Form.Field
        control={control}
        name={`defaultOnVariation`}
        render={({ field }) => {
          return (
            <Form.Item className="py-0 flex-1">
              <Form.Label>
                <Trans
                  i18nKey={'form:feature-flags.serve-targeting'}
                  values={{
                    state: isJapaneseLanguage
                      ? t(`form:experiments.on`)
                      : t(`form:experiments.on`).toUpperCase()
                  }}
                />
              </Form.Label>
              <Form.Control>
                <Dropdown
                  options={options}
                  value={field.value}
                  onChange={value => field.onChange(value)}
                  placeholder={t('form:placeholder-tags')}
                  className="w-full"
                  trigger={makeVariationTrigger(
                    currentVariations,
                    field.value,
                    fallbackLabel
                  )}
                />
              </Form.Control>
              <Form.Message />
            </Form.Item>
          );
        }}
      />
      <Form.Field
        control={control}
        name={`defaultOffVariation`}
        render={({ field }) => {
          return (
            <Form.Item className="py-0 flex-1">
              <Form.Label>
                <Trans
                  i18nKey={'form:feature-flags.serve-targeting'}
                  values={{
                    state: isJapaneseLanguage
                      ? t(`form:experiments.off`)
                      : t(`form:experiments.off`).toUpperCase()
                  }}
                />
              </Form.Label>
              <Form.Control>
                <Dropdown
                  options={options}
                  value={field.value}
                  onChange={value => field.onChange(value)}
                  placeholder={t('form:placeholder-tags')}
                  className="w-full"
                  trigger={makeVariationTrigger(
                    currentVariations,
                    field.value,
                    fallbackLabel
                  )}
                />
              </Form.Control>
              <Form.Message />
            </Form.Item>
          );
        }}
      />
    </div>
  );
};

export default DefaultVariations;
