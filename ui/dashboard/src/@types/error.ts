import { ErrorComponentProps } from '@tanstack/react-router';

export interface Issue {
  expected: unknown;
  received: string;
  code: unknown;
  path: unknown;
  message: string;
}

export interface ErrorComponentExpandProps extends ErrorComponentProps {
  error: ErrorComponentProps['error'] & {
    cause: {
      issues: Issue[];
    };
  };
}
