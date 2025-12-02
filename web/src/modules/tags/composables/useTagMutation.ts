import { useMutation, useQueryClient } from "@tanstack/vue-query";
import { mangocatnotesApi } from "@/shared/api/mangocatnotesApi";
import { ITag } from "@/modules/tags/interfaces/ITag";

const createTag = async (tag: Pick<ITag, "name" | "color">) => {
  const { data } = await mangocatnotesApi.post<ITag>("/tags", tag);
  return data;
};

const editTag = async (tag: Partial<ITag>) => {
  const { data } = await mangocatnotesApi.patch<ITag>(`/tags/${tag.id}`, tag);
  return data;
};

const deleteTag = async (tagId: string) => {
  const { data } = await mangocatnotesApi.delete(`/tags/${tagId}`);
  return data;
};

export const useTagMutation = () => {
  const queryClient = useQueryClient();

  const editTagMutation = useMutation(editTag, {
    onSuccess: () => {
      queryClient.invalidateQueries(["tags"]);
    },
  });

  const createTagMutation = useMutation(createTag, {
    onSuccess: () => {
      queryClient.invalidateQueries(["tags"]);
    },
  });

  const deleteTagMutation = useMutation(deleteTag, {
    onSuccess: () => {
      queryClient.invalidateQueries(["tags"]);
    },
  });

  return {
    createTagMutation,
    createTag: createTagMutation.mutate,
    createTagAsync: createTagMutation.mutateAsync,

    editTagMutation,
    editTag: editTagMutation.mutate,
    editTagAsync: editTagMutation.mutateAsync,

    deleteTagMutation,
    deleteTag: deleteTagMutation.mutate,
    deleteTagAsync: deleteTagMutation.mutateAsync,
  };
};
