export interface Tag {
  id: string;
  createdAt: string;
  updatedAt: string;
}

export interface TagCollection {
  tags: Tag[];
  cursor: string;
  totalCount: string;
}
