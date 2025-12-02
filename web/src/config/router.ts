import {
  createRouter,
  createWebHistory,
  NavigationGuardNext,
  RouteLocationNormalized,
  RouteLocationRaw,
  RouteRecordRaw,
} from "vue-router";
import { authRoutes } from "@/modules/auth/routes/auth.routes";
import { appRoutes } from "@/shared/routes/app.routes";
import { publicRoutes } from "@/shared/routes/public.routes";
import { useTokenValidation } from "@/modules/auth/composables/useTokenValidation";
import { IAuthResponse } from "@/modules/auth/interfaces";
import { useLoading } from "vue-loading-overlay";

const routes = [
  {
    path: "",
    component: () => import("@/layouts/Layout.vue"),
    meta: {
      authRequired: false,
    },
    children: [...publicRoutes, ...authRoutes],
  },

  {
    path: "/app",
    name: "home",
    component: () => import("@/layouts/Layout.vue"),
    redirect: { name: "notes" },
    meta: {
      authRequired: true,
    },
    children: appRoutes,
  },

  {
    path: "/:pathMatch(.*)*",
    name: "not-found",
    component: () => import("@/layouts/Layout.vue"),
    children: [
      {
        path: "",
        name: "not-found",
        component: () => import("@/shared/views/NotFound.vue"),
      },
    ],
  },
] as RouteRecordRaw[];

const router = createRouter({
  routes,
  history: createWebHistory(),
});

const handleRootNavigation = (
  to: RouteLocationNormalized,
  next: NavigationGuardNext,
  redirect: RouteLocationRaw,
) => {
  to.path === "/" ? next(redirect) : next();
};

router.beforeEach(async (to, from, next) => {
  const loading = useLoading().show();
  const { validateToken, jwt } = useTokenValidation();
  const authRequired = to.matched.some((route) => route.meta.authRequired);
  let authenticated: void | IAuthResponse;

  if (authRequired) {
    authenticated = await validateToken();
    loading.hide();

    authenticated
      ? handleRootNavigation(to, next, { name: "notes" })
      : next({ name: "welcome" });
  } else {
    jwt.value
      ? (authenticated = await validateToken())
      : (authenticated = undefined);
    loading.hide();

    authenticated
      ? next({ name: from?.name || "notes" })
      : handleRootNavigation(to, next, { name: "welcome" });
  }
});

export default router;
