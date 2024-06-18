const path = require('path');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');

module.exports = {
  entry: './src/main.tsx',
  output: {
    filename: '[name].[contenthash].bundle.js',
    library: 'BucketeerAdmin',
    globalObject: 'this',
    libraryTarget: 'umd',
    path: path.resolve(__dirname, 'dist'),
    clean: true
  },
  module: {
    rules: [
      {
        test: /\.(png|jpg|gif|ico)$/,
        type: 'asset/resource',
        generator: {
          filename: 'assets/[name][ext]'
        }
      },
      {
        test: /\.svg$/i,
        issuer: /\.[jt]sx?$/,
        use: ['@svgr/webpack']
      },
      {
        test: /\.css$/,
        use: [
          MiniCssExtractPlugin.loader,
          {
            loader: 'css-loader'
          },
          {
            loader: 'postcss-loader',
            options: {
              postcssOptions: {
                // This is necessary to import css from node_modules
                config: './src/postcss.config.js'
              }
            }
          }
        ]
      },
      {
        test: /\.tsx?$/,
        loader: 'ts-loader',
        options: {
          transpileOnly: true
        },
        options: {}
      },
      {
        type: 'javascript/auto',
        test: /\.m?js$/,
        resolve: { fullySpecified: false },
        use: []
      }
    ]
  },
  resolve: {
    extensions: ['.ts', '.tsx', '.js', '.json']
  },
  optimization: {
    splitChunks: {
      chunks: 'all',
      maxInitialRequests: 20, // for HTTP2
      maxAsyncRequests: 20, // for HTTP2
      cacheGroups: {
        service: {
          test: /[\\/]service/
        }
      }
    }
  }
};
