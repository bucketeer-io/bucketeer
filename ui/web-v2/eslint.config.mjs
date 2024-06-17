import eslintConfigPrettier from 'eslint-config-prettier';
import tsParser from '@typescript-eslint/parser';

/** @type {import("eslint").Linter.FlatConfig[]} */
export default [
  {
    files: ['src/**/*.ts', 'src/**/*.tsx'],
    ignores: ['**/proto'],
    languageOptions: {
      parser: tsParser,
      parserOptions: {
        project: './tsconfig.json'
      }
    },
    rules: {}
  },
  eslintConfigPrettier
];
