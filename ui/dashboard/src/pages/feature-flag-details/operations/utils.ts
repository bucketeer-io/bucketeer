import dayjs from 'dayjs';
import { v4 as uuid } from 'uuid';
import { Feature, OpsEventRateClauseOperator } from '@types';
import { ActionTypeMap, ScheduleItem } from './types';

export const createInitialDatetimeClause = (lastTime?: number) => ({
  time: dayjs(lastTime ? new Date(lastTime) : new Date())
    .set('second', 0)
    .set('millisecond', 0)
    .add(1, 'hour')
    .toDate()
});

export const createDatetimeClausesList = (lastTime?: number) => ({
  id: uuid(),
  actionType: ActionTypeMap.ENABLE,
  wasPassed: false,
  ...createInitialDatetimeClause(lastTime)
});

export const createEventRate = (feature: Feature) => ({
  variationId: feature.variations[0].id,
  goalId: '',
  minCount: 50,
  threadsholdRate: 50,
  operator: 'GREATER_OR_EQUAL' as OpsEventRateClauseOperator,
  actionType: ActionTypeMap.DISABLE
});

export const createProgressiveRollout = (feature: Feature) => ({
  template: {
    startDate: createInitialDatetimeClause().time,
    interval: 'HOURLY',
    increments: 10,
    variationId: feature.variations[0].id,
    schedulesList: []
  },
  manual: {
    variationId: feature.variations[0].id,
    schedulesList: [
      {
        scheduleId: uuid(),
        executeAt: createInitialDatetimeClause().time,
        weight: 10
      }
    ]
  }
});

export const numberToOrdinalWord = (n: number, spacingChar = ' ') => {
  try {
    const special = [
      'zeroth',
      'first',
      'second',
      'third',
      'fourth',
      'fifth',
      'sixth',
      'seventh',
      'eighth',
      'ninth',
      'tenth',
      'eleventh',
      'twelfth',
      'thirteenth',
      'fourteenth',
      'fifteenth',
      'sixteenth',
      'seventeenth',
      'eighteenth',
      'nineteenth'
    ];

    const tens = [
      '',
      '',
      'twentieth',
      'thirtieth',
      'fortieth',
      'fiftieth',
      'sixtieth',
      'seventieth',
      'eightieth',
      'ninetieth'
    ];

    const tensPrefix = [
      '',
      '',
      'twenty',
      'thirty',
      'forty',
      'fifty',
      'sixty',
      'seventy',
      'eighty',
      'ninety'
    ];
    if (n > 100) {
      console.error('The number out of range');
      return null;
    }
    if (n === 100) return 'hundredth';

    if (n < 20) return special[n];
    if (n % 10 === 0) return tens[Math.floor(n / 10)];

    return tensPrefix[Math.floor(n / 10)] + spacingChar + special[n % 10];
  } catch {
    console.error('The number out of range');
  }
};

function numberToKanji(num: number) {
  try {
    const kanjiNums = [
      '零',
      '一',
      '二',
      '三',
      '四',
      '五',
      '六',
      '七',
      '八',
      '九'
    ];
    const units = ['', '十', '百'];

    if (num === 0) return kanjiNums[0];

    let str = '';
    if (num >= 100) {
      str += kanjiNums[Math.floor(num / 100)] + units[2];
      num %= 100;
    }
    if (num >= 10) {
      const ten = Math.floor(num / 10);
      if (ten > 1) str += kanjiNums[ten];
      str += units[1];
      num %= 10;
    }
    if (num > 0) {
      str += kanjiNums[num];
    }
    return str;
  } catch {
    console.error('The number out of range');
    return null;
  }
}

export const numberToJapaneseOrdinal = (num: number) => {
  return '第' + numberToKanji(num);
};

export const handleCreateIncrement = ({
  lastSchedule,
  incrementType,
  addValue = 1,
  increment = 10
}: {
  lastSchedule: ScheduleItem;
  incrementType: dayjs.ManipulateType;
  addValue?: number;
  increment?: number;
}) => {
  const executeAt = dayjs(lastSchedule?.executeAt)
    .add(addValue, incrementType)
    .toDate();

  let weight = lastSchedule ? Number(lastSchedule.weight) : 0;

  if (weight + increment >= 100) {
    weight = 100;
  } else {
    weight = weight + increment;
  }
  return {
    scheduleId: uuid(),
    executeAt,
    weight
  };
};

export const getDateTimeDisplay = (value: string) => {
  const date = dayjs(new Date(+value * 1000)).format('YYYY/MM/DD');
  const time = dayjs(new Date(+value * 1000)).format('HH:mm');
  return {
    date,
    time
  };
};
