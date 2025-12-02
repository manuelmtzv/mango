<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { cn } from "@/shared/utils/cn";
import ButtonComponent from "@shared/components/form/ButtonComponent.vue";
import TransitionWrapper from "@shared/components/TransitionWrapper.vue";
import DeleteModal from "@shared/components/modal/DeleteModal.vue";
import useNoteMutation from "@/modules/note/composables/useNoteMutation";
import { useToast } from "vue-toast-notification";

export interface NoteDeleteButtonProps {
  noteId: string;
  show: boolean;
  containerClass?: string;
  onDeletingChange?: (value: boolean) => void;
  afterDelete?: () => void;
}

const deleting = ref(false);
const props = defineProps<NoteDeleteButtonProps>();
const toast = useToast();
const {
  deleteNoteAsync,
  deleteNoteMutation: { isLoading },
} = useNoteMutation();

const showModal = ref(false);

const handleDelete = async () => {
  try {
    await deleteNoteAsync(props.noteId);
    toast.success("Note deleted successfully");

    if (props.afterDelete) {
      props.afterDelete();
    }
  } catch (error) {
    toast.error("Error deleting note");
  }
};

const openModal = () => {
  showModal.value = true;
  deleting.value = true;
};

const closeModal = () => {
  showModal.value = false;
  deleting.value = false;
};

watch(deleting, (value) => {
  if (props.onDeletingChange) {
    props.onDeletingChange(value);
  }
});

const showButton = computed(() => props.show || deleting.value);
const loadingDelete = computed(() => isLoading.value || deleting.value);
</script>

<template>
  <TransitionWrapper>
    <VTooltip :class="cn(containerClass)">
      <ButtonComponent
        v-if="!$slots['button-content']"
        v-show="showButton"
        :loading="loadingDelete"
        :class="['inline-flex', '!p-1']"
        @click.prevent="openModal"
      >
        <template #default>
          <span v-if="!showModal" class="material-symbols-outlined text-[22px]">
            delete
          </span>
        </template>
      </ButtonComponent>

      <slot
        name="button-content"
        :props="{ loadingDelete, showButton, openModal }"
      />

      <template #popper> Delete note </template>
    </VTooltip>
  </TransitionWrapper>

  <DeleteModal
    :is-open="showModal"
    :close-modal="closeModal"
    :on-confirm="handleDelete"
    :loading="isLoading"
    title="Delete Note"
    message="Are you sure you want to delete this note?"
  />
</template>

<style scoped></style>
