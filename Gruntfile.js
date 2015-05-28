module.exports = function(grunt) {
	grunt.initConfig({
		pkg: grunt.file.readJSON('package.json'),
		postcss: {
			options: {
				processors: [
					require('autoprefixer-core')({
						browsers: 'last 2 versions'
					})
				],
				diff: false
			},
			files: {
				src: 'css/*.css'
			}
		},
		watch: {
			styles: {
				files: ['css/*.css'],
				tasks: ['postcss']
			}
		}
	});
	grunt.loadNpmTasks('grunt-postcss');
	grunt.loadNpmTasks('grunt-contrib-watch');

	grunt.registerTask('default', ['watch'])
};