import { createApp } from "vue";
import App from "./App.vue";
import router from "./config/router";
import setPinia from "./config/setPinia";
import setVueQuery from "./config/setVueQuery";
import { LoadingPlugin } from "vue-loading-overlay";
import VueToast, { ToastProps } from "vue-toast-notification";
import { vue3Debounce } from "vue-debounce";
import FloatingVue from "floating-vue";
import { MotionPlugin } from "@vueuse/motion";
import MasonryWall from "@yeger/vue-masonry-wall";

import "./style.css";
import "vue-toast-notification/dist/theme-bootstrap.css";
import "vue-loading-overlay/dist/css/index.css";
import "floating-vue/dist/style.css";
import "@cyhnkckali/vue3-color-picker/dist/style.css";

const app = createApp(App);

setVueQuery(app);
setPinia(app);

app.directive("debounce", vue3Debounce({ lock: true }));

app
  .use(VueToast, {
    position: "top",
  } as ToastProps)
  .use(FloatingVue)
  .use(MotionPlugin)
  .use(LoadingPlugin)
  .use(router)
  .use(MasonryWall);

router.isReady().then(() => {
  app.mount("#app");
});
