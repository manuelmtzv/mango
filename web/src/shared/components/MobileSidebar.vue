<script setup lang="ts">
import { useMotions } from "@vueuse/motion";

interface Props {
  when: boolean;
  username?: string;
  logout?: () => void;
}

interface Emits {
  (event: "update:when", value: boolean): boolean;
}

const props = defineProps<Props>();
const emit = defineEmits<Emits>();
const motions = useMotions();

function setWhen(value: boolean) {
  setTimeout(() => emit("update:when", value), 50);
}

const handleLogout = () => {
  if (props.logout) props.logout();
  setWhen(false);
};
</script>

<template>
  <transition @leave="(_, done) => motions.transparent.leave(done)">
    <div
      v-if="when"
      v-motion="'transparent'"
      :initial="{
        opacity: 0,
        x: -640,
      }"
      :enter="{
        opacity: 100,
        x: 0,
        transition: {
          type: 'spring',
          stiffness: 400,
          damping: 40,
          mass: 0.5,
        },
      }"
      :leave="{
        opacity: 0,
        x: -640,
        transition: {
          type: 'spring',
          stiffness: 80,
          mass: 0.8,
        },
      }"
      class="fixed w-full inset-0 bg-black bg-opacity-80 z-50 text-white backdrop-blur-sm"
    >
      <section class="flex flex-col gap-4 w-[90%] mx-auto min-h-[100svh] py-6">
        <div class="flex items-center my-[2.25rem] mt-[1rem] gap-4">
          <h4 v-if="username" class="text-xl">{{ `Hi, ${username}` }}</h4>

          <button
            class="material-symbols-outlined ml-auto text-[26px]"
            @click="setWhen(false)"
          >
            close
          </button>
        </div>

        <nav class="flex flex-col flex-1 h-full gap-6 w-full">
          <template v-if="!username">
            <RouterLink
              :to="{ name: 'notes' }"
              class="button link"
              @click.prevent="setWhen(false)"
            >
              Home
            </RouterLink>

            <hr />

            <RouterLink
              :to="{ name: 'register' }"
              class="button link"
              @click.prevent="setWhen(false)"
            >
              Sign In
            </RouterLink>

            <RouterLink
              :to="{ name: 'login' }"
              class="button link"
              @click.prevent="setWhen(false)"
            >
              Log In
            </RouterLink>
          </template>

          <template v-else>
            <RouterLink
              :to="{ name: 'notes' }"
              class="button link"
              @click.prevent="setWhen(false)"
            >
              Notes
            </RouterLink>

            <RouterLink
              :to="{ name: 'tags' }"
              class="button link"
              @click.prevent="setWhen(false)"
            >
              Tags
            </RouterLink>

            <button
              class="button link link--logout bg-red-500 hover:bg-red-600 !text-center mt-auto"
              @click="handleLogout"
            >
              Logout
            </button>
          </template>
        </nav>
      </section>
    </div>
  </transition>
</template>

<style scoped>
.link {
  @apply block !w-full text-left;
}
</style>
