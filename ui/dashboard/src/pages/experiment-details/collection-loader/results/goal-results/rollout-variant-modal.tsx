import { useMemo } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { yupResolver } from '@hookform/resolvers/yup';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { BestVariation, FeatureRuleStrategy, FeatureVariation } from '@types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Form from 'components/form';
import DialogModal from 'components/modal/dialog';
import { RadioGroup, RadioGroupItem } from 'components/radio';
import TextArea from 'components/textarea';

export type RolloutVariantModalProps = {
  isOpen: boolean;
  variations: FeatureVariation[];
  isRequireComment: boolean;
  defaultStrategy?: FeatureRuleStrategy;
  bestVariations: BestVariation[];
  onClose: () => void;
  onSubmit: (values: RolloutVariant) => void;
};

export interface RolloutVariant {
  variation: string;
  comment?: string;
}

const RolloutVariantModal = ({
  isOpen,
  defaultStrategy,
  variations,
  isRequireComment,
  bestVariations,
  onClose,
  onSubmit
}: RolloutVariantModalProps) => {
  const { t } = useTranslation(['table', 'common']);
  const formSchema = yup.object().shape({
    variation: yup.string().required(),
    comment: isRequireComment ? yup.string().required() : yup.string()
  });

  const bestVariation = useMemo(
    () =>
      bestVariations?.reduce(
        (acc: BestVariation | undefined, curr: BestVariation) =>
          acc
            ? curr.probability > (acc?.probability || 0)
              ? curr
              : acc
            : curr,
        undefined
      ),
    [bestVariations]
  );

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      variation:
        bestVariation?.id ||
        defaultStrategy?.fixedStrategy?.variation ||
        variations[0]?.id ||
        '',
      comment: ''
    }
  });

  const {
    control,
    formState: { isValid, isSubmitting }
  } = form;

  return (
    <DialogModal
      className="w-[500px]"
      title={t('results.rollout-variant')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <FormProvider {...form}>
        <Form onSubmit={form.handleSubmit(onSubmit)}>
          <div className="flex flex-col w-full items-start p-5 gap-y-5">
            <div className="typo-para-medium text-gray-600 w-full">
              {t('results.choose-variant')}
            </div>
            <Form.Field
              name="variation"
              control={control}
              render={({ field }) => (
                <Form.Item className="py-0">
                  <Form.Control>
                    <div className="flex flex-col w-full gap-y-4">
                      <RadioGroup
                        defaultValue={field.value}
                        onValueChange={field.onChange}
                      >
                        {variations.map(({ id, name, value }) => {
                          const rolloutVariation = bestVariations?.find(
                            item => item.id === id
                          );
                          return (
                            <div key={id} className="flex items-center gap-x-2">
                              <RadioGroupItem value={id} id={id} />
                              <label
                                htmlFor={id}
                                className="flex-1 typo-para-medium text-gray-600"
                              >
                                <Trans
                                  i18nKey={'table:results.variant-percent'}
                                  values={{
                                    name: name || value,
                                    percent: `${rolloutVariation ? (rolloutVariation.probability * 100).toFixed(1) : 0}`
                                  }}
                                />
                              </label>
                            </div>
                          );
                        })}
                      </RadioGroup>
                    </div>
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <Form.Field
              control={form.control}
              name="comment"
              render={({ field }) => (
                <Form.Item className="py-0 w-full">
                  <Form.Label
                    required={isRequireComment}
                    optional={!isRequireComment}
                  >
                    {t('form:comment-for-update')}
                  </Form.Label>
                  <Form.Control>
                    <TextArea
                      placeholder={`${t('form:placeholder-comment')}`}
                      rows={3}
                      {...field}
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
          </div>
          <ButtonBar
            secondaryButton={
              <Button type="submit" disabled={!isValid} loading={isSubmitting}>
                {t(`common:submit`)}
              </Button>
            }
            primaryButton={
              <Button type="button" onClick={onClose} variant="secondary">
                {t(`common:cancel`)}
              </Button>
            }
          />
        </Form>
      </FormProvider>
    </DialogModal>
  );
};

export default RolloutVariantModal;
