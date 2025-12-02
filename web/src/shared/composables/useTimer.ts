import { computed, ref, watch } from "vue";

interface Props {
  defaultTime?: number;
  limit?: number;
  active?: boolean;
}

export const useTimer = ({ active, limit, defaultTime }: Props) => {
  const time = ref(defaultTime ?? 0);
  const resetTimer = () => (time.value = 0);

  const limitReached = computed(() => {
    return limit && time.value >= limit;
  });

  watch(
    time,
    () => {
      if (active && !limitReached.value) {
        const interval = setInterval(() => {
          time.value = time.value + 1;
        }, 1000);

        setTimeout(() => {
          clearInterval(interval);
        }, 1000);
      }
    },
    {
      immediate: true,
    },
  );

  return {
    time,
    resetTimer,
  };
};
