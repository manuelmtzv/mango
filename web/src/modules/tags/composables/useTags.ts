import { ITag } from "@/modules/tags/interfaces/ITag";
import { mangocatnotesApi } from "@/shared/api/mangocatnotesApi";
import { useTagsState } from "./useTagsState";
import { useQuery } from "@tanstack/vue-query";
import { watch } from "vue";
import { IListResponse } from "@/shared/interfaces/IListResponse";

const getTags = async (): Promise<IListResponse<ITag>> => {
  const { data } = await mangocatnotesApi.get<IListResponse<ITag>>("/tags");

  return data;
};

export const useTags = () => {
  const { tags, setTags } = useTagsState();

  const { data, isLoading, isRefetching, refetch } = useQuery(
    ["tags"],
    getTags,
    {
      enabled: true,
    },
  );

  watch(
    data,
    (value) => {
      if (value) {
        setTags(value.data);
      }
    },
    {
      immediate: true,
    },
  );

  return {
    data,
    tags,
    isLoading,
    isRefetching,
    refetch,
  };
};
