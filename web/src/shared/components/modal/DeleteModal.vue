<script setup lang="ts">
import Modal, { ModalProps } from "@shared/components/modal/Modal.vue";
import ButtonComponent from "@/shared/components/form/ButtonComponent.vue";

interface DeleteModalProps extends ModalProps {
  message: string;
  loading?: boolean;
  onConfirm: () => void;
  closeModal: () => void;
}

const props = withDefaults(defineProps<DeleteModalProps>(), {
  title: "Deleting resource",
  message: "Are you sure you want to delete this resource?",
});
</script>

<template>
  <Modal v-bind="props">
    <template #body>
      <p v-if="$props.message">{{ props.message }}</p>
    </template>

    <template #footer>
      <div class="flex justify-between">
        <ButtonComponent class="border-gray-800" @click="closeModal"
          >Cancel</ButtonComponent
        >

        <ButtonComponent
          class="!bg-delete-default hover:!bg-delete-hover"
          :loading="loading"
          @click="onConfirm"
        >
          Confirm
        </ButtonComponent>
      </div>
    </template>
  </Modal>
</template>

<style scoped></style>
