<script setup lang="ts">
import { useTimer } from "@shared/composables/useTimer";
import { useMotions } from "@vueuse/motion";

const showTimeLimit = 24;
const hightOpacityLimit = 8;
const motions = useMotions();
const { time: showTime, resetTimer } = useTimer({
  limit: showTimeLimit,
  active: true,
  defaultTime: showTimeLimit,
});

interface Expose {
  resetTimer: () => void;
}

defineExpose<Expose>({ resetTimer });
</script>

<template>
  <Transition @leave="(_, done) => motions.fade.leave(done)">
    <div
      v-if="showTime < showTimeLimit"
      v-motion="'fade'"
      v-tooltip="
        showTime < hightOpacityLimit ? 'Edited now' : 'Edited recently'
      "
      :initial="{ opacity: 0 }"
      :enter="{ opacity: 1 }"
      :leave="{ opacity: 0 }"
      class="h-fi flex items-center gap-2"
    >
      <span
        :class="[
          'material-symbols-outlined icon',
          showTime < hightOpacityLimit ? 'opacity-100' : 'opacity-50',
        ]"
      >
        cloud_done
      </span>
    </div>
  </Transition>
</template>

<style scoped>
.icon {
  @apply transition-all duration-300 ease-linear leading-none;
}
.recent {
  @apply text-green-500;
}
</style>
