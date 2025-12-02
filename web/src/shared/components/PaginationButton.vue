<script setup lang="ts">
import { cn } from "@/shared/utils/cn";
import { ButtonHTMLAttributes, computed } from "vue";

interface PaginationButtonProps extends ButtonHTMLAttributes {
  page: number;
  currentPage: number;
}

const props = defineProps<PaginationButtonProps>();

const isCurrentPage = computed(() => props.page === props.currentPage);
</script>

<template>
  <button
    :v-bind="props"
    :class="
      cn(
        'p-2 py-1 rounded-md border-2 select-none',
        isCurrentPage && 'text-gray-500 cursor-default',
      )
    "
    :disabled="isCurrentPage"
  >
    {{ props.page }}
  </button>
</template>

<style scoped></style>
