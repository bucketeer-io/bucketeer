import tsEslint from 'typescript-eslint';
import tsParser from '@typescript-eslint/parser';
import eslintPluginPrettierRecommended from 'eslint-plugin-prettier/recommended';

export default [
  {
    files: ['src/**/*.ts', 'test/**/*.ts'],
    ignores: ['**/*.d.ts', 'src/google/', 'src/protoc-gen-openapiv2/', 'src/proto/'],
    languageOptions: {
      parser: tsParser,
      parserOptions: {
        sourceType: 'module',
        project: [`tsconfig.json`, `tsconfig.test.json`],
      },
      globals: {
        node: true,
      },
    },
    plugins: {
      ...tsEslint.configs.recommended,
      eslintPluginPrettierRecommended,
    },
    rules: {
      quotes: ['error', 'single', { avoidEscape: true }],
    },
  },
];
