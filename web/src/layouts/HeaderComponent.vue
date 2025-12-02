<script setup lang="ts">
import { computed, ref } from "vue";
import useAuth from "@/modules/auth/composables/useAuth";
import { useWindowSize } from "@vueuse/core";
import MobileSidebar from "@shared/components/MobileSidebar.vue";
import DesktopNavigation from "@/shared/components/DesktopNavigation.vue";

const { username, logout } = useAuth();
const { width } = useWindowSize();

const isMobile = computed(() => width.value < 640);
const mobileSidebar = ref(false);
</script>

<template>
  <header class="header">
    <div class="content header__container">
      <RouterLink :to="{ name: 'home' }" class="header__logo">
        <picture>
          <source srcset="@/assets/imgs/webp/mangocat.webp" type="image/webp" />
          <img
            class="header__icon"
            src="@/assets/imgs/png/mangocat.png"
            alt="MangoCatNotes Logo"
          />
        </picture>
        <h1 class="header__title">Mangocatnotes</h1>
      </RouterLink>

      <DesktopNavigation v-if="!isMobile" :username="username" />

      <button
        class="inline-flex"
        v-else
        @click="mobileSidebar = !mobileSidebar"
      >
        <span
          :class="[
            'material-symbols-outlined text-[26px]',
            mobileSidebar && 'invisible',
          ]"
        >
          menu
        </span>
      </button>

      <MobileSidebar
        :username="username"
        :logout="logout"
        v-model:when="mobileSidebar"
        :class="[width > 640 && 'hidden']"
      />
    </div>
  </header>
</template>

<style scoped lang="css">
.header {
  @apply bg-header-default py-6 shadow-md;
}
.header__container {
  @apply flex gap-4 items-center justify-between;
}
.header__logo {
  @apply flex gap-4 items-center;
}
.header__icon {
  @apply w-[3rem] sm:w-[4rem];
}
.header__title {
  @apply font-bold text-[1.5rem] [line-height:2.25rem] sm:text-3xl;
}
.header__actions {
  @apply hidden gap-4 items-center sm:flex;
}
</style>
