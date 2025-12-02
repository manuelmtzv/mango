import { RouteRecordRaw } from "vue-router";

export const authRoutes = [
  {
    path: "/login",
    name: "login",
    component: () => import("@/modules/auth/views/Login.vue"),
  },
  {
    path: "/register",
    name: "register",
    component: () => import("@/modules/auth/views/Register.vue"),
  },
] as RouteRecordRaw[];
