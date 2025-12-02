<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";
import { get } from "@vueuse/core";
import useVuelidate from "@vuelidate/core";
import { required, minLength, maxLength, helpers } from "@vuelidate/validators";
import { Vue3ColorPicker } from "@cyhnkckali/vue3-color-picker";
import { useToast } from "vue-toast-notification";
import InputWrapper from "@/shared/components/form/InputWrapper.vue";
import { useTagMutation } from "@/modules/tags/composables/useTagMutation";
import ButtonComponent from "@/shared/components/form/ButtonComponent.vue";
import { getErrorMessage } from "@/shared/utils/getErrorMessage";
import TagPreview from "@/modules/tags/components/TagPreview.vue";

type TagFormEmits = {
  (event: "created"): void;
};

const state = reactive({
  name: "",
  color: "#fff",
});

const rules = {
  name: {
    required: helpers.withMessage("Tag name is required", required),
    minLength: minLength(2),
    maxLength: maxLength(20),
  },
  color: {
    required: helpers.withMessage("Color is required", required),
  },
};

const v$ = useVuelidate(rules, state);
const toast = useToast();
const emit = defineEmits<TagFormEmits>();
const showColorPicker = ref(false);
const {
  createTagMutation: { isLoading },
  createTagAsync,
} = useTagMutation();

async function handleSubmit() {
  v$.value.$validate();
  if (v$.value.$error) {
    return;
  }

  try {
    await createTagAsync({
      name: state.name.trim(),
      color: state.color,
    });
    toast.success("Tag created successfully!");
    emit("created");
  } catch (err) {
    toast.error(getErrorMessage(err));
  }
}

onMounted(() => {
  showColorPicker.value = true;
});
</script>

<template>
  <form @submit.prevent="handleSubmit" class="form !max-w-none w-full">
    <h2 class="text-2xl font-semibold">Create a new Tag</h2>

    <InputWrapper :error="get(v$.name.$errors.at(0)?.$message)">
      <input
        class="input w-full"
        type="text"
        id="name"
        placeholder="Tag name"
        v-model="state.name"
        maxlength="20"
        :spellcheck="false"
      />
    </InputWrapper>

    <fieldset v-if="showColorPicker">
      <Vue3ColorPicker
        v-model="state.color"
        mode="solid"
        type="HEX"
        :show-alpha="false"
        input-type="RGB"
        class="!w-full !shadow-none !rounded-md input"
      />

      <span v-if="v$.color.$error" class="text-red-500 text-sm mt-1">{{
        v$.color.$errors.at(0)?.$message
      }}</span>
    </fieldset>

    <TagPreview :name="state.name" :color="state.color" />

    <div class="flex items-center gap-4">
      <slot name="left-button" />

      <ButtonComponent
        type="submit"
        class="button submit ml-auto text-gray-800"
        :loading="isLoading"
        >Create tag</ButtonComponent
      >
    </div>
  </form>
</template>

<style scoped></style>
