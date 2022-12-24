const webpack = require('webpack');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const OptimizeCSSAssetsPlugin = require('optimize-css-assets-webpack-plugin');
const nrwlConfig = require('@nrwl/react/plugins/webpack.js');

module.exports = (config, context) => {
  nrwlConfig(config);

  config.module.rules.splice(1, 1);

  return {
    ...config,
    module: {
      ...config.module,
      rules: [
        ...config.module.rules,
        {
          test: /\.m?js$/,
          type: 'javascript/auto',
        },
        {
          test: /\.css$/,
          use: [
            MiniCssExtractPlugin.loader,
            {
              loader: 'css-loader',
            },
            {
              loader: 'postcss-loader',
              options: {
                postcssOptions: {
                  // This is necessary to import css from node_modules
                  config: './apps/admin/src/postcss.config.js',
                },
              },
            },
          ],
        },
      ],
    },
    plugins: [
      ...config.plugins,
      new MiniCssExtractPlugin(),
      new webpack.EnvironmentPlugin(['RELEASE_CHANNEL']),
    ],
    optimization: {
      ...config.optimization,
      minimizer: [
        ...config.optimization.minimizer,
        new OptimizeCSSAssetsPlugin({}),
      ],
      splitChunks: {
        ...config.optimization.splitChunks,
        cacheGroups: {
          ...config.optimization.splitChunks.cacheGroups,
          vendor: {
            name: 'vendor',
            chunks: 'all',
            test: (module) => {
              if (module.resource && /\.css$/u.test(module.resource)) {
                return false;
              }
              return module.context && module.context.includes('node_modules');
            },
            enforce: true,
          },
        },
      },
    },
    stats: {
      assets: false,
      entrypoints: false,
      children: false,
      modules: false,
      colors: true,
    },
    devServer: {
      ...config.devServer,
      host: 'localhost',
      port: 8000,
      disableHostCheck: true,
    },
  };
};
