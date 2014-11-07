angular.module("druns", [])
	.constant("WEEKDAYS", {
		sunday: "Sunday",
		monday: "Monday",
		tuesday: "Tuesday",
		wednesday: "Wednesday",
		thursday: "Thursday",
		friday: "Friday",
		saturday: "Saturday"
	})

	.run(function($rootScope, WEEKDAYS) {
		$rootScope.WEEKDAYS = WEEKDAYS
		$rootScope.WEEKDAYS_LIST = [
			WEEKDAYS.sunday,
			WEEKDAYS.monday,
			WEEKDAYS.tuesday,
			WEEKDAYS.wednesday,
			WEEKDAYS.thursday,
			WEEKDAYS.friday,
			WEEKDAYS.saturday
		];
   })

	.service("client", function() {
		var client = {
			data: {
				name: "",
				classes: []
			}
		};

		return {
			getClient: function() {
				return client;
			},
			setClient: function(c) {
				client.data = c;
			}
		}
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
			if (!angular.isObject(response.data)) {
				return $q.reject("An unknown error occurred.");
			}
			return $q.reject(response.data.message);
		}

		function handleSuccess(response) {
			return response.data;
		}
	})

	.controller("scheduleCtrl", function($rootScope, $scope, client, clientService) {
		$scope.clients = [];
		$scope.times = [
			"05:00", "05:30", "06:00", "06:30", "07:00", "07:30", "08:00", "08:30", "09:00", "09:30",
			"10:00", "10:30", "11:00", "11:30", "12:00", "12:30", "13:00", "13:30", "14:00", "14:30",
			"15:00", "15:30", "16:00", "16:30", "17:00", "17:30", "18:00", "18:30", "19:00", "19:30",
			"20:00", "20:30", "21:00", "21:30", "22:00", "22:30", "23:00", "23:30"
		];

		// Returns a list of client objects
		$scope.clientsAt = function(time, weekday) {
			var clients = [];
			var current = moment("1970-01-01 " + time);

			$scope.clients.forEach(function(client) {
				if (!client.classes) {
					return;
				}

				client.classes.some(function(c) {
					if (c.weekday != weekday) {
						return false;
					}

					var begin = moment(c.time);
					var end = angular.copy(begin);
					end.add(c.duration, "minutes");

					if (begin <= current && current < end) {
						clients.push(client);
						return true;
					}

					return false;
				});
			});
			return clients;
		};

		$scope.retrieveClients = function() {
			clientService.retrieveAll()
				.then(
					function(clients) {
						// Convert JSON string time to Date object
						clients.forEach(function(client) {
							if (!client || !client.classes) {
								return;
							}

							client.classes.forEach(function(c) {
								if (c.time) {
									c.time = moment(c.time).toDate();
								}
							});
						});
						$scope.clients = clients;
					}
				);
		};

		$scope.editClient = function(c) {
			client.setClient(c);
			$rootScope.clientFormMode = true;
		};

		$scope.clientColor = function(c) {
			var hash = 0;
	    for (var i = 0; i < c.id.length; i++) {
	       hash = c.id.charCodeAt(i) + ((hash << 5) - hash);
	    }

	    var colour = "#";
	    for (var i = 0; i < 3; colour += ("00" + ((hash >> i++ * 8) & 0xFF).toString(16)).slice(-2));
	    return colour;
		};

		$scope.retrieveClients();
	})

	.controller("clientFormCtrl", function($rootScope, $scope, WEEKDAYS, client, clientService) {
		$scope.client = client.getClient();

		$scope.addClass = function() {
			$scope.client.data.classes.push({
				weekday: WEEKDAYS.sunday,
				time: new Date(1970, 0, 1, 5, 0, 0),
				duration: "30",
			});
		};

		$scope.removeClass = function(index) {
			$scope.client.data.classes.splice(index, 1);
		};

		$scope.save = function() {
			clientService.save($scope.client.data)
				.then(
					function() {
						// TODO: Success
						$rootScope.clientFormMode = false;
					}
				);
		};
	});