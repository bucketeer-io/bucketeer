module.exports = {
  root: true,
  env: { browser: true, es2020: true },
  extends: [
    'eslint:recommended',
    'plugin:@typescript-eslint/recommended',
    'plugin:react-hooks/recommended',
  ],
  rules: {
    'react-refresh/only-export-components': [
      'warn',
      { allowConstantExport: true },
    ],
  },
  ignorePatterns: ['dist', '.eslintrc.cjs'],
  parser: '@typescript-eslint/parser',
  plugins: ['react-refresh'],
  overrides: [
    // Configuration for TypeScript files
    {
      files: ['**/*.ts', '**/*.tsx'],
      plugins: ['@typescript-eslint', 'tailwindcss', 'unused-imports'],
      extends: [
        'airbnb-typescript',
        'plugin:prettier/recommended',
        'plugin:tailwindcss/recommended'
      ],
      parser: '@typescript-eslint/parser',
      parserOptions: {
        project: 'tsconfig.json',
        tsconfigRootDir: __dirname,
        sourceType: 'module'
      },
      rules: {
        'no-plusplus': [2, { allowForLoopAfterthoughts: true }],
        'react/destructuring-assignment': 'off',
        'react/require-default-props': 'off',
        'react/jsx-props-no-spreading': 'off',
        'react-hooks/exhaustive-deps': 'off',
        'react/jsx-filename-extension': 'off',
        '@typescript-eslint/comma-dangle': 'off',
        '@typescript-eslint/no-unused-vars': 'off',
        '@typescript-eslint/naming-convention': [
          'error',
          {
            selector: 'interface',
            format: ['PascalCase']
          }
        ],
        'import/prefer-default-export': 'off',
        'import/extensions': 'off',
        'class-methods-use-this': 'off',
        'import/no-extraneous-dependencies': 'off',
        'import/order': 'off',
        'unused-imports/no-unused-imports': 'warn',
        'unused-imports/no-unused-vars': 'warn',
        'no-nested-ternary': 'off'
      }
    }
  ]
}
