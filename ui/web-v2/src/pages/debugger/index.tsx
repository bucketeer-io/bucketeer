import React, { FC, memo, useCallback, useState } from 'react';
import { useDispatch } from 'react-redux';

import { Header } from '../../components/Header';
import { AppDispatch } from '../../store';
import { DebuggerEvaluateForm } from '../../components/DebuggerEvaluateForm';
import { FormProvider, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { addFormSchema } from './formSchema';
import { evaluateFeatures } from '../../modules/debugger';
import { useCurrentEnvironment } from '../../modules/me';
import { DebuggerResult } from '../../components/DebuggerResult';
import { EvaluateFeaturesResponse } from '../../proto/feature/service_pb';
import { getFeature } from '../../modules/features';
import { Feature } from '../../proto/feature/feature_pb';
import { Reason } from '../../proto/feature/reason_pb';
import { useIntl } from 'react-intl';
import { messages } from '../../lang/messages';

export interface UserEvaluation {
  id: string;
  userId: string;
  featureId: string;
  variationId: string;
  variationName: string;
  featureDetails: Feature.AsObject;
  reason: {
    type: Reason.TypeMap[keyof Reason.TypeMap];
    ruleId: string;
  };
}

export const DebuggerIndexPage: FC = memo(() => {
  const dispatch = useDispatch<AppDispatch>();
  const currentEnvironment = useCurrentEnvironment();
  const { formatMessage: f } = useIntl();

  const [showResults, setShowResults] = useState(false);
  const [userEvaluations, setUserEvaluations] = useState([]);

  const method = useForm({
    resolver: yupResolver(addFormSchema),
    defaultValues: {
      flag: [],
      userId: '',
      userAttributes: []
    },
    mode: 'onChange'
  });
  const { handleSubmit, reset, getValues } = method;

  const handleAddSubmit = useCallback(
    (data) => {
      const { flag, userId, userAttributes } = data;

      const evaluateFeaturesPromises = flag.map((f) =>
        dispatch(
          evaluateFeatures({
            environmentId: currentEnvironment.id,
            flag: f,
            userId: userId,
            userAttributes: userAttributes.map((ua) => [ua.key, ua.value])
          })
        ).then((e) => e.payload as EvaluateFeaturesResponse.AsObject)
      );

      Promise.all(evaluateFeaturesPromises).then((results) => {
        const allEvaluations = results.flatMap(
          (payload) => payload.userEvaluations?.evaluationsList || []
        );

        const featureDetailsPromises = allEvaluations.map((evaluation) =>
          dispatch(
            getFeature({
              environmentId: currentEnvironment.id,
              id: evaluation.featureId
            })
          ).then((e) => {
            const featureDetails = e.payload as Feature.AsObject;
            return {
              ...evaluation,
              featureDetails
            };
          })
        );

        Promise.all(featureDetailsPromises).then((detailedEvaluations) => {
          setUserEvaluations(detailedEvaluations);
          setShowResults(true);
        });
      });
    },
    [dispatch]
  );

  return (
    <>
      <div className="w-full">
        <Header
          title={f(messages.debugger.title)}
          description={f(messages.debugger.description)}
        />
      </div>
      <div className="m-10">
        <div className="bg-white border border-gray-300 rounded-md p-5">
          <FormProvider {...method}>
            <div className={showResults ? 'hidden' : 'block'}>
              <DebuggerEvaluateForm onSubmit={handleSubmit(handleAddSubmit)} />
            </div>
          </FormProvider>
          <div className={showResults ? 'block' : 'hidden'}>
            <DebuggerResult
              userId={getValues('userId')}
              userEvaluations={userEvaluations}
              editFields={() => setShowResults(false)}
              clearAllFields={() => {
                setShowResults(false);
                reset();
              }}
            />
          </div>
        </div>
      </div>
    </>
  );
});
