<script setup lang="ts">
import { cn } from "@/shared/utils/cn";
import LoadingSpin from "@shared/components/LoadingSpin.vue";
import { ButtonHTMLAttributes } from "vue";

interface Props extends /* @vue-ignore */ ButtonHTMLAttributes {
  label?: string;
  loading?: boolean;
}

const props = defineProps<Props>();
</script>

<template>
  <button
    :class="
      cn(
        'flex items-center gap-2 p-2 h-fit rounded-md border transition-colors duration-300',
        props.class,
      )
    "
  >
    <p v-if="label">{{ label }}</p>
    <slot v-if="!label" />

    <LoadingSpin :when="loading" class="!h-5 !w-5" />
  </button>
</template>
