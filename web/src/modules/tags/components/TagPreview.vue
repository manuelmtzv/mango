<script setup lang="ts">
import TagPillEntry from "@/modules/tags/components/TagPillEntry.vue";

type TagPreviewProps = {
  name?: string;
  color: string;
};

defineProps<TagPreviewProps>();
</script>

<template>
  <div class="flex flex-col gap-2 p-4 border rounded-md">
    <h3 class="pb-1 border-b">Preview</h3>

    <TagPillEntry
      v-if="name"
      class-name="w-fit text-base mx-auto cursor-default"
      :tag="name"
      :color="color"
    />

    <span v-else class="text-gray-500">No tag name provided.</span>
  </div>
</template>

<style scoped></style>
