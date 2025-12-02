<script setup lang="ts">
import NoteList from "@modules/note/components/NoteList.vue";
import NotePagination from "@modules/note/components/NotePagination.vue";
import NoteCreateModal from "@/modules/note/components/NoteCreateModal.vue";
</script>

<template>
  <div class="py-8 mb-4 flex-1 flex flex-col justify-between">
    <NoteList>
      <template #actions>
        <NoteCreateModal />
      </template>

      <template #footer>
        <div class="mt-auto">
          <NotePagination />
        </div>
      </template>
    </NoteList>
  </div>
</template>
