var artistApp = angular.module('artists', ['firebase', 'ngAnimate']);

artistApp.controller('listController', ['$scope', '$firebaseArray', function($scope, $firebaseArray) {
    var ref = new Firebase('https://angular-bootcamp.firebaseio.com/');

    $scope.people = $firebaseArray(ref.child('people'));
    $scope.delete = function(itemId) {
        ref.child('people/' + itemId).remove();
    }
    $scope.add = function() {
        if($scope.newItem.name && $scope.newItem.shortname && $scope.newItem.reknown && $scope.newItem.bio) {
            ref.child('people').push($scope.newItem);
            $scope.newItem = {};
        }
    }
}]);

artistApp.controller('detailController', ['$scope', '$routeParams', '$firebaseObject', function($scope, $routeParams, $firebaseObject) {
    var ref = new Firebase('https://angular-bootcamp.firebaseio.com/people');
    $scope.person = $firebaseObject(ref.child($routeParams.person));
}]);