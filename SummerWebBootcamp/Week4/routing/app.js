myApp = angular.module('myApp', ['ngRoute', 'artists']);

myApp.config(['$routeProvider', function($routeProvider) {
    $routeProvider.when('/list', {
        templateUrl: 'list.html',
        controller: 'listController'
    });
    $routeProvider.otherwise({
        redirectTo: '/list'
    });
}]);