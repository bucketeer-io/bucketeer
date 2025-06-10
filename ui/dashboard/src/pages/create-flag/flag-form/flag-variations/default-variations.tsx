import { useFormContext } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { useTranslation } from 'i18n';
import { FeatureVariation } from '@types';
import { AddFlagForm } from 'pages/create-flag/form-schema';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';

const DefaultVariations = () => {
  const { t } = useTranslation(['form', 'common', 'table']);
  const { control, watch } = useFormContext<AddFlagForm>();

  const currentVariations = watch('variations') as FeatureVariation[];

  return (
    <div className="flex items-center w-full gap-x-4 pb-6">
      <Form.Field
        control={control}
        name={`defaultOnVariation`}
        render={({ field }) => {
          const variationIndex = currentVariations?.findIndex(
            item => item.id === field.value
          );

          return (
            <Form.Item className="py-0 flex-1">
              <Form.Label>
                <Trans
                  i18nKey={'form:feature-flags.serve-targeting'}
                  values={{
                    state: 'ON'
                  }}
                />
              </Form.Label>
              <Form.Control>
                <DropdownMenu>
                  <DropdownMenuTrigger
                    placeholder={t(`form:placeholder-tags`)}
                    trigger={
                      <div className="flex items-center gap-x-2 text-gray-700 typo-para-medium">
                        <FlagVariationPolygon index={variationIndex} />
                        <Trans
                          i18nKey={'form:feature-flags.variation'}
                          values={{
                            index:
                              currentVariations?.findIndex(
                                item => item.id === field.value
                              ) + 1
                          }}
                        />
                      </div>
                    }
                    variant="secondary"
                    className="w-full"
                  />
                  <DropdownMenuContent align="start" {...field}>
                    {currentVariations?.map((item, index) => (
                      <DropdownMenuItem
                        {...field}
                        key={index}
                        value={item.id}
                        label={
                          <div className="flex items-center gap-x-2 text-gray-700 typo-para-medium">
                            <FlagVariationPolygon index={index} />
                            <Trans
                              i18nKey={'form:feature-flags.variation'}
                              values={{
                                index: index + 1
                              }}
                            />
                          </div>
                        }
                        onSelectOption={() => {
                          field.onChange(item.id);
                        }}
                      />
                    ))}
                  </DropdownMenuContent>
                </DropdownMenu>
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
          const variationIndex = currentVariations?.findIndex(
            item => item.id === field.value
          );

          return (
            <Form.Item className="py-0 flex-1">
              <Form.Label>
                <Trans
                  i18nKey={'form:feature-flags.serve-targeting'}
                  values={{
                    state: 'OFF'
                  }}
                />
              </Form.Label>
              <Form.Control>
                <DropdownMenu>
                  <DropdownMenuTrigger
                    placeholder={t(`form:placeholder-tags`)}
                    trigger={
                      <div className="flex items-center gap-x-2 text-gray-700 typo-para-medium">
                        <FlagVariationPolygon index={variationIndex} />
                        <Trans
                          i18nKey={'form:feature-flags.variation'}
                          values={{
                            index:
                              currentVariations?.findIndex(
                                item => item.id === field.value
                              ) + 1
                          }}
                        />
                      </div>
                    }
                    variant="secondary"
                    className="w-full"
                  />
                  <DropdownMenuContent align="start" {...field}>
                    {currentVariations?.map((item, index) => (
                      <DropdownMenuItem
                        {...field}
                        key={index}
                        value={item.id}
                        label={
                          <div className="flex items-center gap-x-2 text-gray-700 typo-para-medium">
                            <FlagVariationPolygon index={index} />
                            <Trans
                              i18nKey={'form:feature-flags.variation'}
                              values={{
                                index: index + 1
                              }}
                            />
                          </div>
                        }
                        onSelectOption={() => {
                          field.onChange(item.id);
                        }}
                      />
                    ))}
                  </DropdownMenuContent>
                </DropdownMenu>
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
