<script setup lang="ts">
import Loading from "vue-loading-overlay";
import { useTags } from "@/modules/tags/composables/useTags";
import { MAX_TAGS_PER_USER } from "@/config/constants";
import TagCreateModal from "./TagCreateModal.vue";
import TagList from "@/modules/tags/components/TagList.vue";

const { tags, isLoading } = useTags();
</script>

<template>
  <div class="flex flex-col gap-4">
    <div class="flex items-center justify-between gap-4 border-b pb-4">
      <div>
        <h2 class="font-semibold text-lg">
          Tags administration
          <span
            v-if="!isLoading && tags.length"
            v-tooltip="
              `Every user can have a maximum of ${MAX_TAGS_PER_USER} tags`
            "
            class="text-sm cursor-default"
          >
            {{ `(${tags.length}/${MAX_TAGS_PER_USER})` }}
          </span>
        </h2>

        <p class="text-sm text-gray-800">
          Select a tag to edit its name or color.
        </p>
      </div>

      <TagCreateModal :disabled="tags.length >= MAX_TAGS_PER_USER" />
    </div>

    <TagList v-if="tags.length" :tags="tags" />

    <p v-else-if="!tags.length && !isLoading" class="text-gray-500">
      No tags to show already.
    </p>

    <p v-if="isLoading" class="text-sm text-gray-800">Loading...</p>

    <Loading v-model:active="isLoading" />
  </div>
</template>

<style scoped></style>
