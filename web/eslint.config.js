import js from "@eslint/js";
import svelte from "eslint-plugin-svelte";
import prettier from "eslint-config-prettier";

export default [
  js.configs.recommended,
  ...svelte.configs["flat/recommended"],
  prettier,
  {
    ignores: ["build/", ".svelte-kit/", "node_modules/"],
  },
];
