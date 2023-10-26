import { unwrapUndefinable } from 'option-t/lib/Undefinable/unwrap';
import { FC } from 'react';

import { intl } from '../../../lang';
import { messages } from '../../../lang/messages';
import { GoalResult } from '../../../proto/eventcounter/goal_result_pb';
import { Variation } from '../../../proto/feature/variation_pb';
import { classNames } from '../../../utils/css';
import {
  HeaderCell,
  Table,
  TableBody,
  TableCell,
  TableHeader,
  TableRow,
} from '../Table';

const createHeadCells = (): Array<HeaderCell> => [
  {
    id: 'variation',
    label: intl.formatMessage(messages.experiment.result.variation.label),
    helpText: '',
  },
  {
    id: 'conversion-rate',
    label: intl.formatMessage(messages.experiment.result.conversionRate.label),
    helpText: intl.formatMessage(
      messages.experiment.result.conversionRate.helpText
    ),
  },
  {
    id: 'improvement',
    label: intl.formatMessage(messages.experiment.result.improvement.label),
    helpText: intl.formatMessage(
      messages.experiment.result.improvement.helpText
    ),
  },
  {
    id: 'prob-beat-baseline',
    label: intl.formatMessage(
      messages.experiment.result.probabilityToBeatBaseline.label
    ),
    helpText: intl.formatMessage(
      messages.experiment.result.probabilityToBeatBaseline.helpText
    ),
  },
  {
    id: 'prob-best',
    label: intl.formatMessage(
      messages.experiment.result.probabilityToBest.label
    ),
    helpText: intl.formatMessage(
      messages.experiment.result.probabilityToBest.helpText
    ),
  },
];

interface ConversionRateTableProps {
  goalResult: GoalResult.AsObject;
  baseVariationId: string;
  variations: Map<string, Variation.AsObject>;
}

export const ConversionRateTable: FC<ConversionRateTableProps> = ({
  goalResult,
  baseVariationId,
  variations,
}) => {
  const baseVariationResult = unwrapUndefinable(
    goalResult.variationResultsList.find(
      (el) => el.variationId == baseVariationId
    )
  );
  const baseConversionRate =
    (unwrapUndefinable(baseVariationResult.experimentCount).userCount /
      unwrapUndefinable(baseVariationResult.evaluationCount).userCount) *
    100;

  return (
    <Table>
      <TableHeader cells={createHeadCells()} />
      <TableBody>
        {goalResult.variationResultsList.map((variationResult) => {
          const conversionRate =
            (unwrapUndefinable(variationResult.experimentCount).userCount /
              unwrapUndefinable(variationResult.evaluationCount).userCount) *
            100;
          const cvrProbBeeatBaseline = variationResult.cvrProbBeatBaseline;
          const cvrProbBest = variationResult.cvrProbBest;
          return (
            <TableRow key={variationResult.variationId}>
              <TableCell textLeft={true}>
                {' '}
                {
                  unwrapUndefinable(variations.get(variationResult.variationId))
                    .value
                }{' '}
              </TableCell>
              <TableCell textLeft={true}>
                {' '}
                {Number.isNaN(conversionRate)
                  ? 'n/a'
                  : conversionRate.toFixed(1) + ' %'}{' '}
              </TableCell>
              <TableCell textLeft={true}>
                {baseVariationId === variationResult.variationId
                  ? 'Baseline'
                  : Number.isNaN(conversionRate - baseConversionRate)
                  ? 'n/a'
                  : (conversionRate - baseConversionRate).toFixed(1) +
                    ' %'}{' '}
              </TableCell>
              <TableCell textLeft={true}>
                {' '}
                {baseVariationId === variationResult.variationId
                  ? 'Baseline'
                  : cvrProbBeeatBaseline
                  ? (
                      unwrapUndefinable(variationResult.cvrProbBeatBaseline)
                        .mean * 100
                    ).toFixed(1) + ' %'
                  : '-'}{' '}
              </TableCell>
              <TableCell textLeft={true}>
                {' '}
                {cvrProbBest
                  ? (
                      unwrapUndefinable(variationResult.cvrProbBest).mean * 100
                    ).toFixed(1)
                  : '-'}
                {' % '}
              </TableCell>
            </TableRow>
          );
        })}
      </TableBody>
    </Table>
  );
};
