var myApp = angular.module('ngAppOyster', [
	'ngRoute',
  'loginModule',
  'overviewModule',
]);

myApp.config(['$routeProvider',
  function($routeProvider) {
    $routeProvider.
      when('/overview', {
        templateUrl: 'views/overview.html',
        controller: 'overviewController',
        resolve: {
          logincheck: checkLogin,
        },
      }).
      when('/login', {
        templateUrl: 'views/login.html',
        controller: 'loginController',
      }).
      otherwise({
        redirectTo: '/login'
      });
  }]);

var checkLogin = function($q, $timeout, $http, $location, $rootScope) {
  var deferred = $q.defer();
  $http.get('/session').success(function(resp) {
    $rootScope.errorMessage = null;
    deferred.resolve();
    if ($rootScope.user == null) {
      $rootScope.user = resp;
    }
  }).error(function(resp) {
    deferred.reject();
    $.notify("Invalid login!", "error");
    $location.url('/login');
  });
  return deferred.promise;
};

myApp.controller("navBarController", function($scope, $rootScope) {});

myApp.controller("logoutController", function($scope, $http, $location, $rootScope) {
  $scope.logout = function() {
    if (confirm("Are you sure you want to logout?")) {
      $http.delete("/session").success(function() {
        $location.url("/login");
        $rootScope.user = null;
      });
    }
  };
});
