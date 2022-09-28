import { SerializedError } from '@reduxjs/toolkit';
import { FC, useState } from 'react';
import { useIntl } from 'react-intl';
import { shallowEqual, useSelector } from 'react-redux';

import { messages } from '../..//lang/messages';
import {
  ANALYSIS_LIST_PAGE_SIZE,
  ANALYSIS_USER_METADATA_REGEX,
} from '../../constants/analysis';
import { intl } from '../../lang';
import { AppState } from '../../modules';
import { selectById as selectFeatureById } from '../../modules/features';
import { Cell, Row } from '../../proto/eventcounter/table_pb';
import { Feature } from '../../proto/feature/feature_pb';
import { classNames } from '../../utils/css';
import { ListSkeleton } from '../ListSkeleton';
import { Pagination } from '../Pagination';

interface AnalysisTableProps {
  featureId?: string;
}

export const AnalysisTable: FC<AnalysisTableProps> = ({ featureId }) => {
  const [page, setPage] = useState<number>(1);
  const [feature, getFeatureError] = useSelector<
    AppState,
    [Feature.AsObject | undefined, SerializedError | null]
  >(
    (state) => [
      selectFeatureById(state.features, featureId),
      state.features.getFeatureError,
    ],
    shallowEqual
  );
  const headers = useSelector<AppState, Row.AsObject>(
    (state) => state.goalCounts.headers,
    shallowEqual
  );
  const rows = useSelector<AppState, Row.AsObject[]>(
    (state) => state.goalCounts.rows,
    shallowEqual
  );
  const loading = useSelector<AppState, boolean>(
    (state) => state.goalCounts.loading,
    shallowEqual
  );
  const variationIdx = headers?.cellsList?.findIndex(
    (cell) => cell.value == 'Variation'
  );
  const variations = new Map<string, string>();
  feature?.variationsList.forEach((v) => {
    variations.set(v.id, v.value);
  });
  if (loading) {
    return (
      <div className="m-5 border border-gray-300 bg-white rounded">
        <ListSkeleton />
      </div>
    );
  }
  if (!(rows && headers)) {
    return null;
  }
  return (
    <div className="bg-white rounded-md border">
      <Table>
        <TableHeader headers={headers} />
        <TableBody>
          {rows
            .slice(
              (page - 1) * ANALYSIS_LIST_PAGE_SIZE,
              page * ANALYSIS_LIST_PAGE_SIZE
            )
            .map((row, i) => {
              return (
                <TableRow key={i.toString()}>
                  {row.cellsList.map((cell, j) => {
                    let textLeft = false;
                    let value: string;
                    if (cell.type == Cell.Type.DOUBLE) {
                      value = Number.isNaN(cell.valuedouble)
                        ? 'n/a'
                        : cell.valuedouble.toFixed(1);
                    } else {
                      textLeft = true;
                      // Replace variation ID with variation value.
                      if (variationIdx == j) {
                        const variationValue = variations.get(cell.value);
                        if (!variationValue) {
                          value = `(deleted) ${cell.value}`;
                        } else {
                          value = variationValue;
                        }
                      } else {
                        value = cell.value;
                      }
                    }
                    return (
                      <TableCell
                        key={i.toString() + '-' + j.toString()}
                        textLeft={textLeft}
                      >
                        {' '}
                        {value}{' '}
                      </TableCell>
                    );
                  })}
                </TableRow>
              );
            })}
        </TableBody>
      </Table>
      <Pagination
        maxPage={Math.ceil(rows.length / ANALYSIS_LIST_PAGE_SIZE)}
        currentPage={page}
        onChange={setPage}
      />
    </div>
  );
};

interface TableHeaderProps {
  headers: Row.AsObject;
}

const TableHeader: FC<TableHeaderProps> = ({ headers }) => {
  const { formatMessage: f } = useIntl();
  return (
    <thead className="bg-gray-50">
      <tr>
        {headers.cellsList.map((header) => {
          let value = header.value;
          const match = value.match(ANALYSIS_USER_METADATA_REGEX);
          if (match && match?.length > 1) {
            value = match[1];
          }
          if (value === 'tag') {
            value = f(messages.tags);
          }
          return (
            <td
              key={header.value}
              className={classNames(
                'py-1 px-5',
                'rounded-t-md',
                'text-left text-xs',
                'font-medium text-gray-500 tracking-wider'
              )}
            >
              <div className={classNames('flex flex-row items-center')}>
                {headerMessage(value)}
              </div>
            </td>
          );
        })}
      </tr>
    </thead>
  );
};

function headerMessage(name: string): string {
  let msg = '';
  switch (name) {
    case 'Variation':
      msg = intl.formatMessage(messages.analysis.variation);
      break;
    case 'Evaluation user':
      msg = intl.formatMessage(messages.analysis.evaluationUser);
      break;
    case 'Evaluation total':
      msg = intl.formatMessage(messages.analysis.evaluationTotal);
      break;
    case 'Goal user':
      msg = intl.formatMessage(messages.analysis.goalUser);
      break;
    case 'Goal total':
      msg = intl.formatMessage(messages.analysis.goalTotal);
      break;
    case 'Goal value total':
      msg = intl.formatMessage(messages.analysis.goalValueTotal);
      break;
    case 'Goal value mean':
      msg = intl.formatMessage(messages.analysis.goalValueMean);
      break;
    case 'Goal value variance':
      msg = intl.formatMessage(messages.analysis.goalValueVariance);
      break;
    case 'Conversion rate':
      msg = intl.formatMessage(messages.analysis.conversionRate);
      break;
    default:
      msg = name;
  }
  return msg;
}

export interface TableProps {}

export const Table: FC<TableProps> = ({ children }) => {
  return (
    <div className={classNames('w-full rounded-md')}>
      <table
        className={classNames(
          'min-w-full table-auto leading-normal rounded-md'
        )}
      >
        {children}
      </table>
    </div>
  );
};

export interface TableBodyProps {}

export const TableBody: FC<TableBodyProps> = ({ children }) => {
  return <tbody className="text-sm text-gray-600 rounded-md">{children}</tbody>;
};

export interface TableRowProps {}

export const TableRow: FC<TableRowProps> = ({ children }) => {
  return <tr className="border-t">{children}</tr>;
};

export interface TableCellProps {
  textLeft?: boolean;
}

export const TableCell: FC<TableCellProps> = ({ children, textLeft }) => {
  return (
    <td
      className={classNames(
        'w-[1%] px-5 border-b',
        textLeft ? 'text-left' : 'text-right'
      )}
    >
      {children}
    </td>
  );
};
