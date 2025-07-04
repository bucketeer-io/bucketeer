require('dotenv').config();
const { merge } = require('webpack-merge');
const webpack = require('webpack');
const common = require('./webpack.common.js');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');

module.exports = merge(common, {
  mode: 'development',
  devServer: {
    host: 'localhost',
    port: 8000,
    historyApiFallback: true,
    allowedHosts: 'all'
  },
  plugins: [
    new HtmlWebpackPlugin({
      template: './src/index.html',
      favicon: './src/assets/favicon.ico'
    }),
    new MiniCssExtractPlugin(),
    new webpack.EnvironmentPlugin({
      RELEASE_CHANNEL: 'dev',
      DEV_WEB_API_ENDPOINT: process.env.DEV_WEB_API_ENDPOINT,
      DEV_AUTH_REDIRECT_ENDPOINT: process.env.DEV_AUTH_REDIRECT_ENDPOINT,
      NEW_CONSOLE_ENDPOINT: process.env.NEW_CONSOLE_ENDPOINT,
      DEMO_SIGN_IN_ENABLED: process.env.DEMO_SIGN_IN_ENABLED,
      DEMO_SIGN_IN_EMAIL: process.env.DEMO_SIGN_IN_EMAIL,
      DEMO_SIGN_IN_PASSWORD: process.env.DEMO_SIGN_IN_PASSWORD
    })
  ]
});
