<script setup lang="ts">
import useVuelidate from "@vuelidate/core";
import { required, minLength, maxLength, helpers } from "@vuelidate/validators";
import { ref, reactive, watch } from "vue";
import { useToast } from "vue-toast-notification";
import TagPillEntry from "./TagPillEntry.vue";

type InlineTagsInputProps = {
  modelValue: string[];
  maxTags?: number;
  notIn?: string[];
};

type InlineTagsInputEmits = {
  (event: "update:modelValue", value: string[]): void;
};

const props = defineProps<InlineTagsInputProps>();
const emit = defineEmits<InlineTagsInputEmits>();
const tags = ref<string[]>([]);
const toast = useToast();

const state = reactive({
  title: "",
});

const rules = {
  title: {
    required: helpers.withMessage("Tag title is required", required),
    minLength: minLength(2),
    maxLength: maxLength(20),
  },
};

const v$ = useVuelidate(rules, state);

function handleSubmit() {
  if (props.maxTags && tags.value.length >= props.maxTags) {
    toast.warning(`You can only add ${props.maxTags} tags in total!`);
    return;
  }

  v$.value.$validate();

  if (v$.value.$invalid) {
    return;
  }

  if (tags.value.includes(state.title) || props.notIn?.includes(state.title)) {
    toast.warning("This tag name already exists!");
    return;
  }

  tags.value.push(state.title);
  resetState();
  v$.value.$reset();
}

function resetState() {
  state.title = "";
}

function removeTag(tag: string) {
  tags.value = tags.value.filter((t) => t !== tag);
}

watch(
  tags,
  (value) => {
    emit("update:modelValue", value);
  },
  {
    deep: true,
  },
);
</script>

<template>
  <form class="flex flex-col gap-2" @submit.prevent="handleSubmit">
    <div v-if="tags.length">
      <h2 class="font-semibold mb-1">New tags:</h2>

      <div class="flex gap-2 flex-wrap">
        <TagPillEntry
          class="cursor-pointer"
          v-for="tag in tags"
          :key="tag"
          :tag="tag"
          :selected="true"
          :tooltip="'Remove tag'"
          @click.prevent="removeTag(tag)"
        />
      </div>
    </div>

    <fieldset class="inline-flex items-center gap-4">
      <label class="label w-full text-sm">
        <input
          v-model="state.title"
          class="input"
          type="text"
          name="title"
          placeholder="New tag..."
        />
      </label>

      <button
        class="button w-fit whitespace-nowrap border-gray-800 border-2 text-sm"
      >
        Add new tag
      </button>
    </fieldset>

    <div v-if="v$.title.$error" class="error">
      {{ v$.title.$errors.at(0)?.$message }}
    </div>
  </form>
</template>

<style lang="css" scoped></style>
