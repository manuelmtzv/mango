import { useMutation, useQueryClient } from "@tanstack/vue-query";
import { mangocatnotesApi } from "@/shared/api/mangocatnotesApi";

type EditNoteTags = {
  id: string;
  tags: string[];
};

const editNoteTags = async ({ id, tags }: EditNoteTags) => {
  const { data } = await mangocatnotesApi.patch(`/notes/${id}/tags`, {
    tags,
  });
  return data;
};

const useNoteTagsMutation = () => {
  const queryClient = useQueryClient();

  const editNoteTagsMutation = useMutation(editNoteTags, {
    onSuccess: () => {
      queryClient.invalidateQueries(["notes"]);
      queryClient.invalidateQueries([
        "note",
        editNoteTagsMutation.variables.value?.id,
      ]);
    },
  });

  return {
    editNoteTagsMutation,
    editNoteTags: editNoteTagsMutation.mutate,
    editNoteTagsAsync: editNoteTagsMutation.mutateAsync,
  };
};

export default useNoteTagsMutation;
