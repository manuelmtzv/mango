import { useNotePaginationStore } from "@modules/note/stores/notePaginationStore";
import { storeToRefs } from "pinia";

const useNotePaginationState = () => {
  const notePaginationStore = useNotePaginationStore();
  const state = storeToRefs(notePaginationStore);
  const {
    setPaginate,
    setPaginateValue,
    setPage,
    setIsLoading,
    setTotalPages,
  } = notePaginationStore;

  return {
    ...state,

    setPaginate,
    setPaginateValue,
    setPage,
    setIsLoading,
    setTotalPages,
  };
};

export default useNotePaginationState;
