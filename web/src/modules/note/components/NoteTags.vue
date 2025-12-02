<script setup lang="ts">
import TagPillEntry from "@/modules/tags/components/TagPillEntry.vue";
import ListTransitionWrapper from "@/shared/components/animations/ListTransitionWrapper.vue";
import { ITag } from "@/modules/tags/interfaces/ITag";

type NoteTagsProps = {
  tags: ITag[];
};

const { tags } = defineProps<NoteTagsProps>();
</script>

<template>
  <section class="flex flex-col gap-2">
    <slot name="header" />

    <ListTransitionWrapper class="flex items-center gap-2 mt-1 flex-wrap">
      <TagPillEntry
        v-for="tag in tags"
        :key="tag.id"
        :tag="tag.name"
        :color="tag.color"
        class-name="cursor-default"
      />

      <span v-if="tags.length === 0" class="text-gray-500"
        >No tags attached.</span
      >
    </ListTransitionWrapper>
  </section>
</template>

<style scoped></style>
