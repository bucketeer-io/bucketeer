import { yupResolver } from '@hookform/resolvers/yup';
import React, {
  useCallback,
  FC,
  createElement as jsx,
  memo,
  useEffect,
  useState,
} from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { useDispatch } from 'react-redux';

import { AnalysisForm } from '../../components/AnalysisForm';
import { AnalysisTable } from '../../components/AnalysisTable';
import { Header } from '../../components/Header';
import { messages } from '../../lang/messages';
import { getGoalCount } from '../../modules/goalCounts';
import { useCurrentEnvironment } from '../../modules/me';
import { listUserMetadata } from '../../modules/userMetadata';
import { AppDispatch } from '../../store';
import { addDays } from '../../utils/date';

import { formSchema } from './formSchema';

export const AnalysisIndexPage: FC = memo(() => {
  const { formatMessage: f } = useIntl();
  const dispatch = useDispatch<AppDispatch>();
  const currentEnvironment = useCurrentEnvironment();
  const [featureId, setFeatureId] = useState<string | null>(null);
  const methods = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      startAt: addDays(new Date(), -7),
      endAt: new Date(),
      goalId: null,
      featureId: '',
      featureVersion: 0,
      reason: '',
      userMetadata: [],
    },
    mode: 'onChange',
  });
  const { handleSubmit } = methods;

  const handleExecute = useCallback(
    (data) => {
      setFeatureId(data.featureId);
      dispatch(
        getGoalCount({
          environmentNamespace: currentEnvironment.namespace,
          startAt: data.startAt,
          endAt: data.endAt,
          featureId: data.featureId,
          featureVersion: data.featureVersion,
          reason: data.reason,
          goalId: data.goalId,
          filters: null,
          segments: data.segments,
        })
      );
    },
    [dispatch]
  );

  useEffect(() => {
    dispatch(
      listUserMetadata({
        environmentNamespace: currentEnvironment.namespace,
      })
    );
  }, [dispatch]);

  return (
    <>
      <div className="w-full">
        <Header
          title={f(messages.analysis.header.title)}
          description={f(messages.analysis.header.description)}
        />
      </div>
      <div className="px-5">
        <FormProvider {...methods}>
          <AnalysisForm onSubmit={handleSubmit(handleExecute)} />
        </FormProvider>
        <AnalysisTable featureId={featureId} />
      </div>
    </>
  );
});
