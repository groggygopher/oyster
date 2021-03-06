var myApp = angular.module('ngAppOyster', [
	'ngRoute',
  'loginModule',
  'overviewModule',
  'registerModule',
  'transactionModule',
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
      when('/transactions', {
        templateUrl: 'views/transaction.html',
        controller: 'transactionController',
        resolve: {
          logincheck: checkLogin,
        },
      }).
      when('/login', {
        templateUrl: 'views/login.html',
        controller: 'loginController',
      }).
      when('/register', {
        templateUrl: 'views/register.html',
        controller: 'registerController',
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
      }).error(function(resp) {
        $.notify(resp, "error");
      });
    }
  };
});

myApp.directive('customOnChange', function() {
  return {
    restrict: 'A',
    link: function (scope, element, attrs) {
      var onChangeFunc = scope.$eval(attrs.customOnChange);
      element.bind('change', onChangeFunc);
    }
  };
});