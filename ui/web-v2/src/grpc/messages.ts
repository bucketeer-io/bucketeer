import { BrowserHeaders } from 'browser-headers';
import { Nullable } from 'option-t/lib/Nullable';
import { isNotNull, isNull } from 'option-t/lib/Nullable/Nullable';

import { LocalizedMessage } from '../proto/external/googleapis/googleapis/83e756a66b80b072bd234abcfe89edf459090974/google/rpc/error_details_pb';
import { Status } from '../proto/external/googleapis/googleapis/83e756a66b80b072bd234abcfe89edf459090974/google/rpc/status_pb';

export const DEFAULT_ERROR = 'unknown error';

export const STATUS_DETAILS_BIN_KEY = 'grpc-status-details-bin';

export const LOCALIZED_MESSAGE_TYPE = 'google.rpc.LocalizedMessage';

export interface ServiceError {
  metadata: BrowserHeaders;
}

export const extractErrorMessage = (error: Nullable<ServiceError>): string => {
  if (isNull(error)) {
    return DEFAULT_ERROR;
  }

  const status = deserializeStatus(error.metadata);
  if (isNull(status)) {
    return DEFAULT_ERROR;
  }
  const localizedMessage = deserializeLocalizedMessage(status);
  if (isNull(localizedMessage)) {
    return DEFAULT_ERROR;
  }

  return localizedMessage.getMessage();
};

export const deserializeStatus = (
  metadata: BrowserHeaders
): Nullable<Status> => {
  const statusDetailsBins = metadata.get(STATUS_DETAILS_BIN_KEY);
  if (statusDetailsBins.length === 0) {
    return null;
  }

  return Status.deserializeBinary(
    stringToUint8Array(atob(statusDetailsBins[0]))
  );
};

export const deserializeLocalizedMessage = (
  status: Status
): Nullable<LocalizedMessage> => {
  const localizedMessages = status
    .getDetailsList()
    .map((detail) => {
      if (detail.getTypeName() !== LOCALIZED_MESSAGE_TYPE) {
        return null;
      }
      return detail.unpack(
        LocalizedMessage.deserializeBinary,
        detail.getTypeName()
      );
    })
    .filter(isNotNull);

  if (localizedMessages.length === 0) {
    return null;
  }

  return localizedMessages[0];
};

export const stringToUint8Array = (str: string): Uint8Array => {
  const buf = new ArrayBuffer(str.length);
  const bufView = new Uint8Array(buf);
  for (let i = 0; i < str.length; i++) {
    bufView[i] = str.charCodeAt(i);
  }
  return bufView;
};
