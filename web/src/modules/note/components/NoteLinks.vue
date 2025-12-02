<script setup lang="ts">
import { computed } from "vue";
import { LinkIcon } from "@heroicons/vue/24/solid";
import ListTransitionWrapper from "@/shared/components/animations/ListTransitionWrapper.vue";

type NoteLinksProps = {
  content: string;
};

const props = defineProps<NoteLinksProps>();

const links = computed(() => props.content.match(/(https?:\/\/[^\s]+)/g) || []);

const DISCLAIMER = `
  Please note that the following links are provided for your convenience, and we do not have control over the content they lead to.
`;
</script>

<template>
  <section v-if="links" class="flex flex-col gap-2">
    <div class="flex justify-between gap-2">
      <h3 class="font-semibold">Links:</h3>
      <span
        class="material-symbols-outlined text-gray-600"
        v-tooltip="DISCLAIMER"
      >
        info
      </span>
    </div>

    <hr />

    <ListTransitionWrapper class="flex flex-wrap gap-2 mt-1">
      <a
        v-for="link in links"
        :key="link"
        :href="link"
        target="_blank"
        class="flex items-center gap-2 text-sm px-3 py-1 rounded-full text-white bg-gray-700 relative cursor-pointer select-none transition-opacity hover:opacity-80 ease-in-out"
        >{{ link }} <LinkIcon class="w-4"
      /></a>
    </ListTransitionWrapper>
  </section>
</template>
