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

  .run(function($rootScope, WEEKDAYS, clientService) {
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

    clientService.start();
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
      },
      clear: function(weekday, time) {
        client.data = {
          name: "",
          classes: []
        };

        if (weekday && time) {
          client.data.classes.push({
            weekday: weekday,
            time: time.toDate(),
            duration: "30",
          });
        }
      }
    };
  })

  .service("clients", function() {
    var clients = {
      data: []
    };

    return {
      getClients: function() {
        return clients;
      },
      setClients: function(c) {
        if (Array.isArray(c)) {
          clients.data = c;
        } else {
          console.log("Trying to set a non-array in clients", c);
        }
      },
      addClient: function(c) {
        var newClient = true;
        clients.data.some(function(client, index) {
          if (c.id && client.id == c.id) {
            clients.data[index] = c;
            newClient = false;
            return true;

          } else if (c.temporaryId && client.temporaryId == c.temporaryId) {
            clients.data[index] = c;
            newClient = false;
            return true;
          }

          return false;
        });

        if (newClient) {
          clients.data.push(c);
        }
      }
    };
  })

  .service("clientService", function($http, $q, $timeout, clients, messages) {
    var saveLaterClients = [];

    return({
      retrieveAll: retrieveAll,
      save: save,
      start: function() {
        $timeout(start, 5000);
      }
    });

    function retrieveAll() {
      var request = $http({
        method: "GET",
        url: "/clients"
      });

      request.then(
        function(r) {
          convertJSONDate(r.data);
          clients.setClients(r.data);
          localStorage.setItem("clients", angular.toJson(clients.getClients().data));
        },
        function(r) {
          if (r.status == 400) {
            messages.setMessages(r.data);
          } else {
            console.log("Error", r.status, "while retrieving clients.", r.data);
          }

          var c = angular.fromJson(localStorage.getItem("clients"));
          convertJSONDate(c);
          clients.setClients(c);
        }
      );
    }

    function save(client) {
      var request;
      var newClient = false;

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
        newClient = true;
      }

      return request.then(
        function(r) {
          if (newClient) {
            client.id = r.headers("Location").slice(8);
          }

          clients.addClient(client);
          localStorage.setItem("clients", angular.toJson(clients.getClients().data));

          return {
            success: true,
            saveLater: false
          };
        },
        function(r) {
          if (r.status == 400) {
            messages.setMessages(r.data);

            return {
              success: false,
              saveLater: false
            };

          } else {
            if (!client.id && !client.temporaryId) {
              client.temporaryId = (Math.random() + 1).toString(36).substring(7);
            }

            clients.addClient(client);
            saveLater(client);
            localStorage.setItem("clients", angular.toJson(clients.getClients().data));

            console.log("Error", r.status, "while saving client.", r.data);

            return {
              success: false,
              saveLater: true
            };
          }
        }
      );
    }

    function convertJSONDate(clients) {
      clients.forEach(function(c) {
        if (!c || !c.classes) {
          return;
        }

        c.classes.forEach(function(cl) {
          if (cl.time) {
            cl.time = moment(cl.time).toDate();
          }
        });
      });
    }

    function saveLater(c) {
      var newClient = true;
      saveLaterClients.some(function(client, index) {
        if (c.id && client.id == c.id) {
          saveLaterClients[index] = c;
          newClient = false;
          return true;

        } else if (c.temporaryId && client.temporaryId == c.temporaryId) {
          saveLaterClients[index] = c;
          newClient = false;
          return true;
        }

        return false;
      });

      if (newClient) {
        saveLaterClients.push(c);
      }
    }

    function start() {
      var savingClients = angular.copy(saveLaterClients);
      savingClients.forEach(function(client, index) {
        save(client)
          .then(
            function(r) {
              if (r.success) {
                saveLaterClients.splice(index, 1);
              }
            },
            function() {}
          );
      });
      $timeout(start, 5000);
    }
  })

  .service("messages", function() {
    var messages = {
      data: []
    };

    return {
      getMessages: function() {
        return messages;
      },
      setMessages: function(m) {
        if (Array.isArray(m)) {
          messages.data = m;
        } else {
          console.log("Trying to set a non-array in messages", m);
        }
      }
    }
  })

  .controller("scheduleCtrl", function($rootScope, $scope, client, clients, clientService) {
    $scope.clients = clients.getClients();
    $scope.times = [
      "05:00", "05:30", "06:00", "06:30", "07:00", "07:30", "08:00", "08:30", "09:00", "09:30",
      "10:00", "10:30", "11:00", "11:30", "12:00", "12:30", "13:00", "13:30", "14:00", "14:30",
      "15:00", "15:30", "16:00", "16:30", "17:00", "17:30", "18:00", "18:30", "19:00", "19:30",
      "20:00", "20:30", "21:00", "21:30", "22:00", "22:30", "23:00", "23:30"
    ];

    // Returns a list of client objects
    $scope.clientsAt = function(time, weekday) {
      var filteredClients = [];
      var current = moment("1970-01-01 " + time);

      $scope.clients.data.forEach(function(client) {
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
            filteredClients.push(client);
            return true;
          }

          return false;
        });
      });
      return filteredClients;
    };

    $scope.editClient = function(event, c) {
      event.stopPropagation();
      client.setClient(c);
      $rootScope.clientFormMode = true;
    };

    $scope.newClient = function(weekday, time) {
      client.clear(weekday, moment("1970-01-01 " + time));
      $rootScope.clientFormMode = true;
    };

    $scope.clientColor = function(c) {
      if (!c.id && !c.temporaryId) {
        return "#fff";
      }

      var id = c.id;
      if (!c.id) {
        id = c.temporaryId;
      }

      var hash = 0;
      for (var i = 0; i < id.length; i++) {
         hash = id.charCodeAt(i) + ((hash << 5) - hash);
      }

      var colour = "#";
      for (var i = 0; i < 3; colour += ("00" + ((hash >> i++ * 8) & 0xFF).toString(16)).slice(-2));
      return colour;
    };

    clientService.retrieveAll();
  })

  .controller("clientFormCtrl", function($rootScope, $scope, WEEKDAYS, client, 
    clients, clientService, messages) {

    $scope.client = client.getClient();
    $scope.messages = messages.getMessages();

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
          function(r) {
            if (r.success || r.saveLater) {
              $rootScope.clientFormMode = false;
            }
          },
          function(r) {
            if (r.success || r.saveLater) {
              $rootScope.clientFormMode = false;
            }
          }
        );
    };
  });