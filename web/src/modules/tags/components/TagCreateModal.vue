<script setup lang="ts">
import { ref } from "vue";
import Modal from "@/shared/components/modal/Modal.vue";
import TagForm from "@/modules/tags/components/TagForm.vue";
import ButtonComponent from "@/shared/components/form/ButtonComponent.vue";

type TagCreateModalProps = {
  disabled?: boolean;
};

defineProps<TagCreateModalProps>();
const isOpen = ref(false);

const openModal = () => {
  isOpen.value = true;
};

const closeModal = () => {
  isOpen.value = false;
};
</script>

<template>
  <div>
    <ButtonComponent class="submit" @click="openModal" :disabled="disabled">
      <span class="hidden sm:block">Add new tag</span>
      <span class="material-symbols-outlined block sm:hidden text-[22px]">
        add_circle
      </span>
    </ButtonComponent>

    <Modal
      v-bind="$props"
      :is-open="isOpen"
      :close-modal="closeModal"
      panel-class="max-w-xl"
    >
      <template #body>
        <TagForm @created="closeModal()">
          <template #left-button>
            <ButtonComponent
              class="!bg-white border-gray-800"
              @click.prevent="closeModal"
            >
              Cancel
            </ButtonComponent>
          </template>
        </TagForm>
      </template>
    </Modal>
  </div>
</template>

<style scoped></style>
