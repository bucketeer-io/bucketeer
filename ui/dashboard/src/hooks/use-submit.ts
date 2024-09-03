import { useCallback, useReducer } from 'react';
import type { Reducer } from 'react';
import type { ServerError } from '@types';

interface State<Response> {
  submitting: boolean;
  response?: Response;
  error?: ServerError;
}

type Action<Response> =
  | {
      type: 'submit';
    }
  | {
      type: 'success';
      payload: Response;
    }
  | {
      type: 'error';
      payload: ServerError;
    }
  | {
      type: 'reset';
    };

const reducer: <Response>(
  state: State<Response>,
  action: Action<Response>
) => State<Response> = (state, action) => {
  switch (action.type) {
    case 'submit':
      return {
        ...state,
        submitting: true
      };
    case 'success':
      return {
        submitting: false,
        response: action.payload,
        error: undefined
      };
    case 'error':
      return {
        submitting: false,
        response: undefined,
        error: action.payload
      };
    case 'reset':
      return {
        submitting: false,
        response: undefined,
        error: undefined
      };
    default:
      return state;
  }
};

export type UseSubmitReturnType<Payload, Response> = ReturnType<
  typeof useSubmit<Payload, Response>
>;

export const useSubmit = <Payload, Response>(
  submitter: (payload: Payload) => Promise<Response>,
  options?: { handleComplete: (response: Response) => void }
) => {
  const initiatedialState: State<Response> = {
    submitting: false,
    response: undefined,
    error: undefined
  };

  const [state, dispatch] = useReducer<
    Reducer<State<Response>, Action<Response>>
  >(reducer, initiatedialState);

  const onSubmit = useCallback(async (payload: Payload) => {
    dispatch({ type: 'submit' });
    try {
      const response: Response = await submitter(payload);
      dispatch({ type: 'success', payload: response });
      if (options?.handleComplete) {
        const { handleComplete } = options;
        handleComplete(response);
      }
      return response;
    } catch (_e) {
      const error = _e as ServerError;
      dispatch({ type: 'error', payload: error });
      throw _e;
    }
  }, []);
  const onReset = useCallback(() => {
    dispatch({ type: 'reset' });
  }, []);

  return {
    ...state,
    onSubmit,
    onReset
  };
};
