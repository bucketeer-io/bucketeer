import { DatetimeClause } from '../proto/autoops/clause_pb';

export const getDatetimeClause = (value: Uint8Array) =>
  DatetimeClause.deserializeBinary(value).toObject();
