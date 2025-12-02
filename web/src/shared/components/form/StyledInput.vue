<script setup lang="ts">
import { InputHTMLAttributes } from "vue";

interface IProps extends InputHTMLAttributes {
  value: string;
  debounce: number;
}

defineProps<IProps>();

const emit = defineEmits(["update:modelValue"]);
</script>

<template>
  <label class="label">
    <span>Title:</span>
    <input
      :name="name"
      :placeholder="placeholder"
      class="input"
      type="text"
      :model-value="value"
      @update:modelValue="emit('update:modelValue', $event)"
    />
  </label>
</template>

<style scoped></style>
