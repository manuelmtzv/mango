import { RouteRecordRaw } from "vue-router";

export const tagRoutes: RouteRecordRaw[] = [
  {
    path: "tags",
    name: "tags",
    component: () => import("@/modules/tags/views/Index.vue"),
  },
];
