// eslint.config.js
import { defineConfig } from "eslint/config";
import globals from "globals";
import js from "@eslint/js";
import tseslint from "@typescript-eslint/eslint-plugin";
import tsParser from "@typescript-eslint/parser";
import vuePlugin from "eslint-plugin-vue";
import vueParser from "vue-eslint-parser";

const vueEssential = vuePlugin.configs["flat/essential"];

export default defineConfig([
    // 1. Official base recommended rules for JavaScript
    js.configs.recommended,

    // 2. Global universal formatting rules (applied to all files, maintained in one place)
    {
        rules: {
            "no-multi-spaces": [
                "error",
                {
                    ignoreEOLComments: false,
                    exceptions: {
                        Property: false,
                    },
                },
            ],
            indent: ["error", 4],
            "no-trailing-spaces": ["error"],
            "space-before-blocks": ["error", "always"],
            "space-before-function-paren": [
                "error",
                {
                    named: "never",
                    anonymous: "always",
                    asyncArrow: "always",
                },
            ],
            "space-in-parens": ["error", "never"],
            "space-infix-ops": ["error"],
            "spaced-comment": ["error", "always"],
            "eol-last": ["error", "always"],
            "no-multiple-empty-lines": ["error", { max: 1, maxEOF: 0 }],
            // Enforce Unix LF line breaks, disallow Windows CRLF
            "linebreak-style": ["error", "unix"],
        },
    },

    // 3. Base Vue lint rules
    ...(Array.isArray(vueEssential) ? vueEssential : [vueEssential]),

    // 4. TypeScript / TSX exclusive config
    // General formatting rules are already applied globally, no duplicate definition needed
    {
        files: ["**/*.ts", "**/*.tsx"],
        plugins: {
            "@typescript-eslint": tseslint,
        },
        languageOptions: {
            parser: tsParser,
            parserOptions: {
                ecmaVersion: "latest",
                sourceType: "module",
            },
        },
        rules: {
            ...tseslint.configs.recommended.rules,
            "@typescript-eslint/no-explicit-any": "warn",
            "@typescript-eslint/no-unused-vars": ["error"],
            "no-unused-vars": "off",
        },
    },

    // 5. Vue Single-File Component exclusive config
    {
        files: ["**/*.vue"],
        languageOptions: {
            parser: vueParser,
            parserOptions: {
                ecmaVersion: "latest",
                sourceType: "module",
                parser: tsParser,
                extraFileExtensions: [".vue"],
            },
            globals: {
                ...globals.browser,
                ...globals.es2021,
            },
        },
        plugins: {
            vue: vuePlugin,
        },
        rules: {
            "vue/multi-word-component-names": "off",
        },
    },

    // 6. General JS / CJS / MJS environment config
    // Includes vite.config.js and all other project configuration files
    {
        files: ["**/*.js", "**/*.cjs", "**/*.mjs"],
        languageOptions: {
            globals: {
                ...globals.browser,
                ...globals.es2021,
                ...globals.node,
            },
        },
    },

    // Exclude directories and files from linting
    {
        ignores: ["node_modules", "dist", "build", ".vscode", ".idea", "*.min.js", "*.d.ts", "coverage", ".git"],
    },
]);
