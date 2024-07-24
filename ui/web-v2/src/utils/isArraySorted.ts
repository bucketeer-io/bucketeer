export const isArraySorted = (arr: number[]) => {
  for (let i = 0; i < arr.length - 1; i++) {
    if (arr[i] > arr[i + 1]) {
      return false;
    }
  }
  return true;
};

export const isTimestampArraySorted = (arr: number[]) => {
  const convertToMinutes = (timestamp: number) => Math.floor(timestamp / 60000);

  for (let i = 0; i < arr.length - 1; i++) {
    if (convertToMinutes(arr[i]) >= convertToMinutes(arr[i + 1])) {
      return false;
    }
  }
  return true;
};
