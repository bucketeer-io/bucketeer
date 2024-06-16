import tsParser from '@typescript-eslint/parser';
import importPlugin from 'eslint-plugin-import';

/** @type {import("eslint").Linter.FlatConfig[]} */
export default [
  {
    files: [
      'src/{proto,google}/**/*.js',
      'src/{proto,google}/**/*.ts',
      'src/{proto,google}/**/*.tsx',
    ],
    languageOptions: {
      parser: tsParser,
      parserOptions: {
        project: './tsconfig.json',
      },
    },
    rules: {
      quotes: ['error', 'single'],
    },
  },
  // {
  //   files: ['src/**/*.ts', 'src/**/*.tsx'],
  //   // ignores: ['**/proto'],
  //   languageOptions: {
  //     parser: tsParser,
  //     parserOptions: {
  //       project: './tsconfig.json',
  //     },
  //   },
  //   plugins: { import: importPlugin },
  //   rules: {
  //     quotes: ['error', 'single', { allowTemplateLiterals: true }],
  //     'import/order': [
  //       'warn',
  //       {
  //         alphabetize: {
  //           order: 'asc',
  //         },
  //         'newlines-between': 'always',
  //       },
  //     ],
  //   },
  // },
];
