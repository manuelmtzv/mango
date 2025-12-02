import { RouteRecordRaw } from "vue-router";

export const publicRoutes = [
  {
    path: "/welcome",
    name: "welcome",
    component: () => import("@/shared/views/Welcome.vue"),
  },
] as RouteRecordRaw[];
