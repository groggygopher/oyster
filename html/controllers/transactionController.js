var module = angular.module('transactionModule', []);


module.controller('transactionController', ['$scope', '$http', function($scope, $http) {
  $scope.transactions = [];

  $scope.upload = function() {
    var f = document.getElementById('file').files[0];
    var r = new FileReader();

    r.onloadend = function(e) {
      $http({
        url: "/upload",
        method: "POST",
        data: e.target.result,
      }).success(function (resp) {
        $.notify("Uploaded " + resp.uploaded + " transactions, imported " + resp.imported + " new transactions", "success");
        $scope.update();
      }).error(function (response) {
        $.notify(response, "error");
      });
    }
    r.readAsBinaryString(f);
  }

  $scope.update = function() {
    $http.get('/transactions').success(function(resp) {
      $scope.transactions = resp;
    }).error(function(resp) {
      $.notify(resp, "error");
    });
  }
  $scope.update();

  $scope.formatDate = function(dateStr) {
    var d = new Date(dateStr);
    return d.toDateString();
  }

  $scope.formatCurrency = function(amtDbl) {
    if (amtDbl < 0) {
      return "-$" + Math.abs(amtDbl).toFixed(2);
    }
    return "$" + amtDbl.toFixed(2);
  }

}]);