export default {
  typescript: {
    compile: false,
    rewritePaths: {
				"src/": "__lib/"
			},
  },
  files: ['/**/__tests__/*.js'],
};
