import { grpc } from '@improbable-eng/grpc-web';
import { BrowserHeaders } from 'browser-headers';
import * as jspb from 'google-protobuf';

import {
  getSelectedLanguage,
  LanguageTypes
} from '../lang/getSelectedLanguage';
import { getToken } from '../storage/token';

type MetaData = {
  authorization: string;
  'accept-language': LanguageTypes;
};

export const isSuccess = (output: grpc.UnaryOutput<jspb.Message>): boolean =>
  successResponse(output);

const successResponse = (response: grpc.UnaryOutput<jspb.Message>): boolean => {
  const { message, status } = response;
  return status === grpc.Code.OK && !!message;
};

export const getMetaData = (): MetaData => {
  const token = getToken();
  return {
    authorization: `bearer ${token ? token.accessToken : ''}`,
    'accept-language': getSelectedLanguage()
  };
};

export const getMetaDataForClient = (): BrowserHeaders => {
  const token = getToken();
  return new BrowserHeaders({
    authorization: `bearer ${token ? token.accessToken : ''}`,
    'accept-language': getSelectedLanguage()
  });
};
