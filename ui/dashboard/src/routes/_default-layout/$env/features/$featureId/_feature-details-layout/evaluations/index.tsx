import { createRoute } from '@tanstack/react-router';
import { z } from 'zod';
import { EvaluationTimeRange } from '@types';
import EvaluationPage from 'pages/feature-flag-details/evaluation';
import { EvaluationTab } from 'pages/feature-flag-details/evaluation/types';
import { ErrorState } from 'elements/empty-state/error';
import { Route as FeatureDetailsLayoutRoute } from '../../_feature-details-layout';

const validateSearchSchema = z.object({
  period: z
    .enum(
      [
        EvaluationTimeRange.UNKNOWN,
        EvaluationTimeRange.TWENTY_FOUR_HOURS,
        EvaluationTimeRange.SEVEN_DAYS,
        EvaluationTimeRange.FOURTEEN_DAYS,
        EvaluationTimeRange.THIRTY_DAYS
      ],
      {
        message: 'Invalid period'
      }
    )
    .optional(),
  tab: z
    .enum([EvaluationTab.EVENT_COUNT, EvaluationTab.USER_COUNT], {
      message: 'Invalid Tab'
    })
    .optional()
});

export const Route = createRoute({
  path: 'evaluations',
  validateSearch: validateSearchSchema,
  getParentRoute: () => FeatureDetailsLayoutRoute,
  component: EvaluationPage,
  errorComponent: error => <ErrorState error={error} />
});

export type EvaluationSearch = z.infer<typeof validateSearchSchema>;
