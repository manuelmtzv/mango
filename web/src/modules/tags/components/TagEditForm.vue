<script setup lang="ts">
import { reactive, ref, onMounted } from "vue";
import { get } from "@vueuse/core";
import useVuelidate from "@vuelidate/core";
import { required, minLength, maxLength, helpers } from "@vuelidate/validators";
import { Vue3ColorPicker } from "@cyhnkckali/vue3-color-picker";
import { ITag } from "@/modules/tags/interfaces/ITag";
import InputWrapper from "@/shared/components/form/InputWrapper.vue";
import { useToast } from "vue-toast-notification";
import { useTagMutation } from "@/modules/tags/composables/useTagMutation";
import ButtonComponent from "@/shared/components/form/ButtonComponent.vue";
import TagDeleteButton from "@/modules/tags/components/TagDeleteButton.vue";
import TagPreview from "@/modules/tags/components/TagPreview.vue";

type TagEditFormProps = {
  tag: ITag;
};

type TagEditForm = {
  (event: "edited"): void;
};

const props = defineProps<TagEditFormProps>();
const emit = defineEmits<TagEditForm>();
const toast = useToast();
const {
  editTagMutation: { isLoading },
  editTagAsync,
} = useTagMutation();
const showColorPicker = ref(false);

const state = reactive({
  name: props.tag.name,
  color: props.tag.color,
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

function resetState(
  name: string = props.tag.name,
  color: string = props.tag.color,
) {
  state.name = name;
  state.color = color;
}

async function handleSubmit() {
  v$.value.$validate();

  if (v$.value.$invalid) {
    return;
  }

  const result = await editTagAsync({
    id: props.tag.id,
    name: state.name.trim(),
    color: state.color,
  });

  v$.value.$reset();
  resetState(result.name, result.color);
  toast.success("Tag updated successfully!");
  emit("edited");
}

onMounted(() => {
  showColorPicker.value = true;
});
</script>

<template>
  <form
    class="flex flex-col gap-2 p-4 text-sm w-full"
    @submit.prevent="handleSubmit"
  >
    <h2 class="text-2xl font-semibold">Edit tag</h2>

    <InputWrapper class="gap-1" :error="get(v$.name.$errors.at(0)?.$message)">
      <input
        v-model="state.name"
        type="text"
        class="input"
        placeholder="New name"
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

      <span v-if="v$.color.error" class="text-red-500 text-sm mt-1">{{
        v$.color.$errors.at(0)?.$message
      }}</span>
    </fieldset>

    <TagPreview :name="state.name" :color="state.color" />

    <ButtonComponent
      class="button submit !w-full text-base text-gray-800 flex items-center justify-center gap-2"
      :loading="isLoading"
    >
      Save changes
    </ButtonComponent>

    <TagDeleteButton :tag-id="tag.id" @deleted="() => $emit('edited')" />
  </form>
</template>

<style scoped>
.input {
  @apply w-full border border-gray-300;
}
.tag-label {
  @apply font-semibold pb-1 border-b-2 border-gray-200;
}
</style>
