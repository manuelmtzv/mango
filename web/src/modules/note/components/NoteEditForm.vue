<script setup lang="ts">
import { ref, onMounted, watch } from "vue";
import { INote } from "@/modules/note/interfaces/INote";
import { useRouter } from "vue-router";
import useNoteMutation from "@/modules/note/composables/useNoteMutation";
import useInvalidateQuery from "@/shared/composables/useInvalidateQuery";
import NoteFooterDetails from "./NoteFooterDetails.vue";
import EditedFeedback from "@shared/components/EditedFeedback.vue";

interface IProps {
  note: INote;
}

interface Emits {
  (event: "refetch-note"): Promise<void>;
}

const props = defineProps<IProps>();
const emits = defineEmits<Emits>();
const noteId = props.note.id;
const router = useRouter();
const { editNoteAsync } = useNoteMutation();
const { invalidateQuery } = useInvalidateQuery();

const title = ref<string>(props.note.title || "");
const content = ref<string>(props.note.content);
const contentIsEmpty = ref<boolean>(props.note.content == "");
const editMode = ref<boolean>(false);
const editedFeedbackRef = ref<typeof EditedFeedback>();

const date = ref<string>("");
const time = ref<string>("");

async function handleEdit() {
  if (content.value != "") {
    await editNoteAsync({
      id: noteId,
      title: title.value,
      content: content.value,
    });
    resetValues();
    invalidateQuery(["note", noteId]);
    editedFeedbackRef.value?.resetTimer();
  } else {
    contentIsEmpty.value = true;
  }
  emits("refetch-note");
}

function setEditMode(value: boolean) {
  editMode.value = value;
}

function resetValues() {
  setEditMode(false);
  contentIsEmpty.value = false;
}

onMounted(() => {
  if (props.note.content) {
    title.value = props.note.title || "";
    content.value = props.note.content;
  } else {
    router.push("/404");
  }
});

watch(
  () => props.note,
  (value) => {
    if (value) {
      date.value = new Date(props.note.updatedAt).toLocaleDateString();
      time.value = new Date(props.note.updatedAt).toLocaleTimeString();
    }
  },
  {
    immediate: true,
  },
);
</script>

<template>
  <section class="flex flex-col gap-4 flex-1">
    <form class="form w-full !max-w-none" @submit.prevent="null">
      <label class="label">
        <span>Title:</span>
        <input
          name="title"
          class="input"
          type="text"
          v-model="title"
          placeholder="Learnings..."
          @input.prevent="() => setEditMode(true)"
          v-debounce:500ms="handleEdit"
        />
      </label>
      <label class="label">
        <span>Content:</span>
        <textarea
          class="input !font-normal !min-h-[22rem]"
          name="content"
          placeholder="Today I learned that..."
          v-model="content"
          ref="contentTextareaRef"
          @input.prevent="() => setEditMode(true)"
          v-debounce:500ms="handleEdit"
        ></textarea>
      </label>

      <span v-show="contentIsEmpty" class="error"
        >The content field is required!</span
      >
    </form>

    <slot name="after-form" />

    <NoteFooterDetails :date="date" :time="time" class="mt-auto">
      <template #default>
        <EditedFeedback ref="editedFeedbackRef" />
      </template>
    </NoteFooterDetails>
  </section>
</template>

<style scoped lang="css">
.actions {
  @apply flex items-center justify-between;
}
.actions__notes {
  @apply flex items-center gap-2;
}
.actions__edit {
  @apply bg-edit-default hover:bg-edit-hover;
}
.disabled {
  @apply text-gray-400 !bg-transparent hover:!bg-transparent;
}
</style>
../stores/noteStore
