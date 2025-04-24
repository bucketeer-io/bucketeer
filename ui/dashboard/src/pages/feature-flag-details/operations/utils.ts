import dayjs from 'dayjs';
import { v4 as uuid } from 'uuid';
import { Feature } from '@types';
import { ActionTypeMap } from './types';

export const createInitialDatetimeClause = (lastTime?: number) => ({
  time: dayjs(lastTime ? new Date(lastTime) : new Date())
    .add(1, 'hour')
    .toDate()
});

export const createDatetimeClausesList = (lastTime?: number) => ({
  id: uuid(),
  actionType: ActionTypeMap.ENABLE,
  ...createInitialDatetimeClause(lastTime)
});

export const createEventRate = (feature: Feature) => ({
  variation: feature.variations[0].id,
  goal: '',
  minCount: 50,
  threadsholdRate: 50,
  operator: '<=',
  actionType: ActionTypeMap.ENABLE
});

export const createProgressiveRollout = (feature: Feature) => ({
  template: {
    datetime: createInitialDatetimeClause(),
    interval: '1',
    increments: 10,
    variationId: feature.variations[0].id,
    schedulesList: []
  },
  manual: {
    variationId: feature.variations[0].id,
    schedulesList: [
      {
        executeAt: createInitialDatetimeClause(),
        weight: 10
      }
    ]
  }
});
