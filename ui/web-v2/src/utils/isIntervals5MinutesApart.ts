export const isIntervals5MinutesApart = (dateTimeArray) => {
  for (let i = 1; i < dateTimeArray.length; i++) {
    const timeDifference =
      (dateTimeArray[i] - dateTimeArray[i - 1]) / (1000 * 60); // Convert milliseconds to minutes

    if (timeDifference < 5) {
      return false;
    }
  }

  return true;
};
