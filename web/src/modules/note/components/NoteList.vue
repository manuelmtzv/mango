<script setup lang="ts">
import Loading from "vue-loading-overlay";
import useNotes from "@/modules/note/composables/useNotes";
import NoteEntry from "./NoteEntry.vue";
import { computed } from "vue";
import useNotePaginationState from "@/modules/note/composables/useNotePaginationState";
import LoadingSpin from "@/shared/components/LoadingSpin.vue";

const { sortedByTime: renderingNotes, isLoading, isRefetching } = useNotes();
const { isLoading: loadingPagination } = useNotePaginationState();

const loading = computed(() => isLoading.value || loadingPagination.value);
</script>

<template>
  <div class="list">
    <div class="flex items-center justify-between">
      <h2
        :class="['list-title']"
        v-if="renderingNotes.length > 0 && !isLoading"
      >
        Showing:
        <p class="font-normal">{{ renderingNotes.length }} notes</p>

        <LoadingSpin class="w-4" :dark="true" :when="isRefetching" />
      </h2>

      <p v-else>No notes to show</p>

      <slot name="actions" />
    </div>

    <MasonryWall :items="renderingNotes" :column-width="300" :gap="16">
      <template #default="{ item }">
        <NoteEntry :note="item" />
      </template>
    </MasonryWall>

    <p v-if="isLoading">Loading notes...</p>

    <Loading v-model:active="loading" />
  </div>

  <slot name="footer" />
</template>

<style scoped lang="css">
.list {
  @apply flex flex-col gap-4 w-full;
}
.list-title {
  @apply text-lg font-semibold inline-flex gap-1.5 items-center;
}
.list-container {
  @apply grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4;
}
</style>
