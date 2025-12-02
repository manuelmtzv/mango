import { storeToRefs } from "pinia";
import { useTagStore } from "@/modules/tags/stores";

export const useTagsState = () => {
  const tagStore = useTagStore();
  const state = storeToRefs(tagStore);

  return {
    ...state,
    getTags: tagStore.getTags,
    addTag: tagStore.addTag,
    setTags: tagStore.setTags,
    getTagById: tagStore.getTagById,
    resetState: tagStore.resetState,
  };
};
