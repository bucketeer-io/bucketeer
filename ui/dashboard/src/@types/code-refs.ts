export enum RepositoryType {
  // UNSPECIFIED = 'REPOSITORY_TYPE_UNSPECIFIED',
  GITHUB = 'GITHUB',
  GITLAB = 'GITLAB',
  BITBUCKET = 'BITBUCKET'
  // CUSTOM = 'CUSTOM'
}

export interface CodeRefCollection {
  codeReferences: CodeReference[];
  cursor: string;
  totalCount: string;
}

export interface CodeReference {
  id: string;
  featureId: string;
  filePath: string;
  lineNumber: number;
  codeSnippet: string;
  contentHash: string;
  aliases: string[];
  repositoryName: string;
  repositoryOwner: string;
  repositoryType: RepositoryType;
  repositoryBranch: string;
  commitHash: string;
  environmentId: string;
  createdAt: string;
  updatedAt: string;
  sourceUrl: string;
  branchUrl: string;
  fileExtension: string;
}
