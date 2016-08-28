'use strict';

angular.module('myApp.scoreboard', ['ngRoute'])

.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/user', {
    templateUrl: 'scoreboard/scoreboard.html',
    controller: 'ScoreboardCtrl'
  });
}])

.controller('ScoreboardCtrl', ['$sce', '$scope', '$http', 'UserService', function($sce, $scope, $http, UserService) {

	/* ------ BEGIN INIT ------ */
	$scope.request = {};
	$http.get('/v1.0/user').success( function ( users ) {
		$scope.users = users;
		for(var userIt=0 ; userIt<users.length ; ++userIt){		
			$scope.users[userIt].score = 0;
			var currId = $scope.users[userIt].ID;
			$http.get('/v1.0/user/'+currId+'/validatedChallenges').success((function(userIterator){
				return function(validatedChalls){
					for (var challIt=0; challIt < validatedChalls.length ; ++challIt){
						$http.get('/v1.0/challenge/' + validatedChalls[challIt].ChallengeID).success((function(userIter){
							return function(validatedChall){
								$scope.users[userIter].score += validatedChall.points;
							}
						})(userIterator)).error(function(error){
							alert("An error occured : " + error.message);
						});
					}
				}
			})(userIt)).error(function(error){
				alert('An error occured :' + error.message);
			});
		}

	}).error(function(error){
		alert("An error occured while processing request : " + data.error);
	});
	/* ------ END INIT ------ */
}]);
