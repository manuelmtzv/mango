<script setup lang="ts">
import { cn } from "@/shared/utils/cn";

type TagPillEntryProps = {
  tag: string;
  selected?: boolean;
  className?: string;
  tooltip?: string;
  color?: string;
};

const props = defineProps<TagPillEntryProps>();
</script>

<template>
  <div
    :class="
      cn(
        'text-sm px-2.5 py-0.5 rounded-full text-gray-900 border relative cursor-pointer select-none transition-opacity hover:opacity-80 ease-in-out',
        selected && 'border-yellow-500',
        props.className,
      )
    "
    :style="{
      backgroundColor: props.color,
    }"
    v-tooltip="props.tooltip"
  >
    {{ tag }}
    <span
      :class="[
        'w-3 h-3 rounded-full bg-yellow-300 absolute -top-1 -right-1 leading-none border border-yellow-600',
        selected ? 'block' : 'hidden',
      ]"
    >
    </span>
  </div>
</template>

<style scoped lang="css"></style>
