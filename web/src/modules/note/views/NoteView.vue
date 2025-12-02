<script setup lang="ts">
import { useRoute } from "vue-router";
import { useRouter } from "vue-router";
import Loading from "vue-loading-overlay";
import NoteEditForm from "@modules/note/components/NoteEditForm.vue";
import NotFound from "@shared/views/NotFound.vue";
import useNote from "@/modules/note/composables/useNote";
import NoteTags from "@/modules/note/components/NoteTags.vue";
import NoteDeleteButton from "@/modules/note/components/NoteDeleteButton.vue";
import NoteTagsEditModal from "@/modules/note/components/NoteTagsEditModal.vue";
import NoteLinks from "@/modules/note/components/NoteLinks.vue";

const { id: noteId } = useRoute().params as { id: string };
const router = useRouter();

const {
  data: note,
  isLoading,
  refetch,
} = useNote({
  id: noteId,
});

async function afterDelete() {
  router.push({ name: "home" });
}
</script>

<template>
  <div v-if="!isLoading && note" class="flex flex-col py-6 mb-4 flex-1">
    <div class="flex items-center justify-between gap-4 mb-4">
      <h2 class="text-xl font-semibold">Note View</h2>
      <NoteDeleteButton
        :note-id="note.id"
        :show="true"
        :after-delete="afterDelete"
      >
        <template #button-content="{ props: slotProps }">
          <button
            class="button bg-delete-default hover:bg-delete-hover"
            @click.prevent="slotProps.openModal"
          >
            Delete
          </button>
        </template>
      </NoteDeleteButton>
    </div>

    <NoteEditForm :note="note" @refetch-note="refetch">
      <template #after-form>
        <NoteTags :tags="note.tags">
          <template #header>
            <div class="flex gap-4 justify-between items-center">
              <h3 class="font-semibold">Tags:</h3>

              <NoteTagsEditModal :note="note" />
            </div>

            <hr />
          </template>
        </NoteTags>

        <NoteLinks :content="note.content" />
      </template>
    </NoteEditForm>
  </div>

  <Loading v-model:active="isLoading" />

  <NotFound v-if="!isLoading && !note" />
</template>
