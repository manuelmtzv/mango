<script setup lang="ts">
import useNotePagination from "@/modules/note/composables/useNotePagination";
import PaginationArrow from "@shared/components/PaginationArrow.vue";
import PaginationNavigate from "@shared/components/PaginationNavigate.vue";
import useNotes from "@/modules/note/composables/useNotes";
import useNotePaginationState from "../composables/useNotePaginationState";

useNotePagination();

const {
  paginate,
  totalPages,
  paginateRestrictions,
  setPage,
  setPaginateValue,
} = useNotePaginationState();

const { isLoading } = useNotes();
</script>

<template>
  <div v-if="!isLoading && totalPages > 0" class="pagination">
    <nav class="pagination__arrows">
      <PaginationArrow
        button-type="prev"
        @click="setPage('prev')"
        :allowed="paginateRestrictions.canGoPrev"
      />

      <PaginationArrow
        button-type="next"
        @click="setPage('next')"
        :allowed="paginateRestrictions.canGoNext"
      />
    </nav>

    <PaginationNavigate
      class="hidden md:flex"
      :total-pages="totalPages"
      :current-page="paginate.page"
      :set-page="(newPage: number) => setPaginateValue('page', newPage)"
    />

    <div>
      <span> Page {{ paginate.page }} of {{ totalPages }} </span>
    </div>
  </div>
</template>

<style scoped>
.pagination {
  @apply flex justify-between items-center mt-4;
}
.pagination__arrows {
  @apply flex items-center gap-4 select-none;
}
</style>
