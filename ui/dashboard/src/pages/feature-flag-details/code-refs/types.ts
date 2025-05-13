import { RepositoryType } from '@types';

export interface CodeRefFilters {
  cursor: string;
  pageSize: number;
  repositoryName: string;
  repositoryOwner: string;
  repositoryType: RepositoryType;
  repositoryBranch: string;
  fileExtension: string;
}
