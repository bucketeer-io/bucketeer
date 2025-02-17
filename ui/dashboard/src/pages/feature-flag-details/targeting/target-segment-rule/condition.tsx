import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { yupResolver } from '@hookform/resolvers/yup';
import { useQueryUserSegments } from '@queries/user-segments';
import { getCurrentEnvironment, useAuth } from 'auth';
import { LIST_PAGE_SIZE } from 'constants/app';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { cn } from 'utils/style';
import { IconTrash } from '@icons';
import Button from 'components/button';
import { ReactDatePicker } from 'components/date-time-picker';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import { SegmentConditionType } from '../types';

interface Props {
  type: 'if' | 'and';
  condition: SegmentConditionType;
  isDisabledDelete: boolean;
  onDeleteCondition: () => void;
  onChangeFormField: (field: string, value: string | number | boolean) => void;
}

const formSchema = yup.object().shape({
  situation: yup
    .string()
    .oneOf(['compare', 'user-segment', 'date', 'feature-flag'])
    .required(),
  conditioner: yup.string().required(),
  firstValue: yup.string(),
  secondValue: yup.string(),
  value: yup.string(),
  date: yup.string(),
  flag: yup.string()
});

const ConditionForm = ({
  type,
  condition,
  isDisabledDelete,
  onDeleteCondition,
  onChangeFormField
}: Props) => {
  const { t } = useTranslation(['form', 'common']);

  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const isCompare = condition.situation === 'compare';
  const isUserSegment = condition.situation === 'user-segment';
  const isDate = condition.situation === 'date';
  const isFlag = condition.situation === 'feature-flag';

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: condition
  });

  const situationOptions = [
    {
      label: t('feature-flags.compare'),
      value: 'compare'
    },
    {
      label: t('feature-flags.user-segment'),
      value: 'user-segment'
    },
    {
      label: t('feature-flags.date'),
      value: 'date'
    },
    {
      label: t('feature-flags.feature-flag'),
      value: 'feature-flag'
    }
  ];

  const conditionerDateOptions = [
    {
      label: 'Before',
      value: 'before'
    },
    {
      label: 'After',
      value: 'after'
    }
  ];

  const { data: segmentCollection, isLoading: segmentLoading } =
    useQueryUserSegments({
      params: {
        cursor: String(0),
        pageSize: LIST_PAGE_SIZE,
        environmentId: currentEnvironment.id
      },
      enabled: isUserSegment
    });

  const userSegments = segmentCollection?.segments || [];

  const onSubmit: SubmitHandler<SegmentConditionType> = values => {
    console.log(values);
  };

  return (
    <div className="flex items-center w-full gap-x-4">
      <div
        className={cn(
          'flex-center w-[42px] h-[26px] rounded-[3px] typo-para-small leading-[14px]',
          {
            'bg-accent-pink-50 text-accent-pink-500': type === 'if',
            'bg-gray-200 text-gray-600': type === 'and'
          }
        )}
      >
        {type === 'if' ? t('common:if') : t('common:and')}
      </div>
      <div className="flex items-center w-full flex-1 pl-4 border-l border-primary-500">
        <FormProvider {...form}>
          <Form
            onSubmit={form.handleSubmit(onSubmit)}
            className="flex items-end w-full gap-x-4"
          >
            <Form.Field
              control={form.control}
              name="situation"
              render={({ field }) => (
                <Form.Item className="flex flex-col flex-1 py-0 min-w-[170px] order-1">
                  <Form.Label required>
                    {t('feature-flags.situation')}
                  </Form.Label>
                  <Form.Control>
                    <DropdownMenu>
                      <DropdownMenuTrigger
                        label={
                          situationOptions.find(
                            item => item.value === condition.situation
                          )?.label
                        }
                        className="w-full"
                      />
                      <DropdownMenuContent align="start">
                        {situationOptions.map((item, index) => (
                          <DropdownMenuItem
                            key={index}
                            label={item.label}
                            value={item.value}
                            onSelectOption={value => {
                              field.onChange(value);
                              onChangeFormField('situation', value);

                              onChangeFormField(
                                'conditioner',
                                isCompare || isFlag
                                  ? '='
                                  : isDate
                                    ? 'before'
                                    : 'Is included in'
                              );
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
            {isCompare && (
              <Form.Field
                control={form.control}
                name="firstValue"
                render={({ field }) => (
                  <Form.Item className="flex flex-col flex-1 py-0 min-w-[170px] order-2">
                    <Form.Label required>
                      <Trans
                        i18nKey={'form:feature-flags.value-type'}
                        values={{
                          type: 'First'
                        }}
                      />
                    </Form.Label>
                    <Form.Control>
                      <Input
                        {...field}
                        value={condition.firstValue || field.value}
                        onChange={value => {
                          field.onChange(value);
                          onChangeFormField('firstValue', value);
                        }}
                      />
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />
            )}
            {(isUserSegment || isDate) && (
              <Form.Field
                control={form.control}
                name="value"
                render={({ field }) => (
                  <Form.Item
                    className={cn(
                      'flex flex-col flex-1 py-0 min-w-[170px] order-2',
                      { 'order-3': isUserSegment }
                    )}
                  >
                    <Form.Label required>
                      <Trans
                        i18nKey={'form:feature-flags.value'}
                        values={{
                          type: 'First'
                        }}
                      />
                    </Form.Label>
                    <Form.Control>
                      {isUserSegment ? (
                        <DropdownMenu>
                          <DropdownMenuTrigger
                            disabled={segmentLoading}
                            label={
                              userSegments?.find(
                                item => item.id === condition.value
                              )?.name || ''
                            }
                            placeholder={t('common:select-value')}
                            className="w-full"
                          />
                          <DropdownMenuContent align="start">
                            {userSegments?.map((item, index) => (
                              <DropdownMenuItem
                                key={index}
                                label={item.name}
                                value={item.id}
                                onSelectOption={value => {
                                  field.onChange(value);
                                  onChangeFormField('value', value);
                                }}
                              />
                            ))}
                          </DropdownMenuContent>
                        </DropdownMenu>
                      ) : (
                        <Input
                          {...field}
                          value={condition.value || field.value}
                          onChange={value => {
                            field.onChange(value);
                            onChangeFormField('value', value);
                          }}
                        />
                      )}
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />
            )}
            <Form.Field
              control={form.control}
              name="conditioner"
              render={({ field }) => (
                <Form.Item
                  className={cn(
                    'flex flex-col flex-1 py-0 min-w-[170px] order-3',
                    { 'order-2': condition.situation === 'user-segment' }
                  )}
                >
                  <Form.Label required>
                    {t('feature-flags.conditioner')}
                  </Form.Label>
                  <Form.Control>
                    {condition.situation === 'date' ? (
                      <DropdownMenu>
                        <DropdownMenuTrigger
                          label={
                            conditionerDateOptions.find(
                              item => item.value === condition.conditioner
                            )?.label
                          }
                          className="w-full"
                        />
                        <DropdownMenuContent align="start">
                          {conditionerDateOptions.map((item, index) => (
                            <DropdownMenuItem
                              key={index}
                              label={item.label}
                              value={item.value}
                              onSelectOption={value => {
                                field.onChange(value);
                                onChangeFormField('conditioner', value);
                              }}
                            />
                          ))}
                        </DropdownMenuContent>
                      </DropdownMenu>
                    ) : (
                      <Input
                        disabled={['user-segment', 'feature-flag'].includes(
                          condition.situation
                        )}
                        {...field}
                        value={condition.conditioner || field.value}
                        onChange={value => {
                          field.onChange(value);
                          onChangeFormField('conditioner', value);
                        }}
                      />
                    )}
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            {isCompare && (
              <Form.Field
                control={form.control}
                name="secondValue"
                render={({ field }) => (
                  <Form.Item className="flex flex-col flex-1 py-0 min-w-[170px] order-4">
                    <Form.Label required>
                      <Trans
                        i18nKey={'form:feature-flags.value-type'}
                        values={{
                          type: 'Second'
                        }}
                      />
                    </Form.Label>
                    <Form.Control>
                      <Input
                        {...field}
                        value={condition.secondValue || field.value}
                        onChange={value => {
                          field.onChange(value);
                          onChangeFormField('secondValue', value);
                        }}
                      />
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />
            )}
            {isDate && (
              <Form.Field
                control={form.control}
                name="date"
                render={({ field }) => (
                  <Form.Item className="flex flex-col flex-1 py-0 min-w-[170px] order-4">
                    <ReactDatePicker
                      selected={
                        condition?.date || field.value
                          ? new Date(Number(condition?.date || field.value))
                          : null
                      }
                      onChange={date => {
                        field.onChange(date?.getTime().toString());
                        onChangeFormField(
                          'date',
                          date?.getTime().toString() || ''
                        );
                      }}
                    />
                  </Form.Item>
                )}
              />
            )}
            <Button
              type="button"
              disabled={isDisabledDelete}
              variant={'grey'}
              className="flex items-center h-12 text-gray-500 hover:text-gray-600 order-5"
              onClick={onDeleteCondition}
            >
              <Icon icon={IconTrash} size={'sm'} />
            </Button>
          </Form>
        </FormProvider>
      </div>
    </div>
  );
};

export default ConditionForm;
