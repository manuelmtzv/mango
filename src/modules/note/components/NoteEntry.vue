<script setup lang="ts">
import { computed, ref } from "vue";
import { cutString } from "@/shared/utils/cutString";
import type { INote } from "@/modules/note/interfaces/INote";
import NoteDeleteButton from "@/modules/note/components/NoteDeleteButton.vue";
import TagPillEntry from "@/modules/tags/components/TagPillEntry.vue";
import noteDefaults from "@/modules/note/config/defaults";
import ListTransitionWrapper from "@/shared/components/animations/ListTransitionWrapper.vue";
import { ITag } from "@/modules/tags/interfaces/ITag";

interface IProps {
  note: INote;
}

const props = defineProps<IProps>();

const deletingNote = ref<boolean>(false);
const showDeleteButton = ref<boolean>(false);
const tags = ref<ITag[]>(props.note?.tags ?? []);

const renderTags = computed(() =>
  tags.value.slice(0, noteDefaults.NOTE_ENTRY_MAX_TAGS),
);

const isOverflow = computed(
  () => tags.value.length > noteDefaults.NOTE_ENTRY_MAX_TAGS,
);
</script>

<template>
  <router-link
    :class="['note-entry', deletingNote && 'note-entry--deleting']"
    :to="{ name: 'note-detail', params: { id: note.id } }"
    :key="note.id"
    @mouseover="showDeleteButton = true"
    @mouseleave="showDeleteButton = false"
  >
    <div>
      <h3 class="note-entry__title" v-if="note.title">
        {{ cutString(note.title, 50) }}
      </h3>
      <p class="note-entry__description">{{ cutString(note.content, 75) }}</p>
    </div>

    <template v-if="renderTags.length">
      <ListTransitionWrapper
        :id="note.id"
        class="flex gap-2 items-center flex-wrap mt-auto"
      >
        <TagPillEntry
          v-for="tag in renderTags"
          :key="tag.id"
          :tag="tag.name"
          :color="tag.color"
        />

        <span v-if="isOverflow" class="text-xs text-gray-500">
          +{{ tags.length - noteDefaults.NOTE_ENTRY_MAX_TAGS }}
        </span>
      </ListTransitionWrapper>
    </template>

    <NoteDeleteButton
      container-class="delete-button hidden sm:inline-flex absolute top-2 right-2 bg-white rounded-md"
      :note-id="note.id"
      :show="showDeleteButton"
      :on-deleting-change="
        (deleting: boolean) => {
          deletingNote = deleting;
          showDeleteButton = deleting;
        }
      "
    />
  </router-link>
</template>

<style lang="css">
.note-entry {
  @apply flex flex-col gap-3 p-4 rounded-md border bg-entry-default hover:bg-gray-100 transition-colors duration-300 relative shadow-sm;
}
.note-entry__title {
  @apply font-semibold overflow-hidden;
}
.note-entry__description {
  @apply overflow-hidden;
}
</style>
