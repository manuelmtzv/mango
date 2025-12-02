<script setup lang="ts">
import { useTagMutation } from "@/modules/tags/composables/useTagMutation";
import { getErrorMessage } from "@/shared/utils/getErrorMessage";
import { useToast } from "vue-toast-notification";
import DeleteModal from "@shared/components/modal/DeleteModal.vue";
import ButtonComponent from "@/shared/components/form/ButtonComponent.vue";
import { ref } from "vue";

type TagDeleteButtonprops = {
  tagId: string;
};

type TagDeleteButtonEmits = {
  (event: "deleted"): void;
};

const props = defineProps<TagDeleteButtonprops>();
const emit = defineEmits<TagDeleteButtonEmits>();
const {
  deleteTagMutation: { isLoading },
  deleteTagAsync,
} = useTagMutation();
const toast = useToast();
const showModal = ref(false);

const setDeleteModal = (value: boolean) => {
  showModal.value = value;
};

async function handleDelete() {
  try {
    await deleteTagAsync(props.tagId);

    toast.success("Tag deleted successfully");
    emit("deleted");
  } catch (err) {
    toast.error(getErrorMessage(err));
  }
}
</script>

<template>
  <ButtonComponent
    @click.prevent="setDeleteModal(true)"
    class="button !w-full !bg-delete-default hover:!bg-delete-hover !text-base justify-center text-gray-800"
    :loading="isLoading"
  >
    <span> Delete tag </span>
  </ButtonComponent>

  <DeleteModal
    :is-open="showModal"
    :close-modal="() => setDeleteModal(false)"
    :on-confirm="handleDelete"
    title="Delete Tag"
    message="Are you sure you want to delete this tag?"
  />
</template>

<style scoped></style>
