import { unwrapUndefinable } from 'option-t/lib/Undefinable/unwrap';
import { FC } from 'react';

import { messages } from '../../..//lang/messages';
import { intl } from '../../../lang';
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
    id: 'evaluation-user',
    label: intl.formatMessage(messages.experiment.result.evaluationUser.label),
    helpText: intl.formatMessage(
      messages.experiment.result.evaluationUser.helpText
    ),
  },
  {
    id: 'goals',
    label: intl.formatMessage(messages.experiment.result.goals.label),
    helpText: intl.formatMessage(messages.experiment.result.goals.helpText),
  },
  {
    id: 'goal-user',
    label: intl.formatMessage(messages.experiment.result.goalUser.label),
    helpText: intl.formatMessage(messages.experiment.result.goalUser.helpText),
  },
  {
    id: 'conversion rate',
    label: intl.formatMessage(messages.experiment.result.conversionRate.label),
    helpText: intl.formatMessage(
      messages.experiment.result.conversionRate.helpText
    ),
  },
  {
    id: 'value-sum',
    label: intl.formatMessage(messages.experiment.result.valueSum.label),
    helpText: intl.formatMessage(messages.experiment.result.valueSum.helpText),
  },
  {
    id: 'value-per-user',
    label: intl.formatMessage(messages.experiment.result.valuePerUser.label),
    helpText: intl.formatMessage(
      messages.experiment.result.valuePerUser.helpText
    ),
  },
];

interface GoalResultTableProps {
  goalResult: GoalResult.AsObject;
  variations: Map<string, Variation.AsObject>;
}

export const GoalResultTable: FC<GoalResultTableProps> = ({
  goalResult,
  variations,
}) => {
  return (
    <Table>
      <TableHeader cells={createHeadCells()} />
      <TableBody>
        {goalResult.variationResultsList.map((variationResult) => {
          const conversionRate =
            (unwrapUndefinable(variationResult.experimentCount).userCount /
              unwrapUndefinable(variationResult.evaluationCount).userCount) *
            100;
          const valuePerUser =
            unwrapUndefinable(variationResult.experimentCount).valueSum /
            unwrapUndefinable(variationResult.experimentCount).userCount;
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
                {unwrapUndefinable(
                  variationResult.evaluationCount
                ).userCount.toLocaleString()}{' '}
              </TableCell>
              <TableCell textLeft={true}>
                {' '}
                {unwrapUndefinable(
                  variationResult.experimentCount
                ).eventCount.toLocaleString()}{' '}
              </TableCell>
              <TableCell textLeft={true}>
                {' '}
                {unwrapUndefinable(
                  variationResult.experimentCount
                ).userCount.toLocaleString()}{' '}
              </TableCell>
              <TableCell textLeft={true}>
                {' '}
                {Number.isNaN(conversionRate)
                  ? 'n/a'
                  : conversionRate.toFixed(1) + ' %'}{' '}
              </TableCell>
              <TableCell textLeft={true}>
                {' '}
                {unwrapUndefinable(
                  variationResult.experimentCount
                ).valueSum.toLocaleString()}{' '}
              </TableCell>
              <TableCell textLeft={true}>
                {' '}
                {Number.isNaN(valuePerUser)
                  ? 'n/a'
                  : valuePerUser.toFixed(2)}{' '}
              </TableCell>
            </TableRow>
          );
        })}
      </TableBody>
    </Table>
  );
};
