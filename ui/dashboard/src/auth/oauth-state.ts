import { v4 as uuid } from 'uuid';

// Generates an unguessable value for the OAuth `state` parameter.
export const generateOAuthState = (): string => uuid();

// The `state` returned on the OAuth callback is accepted only when it is a
// non-empty string equal to the value stored when the sign-in was started.
// A missing value on either side (or a value that is not a string, such as a
// repeated query parameter) is rejected.
export const isValidOAuthState = (
  state: unknown,
  expectedState: string | undefined
): boolean =>
  typeof state === 'string' && state !== '' && state === expectedState;
