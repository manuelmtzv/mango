<script setup lang="ts">
import { computed, ref } from "vue";
import { useNoteStore } from "@/modules/note/stores/noteStore";
import useNoteMutation from "@/modules/note/composables/useNoteMutation";
import { useToast } from "vue-toast-notification";
import LoadingSpin from "@shared/components/LoadingSpin.vue";
import TagsSelect from "@/modules/tags/components/TagsSelect.vue";
import { ITag } from "@/modules/tags/interfaces/ITag";
import InlineTagsInput from "@/modules/tags/components/InlineTagsInput.vue";
import { useTags } from "@/modules/tags/composables/useTags";
import { union } from "ramda";
import { MAX_TAGS_PER_USER } from "@/config/constants";

export interface NoteFormProps {
  setModal?: (value: boolean) => void;
}

const props = defineProps<NoteFormProps>();
const noteStore = useNoteStore();
const { createNoteAsync, createNoteMutation } = useNoteMutation();
const title = ref<string>("");
const content = ref<string>("");
const tagsToAdd = ref<ITag[]>([]);
const tagsToCreate = ref<string[]>([]);
const contentIsEmpty = ref<boolean>(false);
const openTagForm = ref<boolean>(false);
const useTagsState = useTags();

async function handleSubmit() {
  if (content.value != "") {
    try {
      const newNote = await createNoteAsync({
        title: title.value,
        content: content.value,
        tags: union(
          tagsToAdd.value.map((tag) => tag.name),
          tagsToCreate.value,
        ),
      });
      noteStore.addNote(newNote);
      resetValues();
      useToast().success("Note created successfully!");

      if (props.setModal) {
        props.setModal(false);
      }
    } catch (err) {
      useToast().error("Something went wrong! Please try again.");
    }
  } else {
    contentIsEmpty.value = true;
  }
}

function resetValues(): void {
  title.value = content.value = "";
  contentIsEmpty.value = false;
}

const tagNames = computed(() => {
  return useTagsState.tags.value.map((tag) => tag.name);
});

const maxAllowedTags = computed(() => {
  return MAX_TAGS_PER_USER - tagNames.value.length;
});
</script>

<template>
  <form class="form" @submit.prevent="handleSubmit">
    <h2 class="form-title">Create new Note</h2>

    <label class="label">
      <span>Title</span>
      <input
        class="input"
        type="text"
        name="title"
        v-model="title"
        placeholder="Learnings..."
      />
    </label>

    <label class="label">
      <p class="font-semibold">
        Content <span v-tooltip="'Required'" class="text-red-500">*</span>
      </p>
      <textarea
        class="input"
        type="text"
        name="content"
        v-model="content"
        placeholder="Today I learned..."
      ></textarea>
      <span v-if="contentIsEmpty" class="error"
        >The content field is required!</span
      >
    </label>

    <TagsSelect v-model="tagsToAdd" :use-tags-state="useTagsState">
      <template #title>
        <div class="inline-flex items-center justify-between gap-4 mb-2">
          <div>
            <h2 class="font-semibold">
              {{
                `Tags (${useTagsState.tags.value.length}/${MAX_TAGS_PER_USER}):`
              }}
            </h2>
            <p class="text-sm">Select tags to attach to the new note.</p>
          </div>

          <button
            v-tooltip="
              maxAllowedTags <= 0
                ? 'You have reached the maximum number of tags allowed.'
                : ''
            "
            :disabled="maxAllowedTags <= 0 || useTagsState.isLoading.value"
            :class="[
              'py-0.5 px-2.5 rounded-md text-gray-800 border border-gray-800 relative cursor-pointer select-none text-sm transition-all duration-300',
              maxAllowedTags <= 0 && 'opacity-50',
            ]"
            @click.prevent="openTagForm = !openTagForm"
          >
            {{ openTagForm ? "Close" : "Add" }}
          </button>
        </div>
      </template>

      <template #default>
        <div v-if="openTagForm" class="pt-4 border-t">
          <InlineTagsInput
            v-model="tagsToCreate"
            :max-tags="maxAllowedTags"
            :not-in="tagNames"
          />
        </div>
      </template>
    </TagsSelect>

    <nav class="flex gap-4 justify-between">
      <button
        v-if="props.setModal"
        class="button cancel flex items-center gap-2 border-gray-800"
        type="button"
        @click="props.setModal(false)"
      >
        <span>Cancel</span>
      </button>

      <button class="button submit flex items-center gap-2" type="submit">
        <span>Save</span>

        <LoadingSpin
          :when="createNoteMutation.isLoading.value"
          class="!h-4 !w-4"
        />
      </button>
    </nav>
  </form>
</template>

<style scoped lang="css"></style>
