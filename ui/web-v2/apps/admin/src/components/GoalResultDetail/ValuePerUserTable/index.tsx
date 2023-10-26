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
    id: 'value-per-user',
    label: intl.formatMessage(messages.experiment.result.valuePerUser.label),
    helpText: intl.formatMessage(
      messages.experiment.result.valuePerUser.helpText
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

interface ValuePerUserTableProps {
  goalResult: GoalResult.AsObject;
  baseVariationId: string;
  variations: Map<string, Variation.AsObject>;
}

export const ValuePerUserTable: FC<ValuePerUserTableProps> = ({
  goalResult,
  baseVariationId,
  variations,
}) => {
  const baseVariationResult = unwrapUndefinable(
    goalResult.variationResultsList.find(
      (el) => el.variationId == baseVariationId
    )
  );
  const baseValuePerUser =
    unwrapUndefinable(baseVariationResult.experimentCount).valueSum /
    unwrapUndefinable(baseVariationResult.experimentCount).userCount;

  return (
    <Table>
      <TableHeader cells={createHeadCells()} />
      <TableBody>
        {goalResult.variationResultsList.map((variationResult) => {
          const valuePerUser =
            unwrapUndefinable(variationResult.experimentCount).valueSum /
            unwrapUndefinable(variationResult.experimentCount).userCount;
          const probBeeatBaseline =
            variationResult.goalValueSumPerUserProbBeatBaseline;
          const probBest = variationResult.goalValueSumPerUserProbBest;
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
                {Number.isNaN(valuePerUser)
                  ? 'n/a'
                  : valuePerUser.toFixed(2)}{' '}
              </TableCell>
              <TableCell textLeft={true}>
                {baseVariationId === variationResult.variationId
                  ? 'Baseline'
                  : Number.isNaN(valuePerUser - baseValuePerUser)
                  ? 'n/a'
                  : (valuePerUser - baseValuePerUser).toFixed(1)}{' '}
              </TableCell>
              <TableCell textLeft={true}>
                {' '}
                {baseVariationId === variationResult.variationId
                  ? 'Baseline'
                  : probBeeatBaseline
                  ? (
                      unwrapUndefinable(
                        variationResult.goalValueSumPerUserProbBeatBaseline
                      ).mean * 100
                    ).toFixed(1) + ' %'
                  : '-'}{' '}
              </TableCell>
              <TableCell textLeft={true}>
                {' '}
                {probBest
                  ? (
                      unwrapUndefinable(
                        variationResult.goalValueSumPerUserProbBest
                      ).mean * 100
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
