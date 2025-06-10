import { ErrorComponentProps } from '@tanstack/react-router';
import dayjs from 'dayjs';
import { ErrorComponentExpandProps } from '@types';
import { isJsonString } from './converts';

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

export const handleGetErrorMessage = (error?: ErrorComponentProps) => {
  if (!error) return null;
  const _error = error as ErrorComponentExpandProps;

  const message =
    _error?.error.cause?.issues[0]?.message ||
    (isJsonString(error.error?.message) &&
      JSON.parse(error.error?.message)?.message) ||
    '';
  return {
    message,
    reset: error?.reset
  };
};
