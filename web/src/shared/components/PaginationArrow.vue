<script setup lang="ts">
import { cn } from "@/shared/utils/cn";
import { ButtonHTMLAttributes } from "vue";

export interface PaginationButtonProps
  extends /* @vue-ignore */ ButtonHTMLAttributes {
  buttonType: "prev" | "next";
  allowed: boolean;
}

const props = defineProps<PaginationButtonProps>();
</script>

<template>
  <button
    v-bind="props"
    :class="
      cn(
        'material-symbols-outlined px-2 py-2 rounded-full inline-flex hover:bg-navigation-hover transition-colors duration-300',
        !props.allowed && 'text-gray-400 cursor-not-allowed',
        props.class,
      )
    "
    :disabled="!props.allowed"
  >
    <span class="material-symbols-outlined">
      {{ props.buttonType === "prev" ? "chevron_left" : "chevron_right" }}
    </span>
  </button>
</template>

<style scoped></style>
