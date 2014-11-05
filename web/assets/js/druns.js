angular.module("druns", [])
	.constant("WEEKDAYS", {
		sunday: "sunday",
		monday: "monday",
		tuesday: "tuesday",
		wednesday: "wednesday",
		thursday: "thursday",
		friday: "friday",
		saturday: "saturday"
	})

	.run(function($rootScope, WEEKDAYS) {
		$rootScope.WEEKDAYS = WEEKDAYS
   })

	.service("clientService", function($http, $q) {
		return({
			retrieveAll: retrieveAll,
			save: save
		})

		function retrieveAll() {
			var request = $http({
				method: "GET",
				url: "/clients"
			});

			return request.then(handleSuccess, handleError);
		}

		function save(client) {
			var request;

			if (client.id && client.id.length > 0) {
				request = $http({
					method: "PUT",
					url: "/client/" + client.id,
					data: client
				});

			} else {
				request = $http({
					method: "POST",
					url: "/client",
					data: client
				});
			}

			return request.then(handleSuccess, handleError);
		}

		function handleError(response) {
			if (!angular.isObject(response.data) || !response.data.message) {
				return $q.reject("An unknown error occurred.");
			}
			return $q.reject(response.data.message);
		}

		function handleSuccess(response) {
			return response.data;
		}
	})

	.controller("scheduleCtrl", function($scope, clientService) {
		$scope.clients = [];
		$scope.times = [
			"5:00", "5:30", "6:00", "6:30", "7:00", "7:30", "8:00", "8:30", "9:00", "9:30", "10:00",
			"10:30", "11:00", "11:30", "12:00", "12:30", "13:00", "13:30", "14:00", "14:30", "15:00",
			"15:30", "16:00", "16:30", "17:00", "17:30", "18:00", "18:30", "19:00", "19:30", "20:00",
			"20:30", "21:00", "21:30", "22:00", "22:30", "23:00", "23:30"
		];

		// Returns a list of client objects
		$scope.clientsAt = function(time, dayOfTheWeek) {
			return [];
		};

		$scope.retrieveClients = function() {
			clientService.retrieveAll()
				.then(
					function(clients) {
						$scope.clients = clients;
					}
				);
		};

		$scope.retrieveClients();
	})

	.controller("clientFormCtrl", function($rootScope, $scope, WEEKDAYS, clientService) {
		$scope.client = {
			name: "",
			classes: []
		};

		$scope.addClass = function() {
			$scope.client.classes.push({
				weekday: WEEKDAYS.sunday,
				time: new Date(1970, 0, 1, 5, 0, 0),
				duration: "30",
			});
		};

		$scope.save = function() {
			clientService.save($scope.client)
				.then(
					function() {
						// TODO: Success
						$rootScope.clientFormMode = false
					}
				);
		};
	});