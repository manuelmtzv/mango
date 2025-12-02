import { App } from "vue";
import { VueQueryPlugin } from "@tanstack/vue-query";

export default function setVueQuery(app: App): App {
  VueQueryPlugin.install(app, {
    queryClientConfig: {
      defaultOptions: {
        queries: {
          cacheTime: 1000 * 60,
        },
      },
    },
  });

  return app;
}
