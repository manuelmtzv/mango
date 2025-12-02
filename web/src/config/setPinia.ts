import { App } from "vue";
import { createPinia } from "pinia";

export const pinia = createPinia();

export default function setPinia(app: App): App {
  app.use(pinia);

  return app;
}
