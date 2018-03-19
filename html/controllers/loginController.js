var module = angular.module('loginModule', []);

module.controller('loginController', ['$scope', '$http', '$rootScope', '$location', function($scope, $http, $rootScope, $location) {
	$scope.user = {};
	
	$scope.login = function(user) {
		$http.post('/session', angular.toJson(user)).success(function(resp) {
			$.notify("Welcome, " + resp.name, "success");
			$rootScope.user = resp;
			$location.url('/overview');
		}).error(function(resp) {
			$.notify("Invalid login", "error");
			$rootScope.user = null;
		});
	};
}]);