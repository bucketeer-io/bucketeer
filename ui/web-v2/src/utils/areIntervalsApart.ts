export const areIntervalsApart = (dateTimeArray, minuteDifference) => {
  for (let i = 1; i < dateTimeArray.length; i++) {
    const differenceInMinutes =
      (dateTimeArray[i] - dateTimeArray[i - 1]) / (1000 * 60); // Convert milliseconds to minutes

    if (differenceInMinutes < minuteDifference) {
      return false;
    }
  }

  return true;
};
