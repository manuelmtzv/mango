import { useQuery } from "@tanstack/vue-query";
import { mangocatnotesApi } from "@/shared/api/mangocatnotesApi";
import { INote } from "@/modules/note/interfaces/INote";

interface Params {
  id: string;
}

const getNote = async (id: string) => {
  const { data } = await mangocatnotesApi.get<INote>(`/notes/${id}`);
  return data;
};

const useNote = (params: Params) => {
  const { data, error, isLoading, isFetching, refetch } = useQuery(
    ["note", params.id],
    () => getNote(params.id),
    {
      retry: false,
      refetchOnWindowFocus: false,
    },
  );

  return {
    data,
    error,
    isLoading,
    isFetching,
    refetch,
  };
};

export default useNote;
