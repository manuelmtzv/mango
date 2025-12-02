import { RouteRecordRaw } from "vue-router";

export const noteRoutes: RouteRecordRaw[] = [
  {
    path: "notes",
    name: "notes",
    component: () => import("@/modules/note/views/Index.vue"),
  },
  {
    path: "notes/:id",
    name: "note-detail",
    component: () => import("@/modules/note/views/NoteView.vue"),
  },
];
