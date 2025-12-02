import { ref } from "vue";
import { ITag } from "../interfaces/ITag";
import { defineStore } from "pinia";

export const useTagStore = defineStore("tagStore", () => {
  const tags = ref<ITag[]>([]);

  return {
    tags,
    getTags(): void {},
    addTag(tag: ITag): void {
      tags.value = [...tags.value, tag];
    },
    setTags(newTags: ITag[] | undefined): void {
      if (newTags) tags.value = newTags;
    },
    getTagById(id: string): ITag {
      const tag = tags.value.filter((tag) => {
        return tag.id === id;
      })[0];

      return tag || null;
    },
    resetState(): void {
      tags.value = [];
    },
  };
});
