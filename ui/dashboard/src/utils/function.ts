import { ENVIRONMENT_WITH_EMPTY_ID } from 'constants/app';
import dayjs from 'dayjs';
import * as yaml from 'js-yaml';
import { Environment } from '@types';
import { areJsonStringsEqual } from './converts';

export const copyToClipBoard = (text: string) => {
  if (navigator.clipboard) {
    navigator.clipboard.writeText(text);
  } else {
    const textArea = document.createElement('textarea');
    textArea.value = text;
    textArea.style.top = '0';
    textArea.style.left = '0';
    textArea.style.position = 'fixed';
    document.body.appendChild(textArea);
    textArea.focus();
    textArea.select();
    try {
      document.execCommand('copy');
    } catch (error) {
      console.error(error);
    } finally {
      document.body.removeChild(textArea);
    }
  }
};

export const isArraySorted = (arr: number[]) => {
  for (let i = 0; i < arr.length - 1; i++) {
    if (arr[i] > arr[i + 1]) {
      return false;
    }
  }
  return true;
};

export const isTimestampArraySorted = (
  arr: number[]
): {
  isSorted: boolean;
  index: number;
  isEquals: boolean;
} => {
  const convertToMinutes = (timestamp: number) => Math.floor(timestamp / 60000);

  for (let i = 0; i < arr.length - 1; i++) {
    if (
      arr[i + 1] &&
      convertToMinutes(arr[i]) >= convertToMinutes(arr[i + 1])
    ) {
      return {
        isSorted: false,
        index: i,
        isEquals: convertToMinutes(arr[i]) === convertToMinutes(arr[i + 1])
      };
    }
  }
  return {
    isSorted: true,
    index: -1,
    isEquals: false
  };
};

export const areIntervalsApart = (
  dateTimeArray: number[],
  minuteDifference: number
) => {
  for (let i = 1; i < dateTimeArray.length; i++) {
    const differenceInMinutes =
      (dateTimeArray[i] - dateTimeArray[i - 1]) / (1000 * 60); // Convert milliseconds to minutes
    if (differenceInMinutes < minuteDifference) {
      return false;
    }
  }

  return true;
};

export const hasDuplicateTimestamps = (arr: number[]) => {
  const convertToMinutes = (timestamp: number) => Math.floor(timestamp / 60000);
  const seenTimestamps = new Set<number>();

  for (const timestamp of arr) {
    const minutes = convertToMinutes(timestamp);
    if (seenTimestamps.has(minutes)) {
      return true;
    }
    seenTimestamps.add(minutes);
  }
  return false;
};

export const isSameOrBeforeDate = (date: Date, conditionDate = new Date()) => {
  return (
    dayjs(date).isSame(conditionDate) || dayjs(date).isBefore(conditionDate)
  );
};

export const onFormatEnvironments = (environments: Environment[]) => {
  let emptyEnvironmentId = '';
  const formattedEnvironments = environments.map((item, index) => {
    if (!item.id) emptyEnvironmentId = `${ENVIRONMENT_WITH_EMPTY_ID}${index}`;
    return {
      ...item,
      id: item.id ? item.id : emptyEnvironmentId
    };
  });
  return { emptyEnvironmentId, formattedEnvironments };
};

export const checkEnvironmentEmptyId = (environmentId: string) =>
  environmentId.includes(ENVIRONMENT_WITH_EMPTY_ID) ? '' : environmentId;

export const onChangeFontWithLocalized = (isLanguageJapanese: boolean) => {
  const htmlElement = document.getElementsByTagName('html')[0];
  if (htmlElement) {
    htmlElement.classList[isLanguageJapanese ? 'add' : 'remove'](
      'japanese-language'
    );
    htmlElement.style.setProperty(
      'font-family',
      isLanguageJapanese ? 'Noto Sans JP, sans-serif' : 'Sofia Pro, sans-serif',
      'important'
    );
  }
};

export const checkFieldDirty = (obj: { [key: string]: boolean }): boolean => {
  if (!obj) return false;
  if (typeof obj === 'boolean' && obj === true) return true;
  for (const value of Object.values(obj)) {
    if (typeof value === 'boolean' && value === true) return true;
    if (Array.isArray(value)) {
      for (const item of value) {
        if (typeof item === 'object') {
          if (checkFieldDirty(item as { [key: string]: boolean })) return true;
        }
      }
    }
    if (typeof value === 'object' && value !== null) {
      if (checkFieldDirty(value)) return true;
    }
  }
  return false;
};

export const normalizeYamlIndentation = (yamlString: string) => {
  try {
    const data = yaml.load(yamlString);
    return yaml.dump(data);
  } catch {
    return yamlString;
  }
};

export const isUniqueValue = (text1: string, text2: string, type: string) => {
  switch (type) {
    case 'YAML':
      return (
        normalizeYamlIndentation(text1) !== normalizeYamlIndentation(text2)
      );
    case 'JSON':
      return !areJsonStringsEqual(text1, text2);
    default:
      return text1.trim() !== text2.trim();
  }
};
