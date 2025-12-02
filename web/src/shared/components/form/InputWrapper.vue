<script setup lang="ts">
import { cn } from "@/shared/utils/cn";

type InputWrapperProps = {
  label?: string;
  error?: string;
  className?: string;
};

defineProps<InputWrapperProps>();
</script>

<template>
  <label :class="cn('inline-flex flex-col gap-2', $props.className)">
    <slot v-if="$slots['label']" name="label" />

    <p v-else>{{ label ?? "" }}</p>

    <slot />

    <span v-if="error" class="text-red-500 text-sm font-normal">{{
      error
    }}</span>
  </label>
</template>

<style scoped></style>
