<script setup lang="ts">
import { reactive } from "vue";
import { IRegisterForm } from "@/modules/auth/interfaces";
import useAuth from "@/modules/auth/composables/useAuth";
import LoadingSpin from "@shared/components/LoadingSpin.vue";
import { useRegisterValidation } from "../composables/useRegisterValidation";

const { register, error, isLoading } = useAuth();
const registerForm: IRegisterForm = reactive({
  name: "",
  username: "",
  email: "",
  password: "",
});

const { v$ } = useRegisterValidation(registerForm);

const onSubmit = () => {
  v$.value.$validate();
  if (v$.value.$invalid) return;

  register(registerForm);
};
</script>

<template>
  <div class="py-6 mb-4">
    <form
      class="max-w-md mx-auto grid gap-4 p-4 border rounded-md mt-4"
      @submit.prevent="onSubmit"
    >
      <h2 class="text-2xl font-semibold text-center">Register form</h2>

      <hr />

      <p v-if="error" class="error">{{ error }}</p>

      <label class="label">
        Name:
        <input v-model="registerForm.name" class="input" type="text" />
        <span v-if="v$.name.$error" class="error--small">
          {{ v$.name.$errors.at(0)?.$message }}
        </span>
      </label>

      <label class="label">
        Username:
        <input
          v-model="registerForm.username"
          class="input"
          type="text"
          autocomplete="username"
        />
        <span v-if="v$.username.$error" class="error--small">
          {{ v$.username.$errors.at(0)?.$message }}
        </span>
      </label>

      <label class="label">
        Email:
        <input v-model="registerForm.email" class="input" type="email" />
        <span v-if="v$.email.$error" class="error--small">
          {{ v$.email.$errors.at(0)?.$message }}
        </span>
      </label>

      <label class="label">
        Password:
        <input
          v-model="registerForm.password"
          class="input"
          type="password"
          autocomplete="current-password"
        />
        <span v-if="v$.password.$error" class="error--small">
          {{ v$.password.$errors.at(0)?.$message }}
        </span>
      </label>

      <button
        type="submit"
        class="button submit !w-full inline-flex justify-center gap-2 items-center"
      >
        <span>Submit</span>

        <LoadingSpin :when="isLoading" class="!h-4 !w-4" />
      </button>
    </form>
  </div>
</template>

<style scoped></style>
@/composables/auth/useAuth
