import { useMutation, useQueryClient } from "@tanstack/vue-query";
import { mangocatnotesApi } from "@/shared/api/mangocatnotesApi";
import { INote } from "@/modules/note/interfaces/INote";

const createNote = async (
  note: Pick<INote, "title" | "content"> & { tags: string[] },
) => {
  const { data } = await mangocatnotesApi.post<INote>("/notes", note);
  return data;
};

const editNote = async (note: Partial<INote> & { id: string }) => {
  const { data } = await mangocatnotesApi.patch(`/notes/${note.id}`, note);
  return data;
};

const deleteNote = async (id: string) => {
  const { data } = await mangocatnotesApi.delete(`/notes/${id}`);
  return data;
};

const useNoteMutation = () => {
  const queryClient = useQueryClient();

  const createNoteMutation = useMutation(createNote, {
    onSuccess: () => {
      queryClient.invalidateQueries(["notes"]);
    },
  });

  const editNoteMutation = useMutation(editNote, {
    onSuccess: () => {
      queryClient.invalidateQueries(["notes"]);
    },
  });

  const deleteNoteMutation = useMutation(deleteNote, {
    onSuccess: () => {
      queryClient.invalidateQueries(["notes"]);
    },
  });

  return {
    createNoteMutation,
    createNote: createNoteMutation.mutate,
    createNoteAsync: createNoteMutation.mutateAsync,

    editNoteMutation,
    editNote: editNoteMutation.mutate,
    editNoteAsync: editNoteMutation.mutateAsync,

    deleteNoteMutation,
    deleteNote: deleteNoteMutation.mutate,
    deleteNoteAsync: deleteNoteMutation.mutateAsync,
  };
};

export default useNoteMutation;
