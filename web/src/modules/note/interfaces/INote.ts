import { ITag } from "@/modules/tags/interfaces/ITag";

export interface INote {
  id: string;
  title?: string;
  content: string;
  createdAt: string;
  updatedAt: string;
  tagIDs?: string[];
  tags?: ITag[];
}
