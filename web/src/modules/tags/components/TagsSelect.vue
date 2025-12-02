<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { useTags } from "@/modules/tags/composables/useTags";
import { ITag } from "@/modules/tags/interfaces/ITag";
import TagPillEntry from "./TagPillEntry.vue";

const MAX_RENDERING = 30;

type TagsSelectProps = {
  modelValue: ITag[];
  useTagsState: ReturnType<typeof useTags>;
};

type TagsSelectEmits = {
  (event: "update:modelValue", value: ITag[]): void;
};

defineProps<TagsSelectProps>();
const emit = defineEmits<TagsSelectEmits>();
const { tags, isLoading } = useTags();
const selectedTags = ref<ITag[]>([]);
const limitRendering = ref(true);

const handleTagTap = (tag: ITag) => {
  const index = selectedTags.value.findIndex((t) => t.id === tag.id);
  if (index === -1) {
    selectedTags.value.push(tag);
  } else {
    selectedTags.value.splice(index, 1);
  }
};

const checkSelected = (tag: ITag) => {
  return selectedTags.value.some((t) => t.id === tag.id);
};

const renderingTags = computed(() => {
  return limitRendering.value ? tags.value.slice(0, MAX_RENDERING) : tags.value;
});

watch(
  selectedTags,
  (value) => {
    emit("update:modelValue", value);
  },
  {
    deep: true,
  },
);
</script>

<template>
  <div class="flex flex-col gap-2">
    <slot v-if="$slots['title']" name="title" />

    <h2 v-else class="font-semibold">Tags:</h2>

    <div class="flex flex-col gap-4 p-4 border rounded-md">
      <div
        :v-if="!isLoading"
        class="v-collapse flex items-center gap-2 flex-wrap justify-around"
      >
        <template v-if="renderingTags.length">
          <TagPillEntry
            :class-name="tag.color && `bg-[${tag.color}]`"
            v-for="tag in renderingTags"
            :key="tag.id"
            :tag="tag.name"
            :selected="checkSelected(tag)"
            :color="tag.color"
            @click="handleTagTap(tag)"
          />

          <TagPillEntry
            v-if="tags.length > MAX_RENDERING"
            @click.prevent="limitRendering = !limitRendering"
            class="text-sm text-white bg-gray-800 cursor-pointer text-center hover:bg-gray-700"
            :tag="limitRendering ? 'Show more tags' : 'Show less'"
          />
        </template>

        <p v-else class="text-sm text-gray-800">No tags to show.</p>
      </div>

      <p v-if="isLoading" class="text-sm text-gray-800">Loading...</p>

      <slot />
    </div>
  </div>
</template>

<style scoped></style>
