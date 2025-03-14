import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { yupResolver } from '@hookform/resolvers/yup';
import { useQueryTags } from '@queries/tags';
// import { useQueryClient } from '@tanstack/react-query';
import { LIST_PAGE_SIZE } from 'constants/app';
// import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import { cloneDeep } from 'lodash';
import { v4 as uuid } from 'uuid';
import {
  IconFlagJSON,
  IconFlagNumber,
  IconFlagString,
  IconFlagSwitch,
  IconInfo
} from '@icons';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import { FlagDataType } from 'pages/feature-flags/types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import SlideModal from 'components/modal/slide';
import TextArea from 'components/textarea';
import { formSchema } from './formSchema';
import Variations, { VariationType } from './variations';

interface AddFlagModalProps {
  isOpen: boolean;
  onClose: () => void;
}

type ServeType = {
  id?: string;
  value?: string;
};

export interface AddFlagForm {
  name: string;
  flagId: string;
  description?: string;
  tags: string[];
  flagType: FlagDataType;
  variations?: VariationType[];
  serveOn: ServeType;
  serveOff: ServeType;
}

const defaultVariations: VariationType[] = [
  {
    id: uuid(),
    value: 'True',
    name: '',
    description: ''
  },
  {
    id: uuid(),
    value: 'False',
    name: '',
    description: ''
  }
];

export const flagTypeOptions = [
  {
    label: 'Boolean',
    value: 'boolean',
    icon: IconFlagSwitch
  },
  {
    label: 'String',
    value: 'string',
    icon: IconFlagString
  },
  {
    label: 'Number',
    value: 'number',
    icon: IconFlagNumber
  },
  {
    label: 'JSON',
    value: 'json',
    icon: IconFlagJSON
  }
];

