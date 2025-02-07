import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { yupResolver } from '@hookform/resolvers/yup';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { cn } from 'utils/style';
import Divider from 'components/divider';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import Input from 'components/input';

type SituationType = 'compare' | 'user-segment' | 'date' | 'feature-flag';

interface Props {
  type: 'if' | 'and';
  situation: SituationType;
}

interface ConditionForm {
  situation: string;
  firstValue?: string;
  conditioner: string;
  secondValue?: string;
  value?: string;
  date?: string;
}

const formSchema = yup.object().shape({
  situation: yup.string().required(),
  conditioner: yup.string().required(),
  firstValue: yup.string(),
  secondValue: yup.string(),
  value: yup.string(),
  date: yup.string()
});

const Condition = ({ type, situation = 'compare' }: Props) => {
  const { t } = useTranslation(['form', 'common']);

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      situation,
      conditioner: situation === 'compare' ? '=' : '',
      firstValue: '1',
      secondValue: '',
      value: '',
      date: ''
    }
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

  const onSubmit: SubmitHandler<ConditionForm> = values => {
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
      <Divider vertical className="flex border-primary-500 w-[1px] h-[70px]" />
      <div className="flex items-center w-full flex-1 pl-4">
        <FormProvider {...form}>
          <Form
            onSubmit={form.handleSubmit(onSubmit)}
            className="flex items-end w-full gap-x-4"
          >
            <Form.Field
              control={form.control}
              name="situation"
              render={({ field }) => (
                <Form.Item className="flex flex-col flex-1 py-0 min-w-[170px]">
                  <Form.Label required>
                    {t('feature-flags.situation')}
                  </Form.Label>
                  <Form.Control>
                    <DropdownMenu>
                      <DropdownMenuTrigger
                        label={
                          situationOptions.find(
                            item => item.value === situation
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
              name="firstValue"
              render={({ field }) => (
                <Form.Item className="flex flex-col flex-1 py-0 min-w-[170px]">
                  <Form.Label required>
                    <Trans
                      i18nKey={'form:feature-flags.value-type'}
                      values={{
                        type: 'First'
                      }}
                    />
                  </Form.Label>
                  <Form.Control>
                    <Input {...field} />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <Form.Field
              control={form.control}
              name="conditioner"
              render={({ field }) => (
                <Form.Item className="flex flex-col flex-1 py-0 min-w-[170px]">
                  <Form.Label required>
                    {t('feature-flags.conditioner')}
                  </Form.Label>
                  <Form.Control>
                    <Input {...field} />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <Form.Field
              control={form.control}
              name="secondValue"
              render={({ field }) => (
                <Form.Item className="flex flex-col flex-1 py-0 min-w-[170px]">
                  <Form.Label required>
                    <Trans
                      i18nKey={'form:feature-flags.value-type'}
                      values={{
                        type: 'Second'
                      }}
                    />
                  </Form.Label>
                  <Form.Control>
                    <Input {...field} />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
          </Form>
        </FormProvider>
      </div>
    </div>
  );
};

export default Condition;
