export interface IPaginateNotes {
  page: number;
  limit: number;
  search?: string;
  tags?: string[];
}

export type PaginateNotesKeys = keyof IPaginateNotes;
