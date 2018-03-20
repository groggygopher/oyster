var module = angular.module('registerModule', []);


module.controller('registerController', ['$rootScope', '$scope', '$http', '$location', function($rootScope, $scope, $http, $location) {
  $scope.user = {};
  
  $scope.register = function(user) {
    $http.put('/session', angular.toJson(user)).success(function(resp) {
      $.notify("Welcome, " + resp.name, "success");
      $rootScope.user = resp;
      $location.url('/overview');
    }).error(function(resp) {
      $.notify(resp, "error");
      $rootScope.user = null;
    });
  };
}]);