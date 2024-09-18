import eslintConfigPrettier from 'eslint-config-prettier';
import tsParser from '@typescript-eslint/parser';
import tseslint from 'typescript-eslint';

/** @type {import("eslint").Linter.FlatConfig[]} */
export default [
  {
    files: ['src/**/*.ts', 'src/**/*.tsx', 'index.ts']
  },
  {
    ignores: ['src/google/', 'src/protoc-gen-openapiv2/', 'src/proto/']
  },
  ...tseslint.configs.recommended,
  {
    languageOptions: {
      parser: tsParser,
      parserOptions: {
        project: './tsconfig.json',
        tsconfigRootDir: __dirname,
      }
    }
  },
  eslintConfigPrettier
];
