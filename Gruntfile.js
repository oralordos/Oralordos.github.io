module.exports = function(grunt) {
	grunt.initConfig({
		pkg: grunt.file.readJSON('package.json'),
		postcss: {
			options: {
				processors: [
					require('autoprefixer-core')({
						browsers: '> 5%',
						cascade: false
					})
				],
				diff: false
			},
			files: {
				src: ['*.css', 'SummerWebBootcamp/**/*.css']
			}
		},
		watch: {
			styles: {
				files: ['*.css'],
				tasks: ['postcss']
			}
		}
	});
	grunt.loadNpmTasks('grunt-postcss');
	grunt.loadNpmTasks('grunt-contrib-watch');

	grunt.registerTask('default', ['postcss'])
};