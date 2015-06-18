var artistApp = angular.module('artists', ['firebase', 'ngAnimate']);

artistApp.controller('listController', ['$scope', '$firebaseArray', function($scope, $firebaseArray) {
    var ref = new Firebase('https://angular-bootcamp.firebaseio.com/');

    $scope.people = $firebaseArray(ref.child('people'));
}]);

artistApp.controller('detailController', ['$scope', '$routeParams', '$firebaseObject', function($scope, $routeParams, $firebaseObject) {
    var ref = new Firebase('https://angular-bootcamp.firebaseio.com/people');
    $scope.person = $firebaseObject(ref.child($routeParams.person));
}]);