<template>
  <textarea
    class="input"
    ref="textareaRef"
    v-model="internalValue"
    @input="handleInput"
  ></textarea>
</template>

<script lang="ts">
import { ref, watch, nextTick, SetupContext, onMounted } from "vue";

export default {
  props: {
    value: {
      type: String,
      required: true,
    },
  },
  setup(props: { value: string }, context: SetupContext) {
    const internalValue = ref(props.value);
    const textareaRef = ref<HTMLTextAreaElement | null>(null);

    onMounted(() => {
      adjustTextareaHeight();
    });

    const adjustTextareaHeight = () => {
      if (textareaRef.value) {
        textareaRef.value.style.height = "auto";
        textareaRef.value.style.height = `${textareaRef.value.scrollHeight}px`;
      }
    };

    watch(internalValue, () => {
      nextTick(() => {
        adjustTextareaHeight();
      });
    });

    const handleInput = (event: Event) => {
      const target = event.target as HTMLTextAreaElement;
      internalValue.value = target.value;
      context.emit("update:value", internalValue.value);
    };

    return {
      internalValue,
      textareaRef,
      handleInput,
    };
  },
};
</script>
