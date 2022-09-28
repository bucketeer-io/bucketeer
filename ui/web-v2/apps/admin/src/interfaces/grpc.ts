export interface RequestResult<T> {
  isOk: boolean;
  id?: string;
  errorMessage?: string;
  response?: T;
}
