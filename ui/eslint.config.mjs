// @ts-check

import eslint from '@eslint/js';
import { defineConfig, globalIgnores } from 'eslint/config';
import tseslint from 'typescript-eslint';
import reactPlugin from 'eslint-plugin-react';

export default defineConfig(
  globalIgnores([
    "build",
    ".react-router",
    "node_modules",
  ]),
  eslint.configs.recommended,
  tseslint.configs.recommended,
  reactPlugin.configs.flat.recommended,
  reactPlugin.configs.flat['jsx-runtime'],
);
