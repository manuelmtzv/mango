<script setup lang="ts">
interface IProps {
  date: string;
  time: string;
}

defineProps<IProps>();
</script>

<template>
  <div class="flex items-center gap-4 justify-between min-h-[2rem]">
    <p>
      <span class="important">Last note update:</span> {{ date }} | {{ time }}
    </p>

    <slot />
  </div>
</template>

<style scoped></style>
