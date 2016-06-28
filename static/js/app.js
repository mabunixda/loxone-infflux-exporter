var app = angular.module('myApp', ['ngRoute']);

app.config(function($routeProvider) {
  $routeProvider

  .when('/', {
    templateUrl : 'pages/home.html',
    controller  : 'HomeController'
  })

  .when('/configuration', {
    templateUrl : 'pages/configuration.html',
    controller  : 'ConfigurationController'
  })  

  .otherwise({redirectTo: '/'});
});

app.controller('HomeController', function($scope) {
  $scope.message = 'Hello from HomeController';
});

app.controller('ConfigurationController', function($scope) {  
  loadLoxConfig();
  var bla = localStorage["LoxoneConfig"];
  loxConfig = JSON.parse(bla);       
  $scope.message =loxConfig.msInfo.msName;
});
