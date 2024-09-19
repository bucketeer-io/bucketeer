export default {
  typescript: {
    compile: false,
    rewritePaths: {
				"src/": "__lib/"
			},
  },
  files: ['__test/**/__tests__/*.js'],
};
