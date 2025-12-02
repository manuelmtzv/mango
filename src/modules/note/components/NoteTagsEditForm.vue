<script setup lang="ts">
import { computed, ref } from "vue";
import { useToast } from "vue-toast-notification";
import { useTags } from "@/modules/tags/composables/useTags";
import { INote } from "@/modules/note/interfaces/INote";
import TagPillEntry from "@/modules/tags/components/TagPillEntry.vue";
import ListTransitionWrapper from "@/shared/components/animations/ListTransitionWrapper.vue";
import { ITag } from "@/modules/tags/interfaces/ITag";
import ButtonComponent from "@/shared/components/form/ButtonComponent.vue";
import useNoteTagsMutation from "@/modules/note/composables/useNoteTagsMutation";

type NoteTagsEditFormProps = {
  note: INote;
};

type NoteTagsEditFormEmits = {
  (event: "cancel"): void;
  (event: "save"): void;
};

const props = defineProps<NoteTagsEditFormProps>();
const emit = defineEmits<NoteTagsEditFormEmits>();
const { tags, isLoading } = useTags();
const tagsToBeAttached = ref<ITag[]>(
  props.note.tags ? [...props.note.tags] : [],
);
const {
  editNoteTagsAsync,
  editNoteTagsMutation: { isLoading: loadingUpdate },
} = useNoteTagsMutation();
const toast = useToast();

const handleAddTag = (tag: ITag) => {
  if (!tagsToBeAttached.value.some((t) => t.name === tag.name)) {
    tagsToBeAttached.value.push(tag);
  } else {
    tagsToBeAttached.value = tagsToBeAttached.value.filter((t) => t !== tag);
  }
};

const renderTags = computed(() => {
  return tagsToBeAttached.value;
});

const remainingTags = computed(() => {
  return tags.value.filter((tag) => {
    return !tagsToBeAttached.value.some((noteTag) => noteTag.id === tag.id);
  });
});

async function handleSave() {
  const tagIds = tagsToBeAttached.value.map((tag) => tag.id);

  try {
    await editNoteTagsAsync({ id: props.note.id, tags: tagIds });
    emit("save");
    toast.success("Tags saved successfully");
  } catch (error) {
    toast.error("Failed to save tags");
  }
}
</script>

<template>
  <div class="flex flex-col gap-4">
    <section>
      <div>
        <h3 class="text-lg font-semibold mb-2">Note tags</h3>
        <p class="text-sm">Select tags to be detatch from the note</p>
      </div>

      <hr class="my-2" />

      <ListTransitionWrapper class="flex flex-wrap gap-2">
        <TagPillEntry
          v-for="tag in renderTags"
          :key="tag.id"
          :tag="tag.name"
          :color="tag.color"
          :selected="true"
          @click="handleAddTag(tag)"
        />
      </ListTransitionWrapper>

      <span v-if="!renderTags.length" class="text-sm">
        No tags attached to the note.
      </span>
    </section>

    <section>
      <div>
        <h3 class="text-lg font-semibold">Available tags</h3>
        <p class="text-sm">Select tags to be attached to the note</p>
      </div>

      <hr class="my-2" />

      <ListTransitionWrapper class="flex flex-wrap gap-2">
        <template v-if="!isLoading">
          <TagPillEntry
            v-for="tag in remainingTags"
            :key="tag.id"
            :tag="tag.name"
            :color="tag.color"
            @click="handleAddTag(tag)"
          />

          <span v-if="isLoading">Loading tags...</span>
        </template>

        <span v-if="!remainingTags.length && !isLoading" class="text-sm">
          No more tags available.
        </span>
      </ListTransitionWrapper>
    </section>

    <nav class="flex items-center justify-between">
      <ButtonComponent @click="emit('cancel')" class="border-gray-800">
        Cancel
      </ButtonComponent>

      <ButtonComponent
        class="bg-navigation-default hover:bg-navigation-hover"
        :loading="loadingUpdate"
        :disabled="loadingUpdate"
        @click="handleSave"
      >
        Save
      </ButtonComponent>
    </nav>
  </div>
</template>

<style scoped></style>
