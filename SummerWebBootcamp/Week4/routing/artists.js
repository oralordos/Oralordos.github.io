var artistApp = angular.module('artists', ['firebase', 'ngAnimate']);

artistApp.controller('listController', ['$scope', '$firebaseArray', function($scope, $firebaseArray) {
    var ref = new Firebase('https://angular-bootcamp.firebaseio.com/');

    $scope.people = $firebaseArray(ref.child('people'));
}]);