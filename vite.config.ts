import { fileURLToPath, URL } from "node:url";

import { defineConfig } from "vite";
import Components from "@uni-helper/vite-plugin-uni-components";
import { WotResolver } from "./wot-ui-resolver";
import Uni from "@uni-helper/plugin-uni";

export default defineConfig({
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
  plugins: [
    Components({
      dts: true,
      resolvers: [WotResolver()],
    }),
    Uni(),
  ],
  css: {
    preprocessorOptions: {
      scss: {
        silenceDeprecations: ["legacy-js-api"],
      },
    },
  },
});
