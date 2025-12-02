<script setup lang="ts">
import { cn } from "@/shared/utils/cn";
import { HTMLAttributes } from "vue";

interface PaginationNavigateProps extends /* @vue-ignore */ HTMLAttributes {
  totalPages: number;
  currentPage: number;
  setPage: (page: number) => void;
}

const props = defineProps<PaginationNavigateProps>();
</script>

<template>
  <nav :v-bind="props" :class="cn('flex items-center gap-3', props.class)">
    <label for="pagination-select">Current page: </label>

    <select
      id="pagination-select"
      :value="props.currentPage"
      @change="
        (event) =>
          props.setPage(Number((event.target as HTMLInputElement).value))
      "
      class="p-2 py-0.5 rounded-md border-2 focus:outline-none bg-gray-100 hover:bg-gray-200 transition-colors duration-300"
    >
      <option
        v-for="page in props.totalPages"
        :key="page"
        :value="page"
        :selected="page === props.currentPage"
      >
        {{ page }}
      </option>
    </select>
  </nav>
</template>

<style scoped></style>
