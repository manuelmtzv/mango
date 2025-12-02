import useNotes from "@/modules/note/composables/useNotes";
import { watch } from "vue";
import { useToast } from "vue-toast-notification";
import useNotePaginationState from "./useNotePaginationState";

const useNotePagination = () => {
  const { data, refetch } = useNotes();
  const toast = useToast();

  const { setIsLoading, setPaginateValue, ...state } = useNotePaginationState();

  async function refetchPagination() {
    setIsLoading(true);
    await refetch();
  }

  watch(
    () => state.paginate,
    async () => {
      await refetchPagination();
    },
    {
      deep: true,
    },
  );

  watch(
    () => data.value?.data,
    (value) => {
      if (!value) return;

      const { paginate, totalPages } = state;

      if (!value.length && paginate.value.page > 1) {
        toast.info(
          "No notes found on last page; redirecting to the last one.",
          {
            duration: 5000,
          },
        );

        setPaginateValue("page", totalPages.value);
      }
    },
  );

  return {
    refetchPagination,
  };
};

export default useNotePagination;
