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

export const isTimestampArraySorted = (arr: number[]) => {
  const convertToMinutes = (timestamp: number) => Math.floor(timestamp / 60000);

  for (let i = 0; i < arr.length - 1; i++) {
    if (convertToMinutes(arr[i]) >= convertToMinutes(arr[i + 1])) {
      return false;
    }
  }
  return true;
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
