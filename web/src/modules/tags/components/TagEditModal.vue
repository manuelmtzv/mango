<script setup lang="ts">
import { watch } from "vue";
import Modal from "@/shared/components/modal/Modal.vue";
import TagEditForm from "@/modules/tags/components/TagEditForm.vue";
import { ITag } from "@/modules/tags/interfaces/ITag";

type TagEditModalProps = {
  isOpen: boolean;
  tag: ITag;
};

type TagEditModalEmits = {
  (event: "close-modal"): void;
  (event: "update:isOpen", value: boolean): void;
};

const props = defineProps<TagEditModalProps>();
const emit = defineEmits<TagEditModalEmits>();

const closeModal = () => {
  emit("update:isOpen", false);
};

watch(
  () => props.isOpen,
  (value) => {
    if (!value) {
      emit("close-modal");
    }
  },
);
</script>

<template>
  <Modal
    :is-open="isOpen"
    :close-modal="closeModal"
    panel-class="w-full max-w-md p-2"
  >
    <template #body>
      <TagEditForm :tag="tag" @edited="closeModal" />
    </template>
  </Modal>
</template>

<style scoped></style>