const AddFlagModal = ({ isOpen, onClose }: AddFlagModalProps) => {
  //   const queryClient = useQueryClient();
  const { t } = useTranslation(['common', 'form']);
  //   const { notify } = useToast();

  const { data: collection, isLoading: isLoadingEnvs } = useQueryTags({
    params: {
      pageSize: LIST_PAGE_SIZE,
      cursor: String(0)
    }
  });

  const tags = collection?.tags || [];

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      name: '',
      flagId: '',
      description: '',
      flagType: 'boolean',
      tags: [],
      variations: defaultVariations as VariationType[],
      serveOn: {
        id: defaultVariations[0].id,
        value: defaultVariations[0].value
      },
      serveOff: {
        id: defaultVariations[1].id,
        value: defaultVariations[1].value
      }
    }
  });
  const { watch } = form;

  const currentFlagOption = flagTypeOptions.find(
    item => item.value === watch('flagType')
  );

  const currentVariations = watch('variations') as VariationType[];

  const onSubmit: SubmitHandler<AddFlagForm> = values => {
    console.log(values);
  };

  return (
    <SlideModal title={t('new-flag')} isOpen={isOpen} onClose={onClose}>
      <div className="w-full p-5 pb-28">
        <p className="text-gray-700 typo-head-bold-small mb-2">
          {t('form:general-info')}
        </p>
        <FormProvider {...form}>
          <Form onSubmit={form.handleSubmit(onSubmit)}>
            <Form.Field
              control={form.control}
              name="name"
              render={({ field }) => (
                <Form.Item className="py-2.5">
                  <Form.Label required className="!mb-2">
                    {t('name')}
                  </Form.Label>
                  <Form.Control>
                    <Input
                      placeholder={`${t('form:placeholder-name')}`}
                      {...field}
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <Form.Field
              control={form.control}
              name="flagId"
              render={({ field }) => (
                <Form.Item className="py-2.5">
                  <Form.Label required className="relative w-fit !mb-2">
                    {t('form:feature-flags.flag-id')}
                    <Icon
                      icon={IconInfo}
                      size="xs"
                      color="gray-500"
                      className="absolute -right-6"
                    />
                  </Form.Label>
                  <Form.Control>
                    <Input
                      placeholder={`${t('form:feature-flags.placeholder-flag')}`}
                      {...field}
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <Form.Field
              control={form.control}
              name="description"
              render={({ field }) => (
                <Form.Item className="py-2.5">
                  <Form.Label optional className="!mb-2">
                    {t('form:description')}
                  </Form.Label>
                  <Form.Control>
                    <TextArea
                      placeholder={t('form:placeholder-desc')}
                      rows={4}
                      {...field}
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <Form.Field
              control={form.control}
              name={`tags`}
              render={({ field }) => (
                <Form.Item className="py-2.5">
                  <Form.Label required className="!mb-2">
                    {t('tags')}
                  </Form.Label>
                  <Form.Control>
                    <DropdownMenu>
                      <DropdownMenuTrigger
                        placeholder={t(`form:placeholder-tags`)}
                        label={
                          tags?.find(item => item.id === field.value[0])
                            ?.name || ''
                        }
                        variant="secondary"
                        className="w-full"
                      />
                      <DropdownMenuContent
                        className="w-[502px]"
                        align="start"
                        {...field}
                      >
                        {tags.map((item, index) => (
                          <DropdownMenuItem
                            {...field}
                            key={index}
                            value={item.id}
                            label={item.name}
                            onSelectOption={value => {
                              field.onChange([value]);
                            }}
                          />
                        ))}
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <Divider className="mt-2.5 mb-4" />
            <p className="text-gray-700 typo-head-bold-small mb-2">
              {t('form:feature-flags.flag-variations')}
            </p>
            <Form.Field
              control={form.control}
              name={`flagType`}
              render={({ field }) => (
                <Form.Item className="py-2">
                  <Form.Label required className="relative w-fit !mb-2">
                    {t('form:feature-flags.flag-type')}
                    <Icon
                      icon={IconInfo}
                      size="xs"
                      color="gray-500"
                      className="absolute -right-6"
                    />
                  </Form.Label>
                  <Form.Control>
                    <DropdownMenu>
                      <DropdownMenuTrigger
                        placeholder={t(`form:feature-flags.flag-type`)}
                        trigger={
                          <div className="flex items-center gap-x-2">
                            {currentFlagOption?.icon && (
                              <Icon
                                icon={currentFlagOption?.icon}
                                size={'md'}
                                className="flex-center"
                              />
                            )}
                            <p>{currentFlagOption?.label}</p>
                          </div>
                        }
                        disabled={isLoadingEnvs}
                        variant="secondary"
                        className="w-full"
                      />
                      <DropdownMenuContent
                        className="w-[502px]"
                        align="start"
                        {...field}
                      >
                        {flagTypeOptions.map((item, index) => (
                          <DropdownMenuItem
                            {...field}
                            key={index}
                            icon={item.icon}
                            value={item.value}
                            label={item.label}
                            onSelectOption={value => {
                              const cloneVariations =
                                cloneDeep(defaultVariations);
                              const newVariations =
                                value === 'boolean'
                                  ? cloneVariations
                                  : cloneVariations.map(item => ({
                                      ...item,
                                      value: ''
                                    }));
                              form.setValue('variations', newVariations);
                              field.onChange(value);
                            }}
                          />
                        ))}
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <Form.Field
              control={form.control}
              name="variations"
              render={({ field }) => (
                <Form.Item>
                  <Form.Control>
                    <Variations
                      flagType={watch('flagType')}
                      variations={currentVariations}
                      onChangeVariations={field.onChange}
                    />
                  </Form.Control>
                </Form.Item>
              )}
            />
            <div className="flex items-center w-full gap-x-4">
              <Form.Field
                control={form.control}
                name={`serveOn`}
                render={({ field }) => (
                  <Form.Item className="py-2.5 flex-1">
                    <Form.Label className="!mb-2">
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
                          trigger={(() => {
                            const variationIndex = currentVariations?.findIndex(
                              item => item.id === field.value.id
                            );
                            return (
                              <div className="flex items-center gap-x-2">
                                <FlagVariationPolygon
                                  color={
                                    variationIndex !== -1
                                      ? (variationIndex + 1) % 3 === 0
                                        ? 'green'
                                        : (variationIndex + 1) % 2 === 0
                                          ? 'pink'
                                          : 'blue'
                                      : 'blue'
                                  }
                                />
                                <Trans
                                  i18nKey={'form:feature-flags.variation'}
                                  values={{
                                    index:
                                      variationIndex !== -1
                                        ? variationIndex + 1
                                        : 1
                                  }}
                                />
                              </div>
                            );
                          })()}
                          variant="secondary"
                          className="w-full"
                        />
                        <DropdownMenuContent align="start" {...field}>
                          {currentVariations?.map((item, index) => (
                            <DropdownMenuItem
                              {...field}
                              key={index}
                              value={item.id}
                              label={`Variation ${index + 1}`}
                              onSelectOption={value => {
                                field.onChange({
                                  id: value,
                                  value: item.value
                                });
                              }}
                            />
                          ))}
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />
              <Form.Field
                control={form.control}
                name={`serveOff`}
                render={({ field }) => (
                  <Form.Item className="py-2.5 flex-1">
                    <Form.Label className="!mb-2">
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
                          trigger={(() => {
                            const variationIndex = currentVariations?.findIndex(
                              item => item.id === field.value.id
                            );
                            return (
                              <div className="flex items-center gap-x-2">
                                <FlagVariationPolygon
                                  color={
                                    variationIndex !== -1
                                      ? (variationIndex + 1) % 3 === 0
                                        ? 'green'
                                        : (variationIndex + 1) % 2 === 0
                                          ? 'pink'
                                          : 'blue'
                                      : 'pink'
                                  }
                                />
                                <Trans
                                  i18nKey={'form:feature-flags.variation'}
                                  values={{
                                    index:
                                      variationIndex !== -1
                                        ? variationIndex + 1
                                        : 2
                                  }}
                                />
                              </div>
                            );
                          })()}
                          variant="secondary"
                          className="w-full"
                        />
                        <DropdownMenuContent align="start" {...field}>
                          {currentVariations?.map((item, index) => (
                            <DropdownMenuItem
                              {...field}
                              key={index}
                              value={item.id}
                              label={`Variation ${index + 1}`}
                              onSelectOption={value => {
                                field.onChange({
                                  id: value,
                                  value: item.value
                                });
                              }}
                            />
                          ))}
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />
            </div>
            <div className="absolute left-0 bottom-0 bg-gray-50 w-full rounded-b-lg">
              <ButtonBar
                primaryButton={
                  <Button variant="secondary" onClick={onClose}>
                    {t(`cancel`)}
                  </Button>
                }
                secondaryButton={
                  <Button
                    type="submit"
                    disabled={!form.formState.isDirty}
                    loading={form.formState.isSubmitting}
                  >
                    {t(`submit`)}
                  </Button>
                }
              />
            </div>
          </Form>
        </FormProvider>
      </div>
    </SlideModal>
  );
};

export default AddFlagModal;
