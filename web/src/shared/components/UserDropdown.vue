<script setup lang="ts">
import useAuth from "@/modules/auth/composables/useAuth";

const { username, logout } = useAuth();
</script>

<template>
  <VDropdown>
    <button class="select-none">{{ username }}</button>

    <template #popper>
      <div class="w-full max-w-[16rem] p-2">
        <button
          class="px-5 py-0.5 border rounded-md hover:bg-red-300"
          @click="logout"
        >
          Log out
        </button>
      </div>
    </template>
  </VDropdown>
</template>

<style scoped></style>
