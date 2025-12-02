<script setup lang="ts">
import UserDropdown from "@/shared/components/UserDropdown.vue";

type DesktopNavigationProps = {
  username?: string;
};

defineProps<DesktopNavigationProps>();
</script>

<template>
  <div class="header__actions">
    <template v-if="!username">
      <RouterLink :to="{ name: 'register' }" class="desktop-navlink"
        >Sign In</RouterLink
      >
      <RouterLink :to="{ name: 'login' }" class="desktop-navlink"
        >Log In</RouterLink
      >
    </template>

    <template v-else>
      <RouterLink :to="{ name: 'notes' }" class="desktop-navlink"
        >Notes</RouterLink
      >

      <RouterLink :to="{ name: 'tags' }" class="desktop-navlink">
        Tags
      </RouterLink>

      <UserDropdown />
    </template>
  </div>
</template>

<style scoped>
.router-link-active {
  @apply border-gray-800;
}
</style>
