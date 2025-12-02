import { computed } from "vue";

type UseRenderablePagesProps = {
  currentPage: number;
  totalPages: number;
};

export const useRenderablePages = ({
  currentPage,
  totalPages,
}: UseRenderablePagesProps) => {
  return computed(() => {
    const MAX_VISIBLE_PAGES = 5;

    const allPages = Array.from({ length: totalPages }, (_, i) => i + 1);

    if (totalPages <= MAX_VISIBLE_PAGES) {
      return allPages;
    }

    const visiblePages = allPages.slice(
      Math.max(currentPage - (currentPage > 3 ? 3 : 1), 0),
      currentPage + 2,
    );

    const showEllipsisBeforeLast = !visiblePages.includes(
      allPages.slice(-2)[0] || 0,
    );

    const renderablePages = [
      ...visiblePages,
      ...(showEllipsisBeforeLast ? ["..."] : []),
      ...allPages.slice(-2),
    ].filter(Boolean);

    return renderablePages;
  });
};
