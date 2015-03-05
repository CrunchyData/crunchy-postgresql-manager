(function() {
    var app = angular.module('BackupApp', ['ngTable', 'ngCookies', 'ui.bootstrap' ]);


 app.run(function($rootScope) {
    $rootScope.$on('LoadingEvent', function(event, args) {
        $rootScope.$broadcast('LoadingEventTarget', args);
        $rootScope.isLoading = true;
    });
    $rootScope.$on('DoneLoadingEvent', function(event, args) {
        $rootScope.$broadcast('DoneLoadingEventTarget', args);
        $rootScope.isLoading = false;
    });
    $rootScope.$on('deleteScheduleEvent', function(event, args) {
        $rootScope.$broadcast('deleteScheduleTarget', args);
    });
    $rootScope.$on('updateSchedulePage', function(event, args) {
        $rootScope.$broadcast('updateSchedulePageTarget', args);
    });
    $rootScope.$on('noSchedule', function(event, args) {
        $rootScope.$broadcast('noScheduleTarget', args);
    });
    $rootScope.$on('setSchedule', function(event, args) {
        $rootScope.$broadcast('setScheduleTarget', args);
    });
    $rootScope.$on('createScheduleEvent', function(event, args) {
        $rootScope.$broadcast('createScheduleTarget', args);
        console.log('broadcast of createScheduleTarget');
    });
});

    
   app.controller('GetContainerController', function($scope, $http, $rootScope, $modal, $cookies, $cookieStore, $filter, ngTableParams) {

	$scope.oneAtATime = true;
	$scope.status = {
		isFirstOpen: true,
		isFirstDisabled: false
	};
    $scope.currentUser = [];

    $scope.currentUser = $cookieStore.get('cpmuser');

    console.log('currentUser is ' + $scope.currentUser);
    $scope.profiles = [{name: 'pg_basebackup'}, {name:'pg_other'}];
    $scope.currentProfileName = $scope.profiles[0];
    //we are using NO and YES for the enabled checkbox value
    $scope.isLoading = false;
    $scope.results = [];
    $scope.clusters = [];
    $scope.myCluster = [];
    $scope.myServer = [];
    $scope.servers = [];
    $scope.currentContainer = [];

    $scope.stats = [];
    $scope.selectedStats = [];
    $scope.users = [];
    $scope.data = [];

    //keeps track of currently selected values onscreen
    $scope.dowlist = {};
    $scope.hourlist = {};
    $scope.minlist = {};
    $scope.monthlist = {};
    $scope.domlist = {};
    $scope.domlist2 = {};


    $scope.thething = [ {'name': 'thething', 'checked' : false}];

   $scope.dayofweek = [ 
   { 'name': 'SUN', 'checked': false }, 
	{ 'name': 'MON', 'checked': false }, 
	{'name': 'TUE', 'checked': false }, 
	{'name': 'WED', 'checked': false },
	{'name': 'THU', 'checked': false },
	{'name': 'FRI', 'checked': false },
	{'name': 'SAT', 'checked': false },];

   $scope.hours = [ 
   { 'name': '0', 'checked': false }, 
	{ 'name': '1', 'checked': false }, 
	{'name': '2', 'checked': false }, 
	{'name': '3', 'checked': false },
	{'name': '4', 'checked': false },
	{'name': '5', 'checked': false },
	{'name': '6', 'checked': false },
	{'name': '7', 'checked': false },
	{'name': '8', 'checked': false },
	{'name': '9', 'checked': false },
	{'name': '10', 'checked': false },
	{'name': '11', 'checked': false },
	{'name': '12', 'checked': false },
	{'name': '13', 'checked': false },
	{'name': '14', 'checked': false },
	{'name': '15', 'checked': false },
	{'name': '16', 'checked': false },
	{'name': '17', 'checked': false },
	{'name': '18', 'checked': false },
	{'name': '19', 'checked': false },
	{'name': '20', 'checked': false },
	{'name': '21', 'checked': false },
	{'name': '22', 'checked': false },
	{'name': '23', 'checked': false },
	];

   $scope.theminutes = [ 
   { 'name': '00', 'checked': false }, 
	{ 'name': '05', 'checked': false }, 
	{'name': '10', 'checked': false }, 
	{'name': '15', 'checked': false },
	{'name': '20', 'checked': false },
	{'name': '25', 'checked': false },
	{'name': '30', 'checked': false },
	{'name': '35', 'checked': false },
	{'name': '40', 'checked': false },
	{'name': '45', 'checked': false },
	{'name': '50', 'checked': false },
	{'name': '55', 'checked': false },
	];

   $scope.themonths = [ 
   { 'name': 'Jan', 'checked': false }, 
	{ 'name': 'Feb', 'checked': false }, 
	{'name': 'Mar', 'checked': false }, 
	{'name': 'Apr', 'checked': false },
	{'name': 'May', 'checked': false },
	{'name': 'Jun', 'checked': false },
	{'name': 'Jul', 'checked': false },
	{'name': 'Aug', 'checked': false },
	{'name': 'Sep', 'checked': false },
	{'name': 'Oct', 'checked': false },
	{'name': 'Nov', 'checked': false },
	{'name': 'Dec', 'checked': false },
	];

   $scope.dayofmonth = [ 
   { 'name': '0', 'checked': false }, 
	{ 'name': '1', 'checked': false }, 
	{'name': '2', 'checked': false }, 
	{'name': '3', 'checked': false },
	{'name': '4', 'checked': false },
	{'name': '5', 'checked': false },
	{'name': '6', 'checked': false },
	{'name': '7', 'checked': false },
	{'name': '8', 'checked': false },
	{'name': '9', 'checked': false },
	{'name': '10', 'checked': false },
	{'name': '11', 'checked': false },
	{'name': '12', 'checked': false },
	{'name': '13', 'checked': false },
	{'name': '14', 'checked': false },
	{'name': '15', 'checked': false },
	{'name': '16', 'checked': false },
	];
   $scope.dayofmonth2 = [ 
   { 'name': '17', 'checked': false }, 
	{ 'name': '18', 'checked': false }, 
	{'name': '19', 'checked': false }, 
	{'name': '20', 'checked': false },
	{'name': '21', 'checked': false },
	{'name': '22', 'checked': false },
	{'name': '23', 'checked': false },
	{'name': '24', 'checked': false },
	{'name': '25', 'checked': false },
	{'name': '26', 'checked': false },
	{'name': '27', 'checked': false },
	{'name': '28', 'checked': false },
	{'name': '29', 'checked': false },
	{'name': '30', 'checked': false },
	{'name': '31', 'checked': false },
	];
   $scope.selectdom = function() {
	   $scope.domlist = $filter('filter')($scope.dayofmonth, 
			   {checked:true});   
	   console.log('selected domlist=' + $scope.domlist);
   }
   $scope.selectdom2 = function() {
	   $scope.domlist2 = $filter('filter')($scope.dayofmonth2, 
			   {checked:true});   
	   console.log('selected domlist2=' + $scope.domlist2);
   }

   $scope.selectdow = function() {
	   $scope.dowlist = $filter('filter')($scope.dayofweek, 
			   {checked:true});   
	   console.log('selected dowlist=' + $scope.dowlist);
   }

   $scope.selecthour = function() {
	   $scope.hourlist = $filter('filter')($scope.hours, 
			   {checked:true});   
	   console.log('selected hourlist=' + $scope.hourlist);
   }

   $scope.selectmin = function() {
	   $scope.minlist = $filter('filter')($scope.theminutes, 
			   {checked:true});   
	   console.log('selected minlist=' + $scope.minlist);
   }

   $scope.selectmonth = function() {
	   $scope.monthlist = $filter('filter')($scope.themonths, 
			   {checked:true});   
	   console.log('selected monthlist=' + $scope.monthlist);
   }

    $scope.containerid = window.containerid;

        var token = $cookieStore.get('cpm_token');

    $http.get($cookieStore.get('AdminURL') + '/servers/' + token).
    success(function(data, status, headers, config) {
        $scope.servers = data;
        console.log('got servers len=' + data.length);
        $scope.myServer = $scope.servers[0];
    }).error(function(data, status, headers, config) {
        console.log('error in fetch of servers');
    });

     $http.get($cookieStore.get('AdminURL') + '/node/' + window.containerid + '.' + token).success(function(data, status, headers, config) {
		$scope.currentContainer = data;
	}).error(function(data, status, headers, config) {
		alert('error in get container');
	});


   $scope.tableParams = new ngTableParams({
        page: 1, // show first page
        count: 10 // count per page
    }, {
        total: $scope.stats.length, // length of data
        getData: function($defer, params) {
            console.log('getData called stats=' + $scope.stats.length);
            // use build-in angular filter
            var orderedData = $scope.stats;

            params.total(orderedData.length); // set total for recalc pagination
            $defer.resolve($scope.users = orderedData.slice((params.page() - 1) * params.count(), params.page() * params.count()));
        }
    });


  //fix around ng-table bug?
    $scope.tableParams.settings().$scope = $scope;

    $scope.checkboxes = {
        'checked': false,
        items: {}
    };

    // watch for check all checkbox
    $scope.$watch('checkboxes.checked', function(value) {
        angular.forEach($scope.users, function(item) {
            if (angular.isDefined(item.ID)) {
                $scope.checkboxes.items[item.ID] = value;
            }
        });
    });

// watch for data checkboxes
    $scope.$watch('checkboxes.items', function(values) {
        if (!$scope.users) {
            return;
        }
        var checked = 0,
            unchecked = 0,
            total = $scope.users.length;
        angular.forEach($scope.users, function(item) {
            checked += ($scope.checkboxes.items[item.ID]) || 0;
            unchecked += (!$scope.checkboxes.items[item.ID]) || 0;
        });
        if ((unchecked == 0) || (checked == 0)) {
            $scope.checkboxes.checked = (checked == total);
        }
        // grayed checkbox
        angular.element(document.getElementById("select_all")).prop("indeterminate", (checked != 0 && unchecked != 0));
    }, true);


    $scope.handleRefresh = function() {
        console.log('refresh now on ' + $scope.currentSchedule.Name);
	postit($scope.currentSchedule);
    };

    $scope.handleBackupNowClick = function() {
        console.log('backup now on ' + $scope.currentContainer.Name);

        var modalInstance = $modal.open({
            templateUrl: 'backupnowmodal.html',
            controller: BackupNowCtrl,
            resolve: {
                value: function() {
                    return $scope.currentSchedule;
                }
            }
        });
        modalInstance.result.then(function(response) {});
    };

    $scope.handleCreateClick = function() {
        console.log('create now on ' + $scope.currentContainer.Name);

        var modalInstance = $modal.open({
            templateUrl: 'createschedulemodal.html',
            controller: CreateScheduleCtrl,
            resolve: {
                value: function() {
                    return $scope.currentContainer;
                }
            }
        });
        modalInstance.result.then(function(response) {});

    };

    $scope.handleMinusClick = function() {
        console.log('delete schedule called on ' + $scope.currentSchedule.Name);

        var modalInstance = $modal.open({
            templateUrl: 'deleteschedulemodal.html',
            controller: DeleteScheduleCtrl,
            resolve: {
                value: function() {
                    return $scope.currentSchedule;
                }
            }
        });
        modalInstance.result.then(function(response) {});

    };

    $scope.handleUpdateClick = function() {
        console.log('update now on ' + $scope.currentContainer.Name);

	updateCurrentSchedule();

        var modalInstance = $modal.open({
            templateUrl: 'updateschedulemodal.html',
            controller: UpdateScheduleCtrl,
            resolve: {
                value: function() {
                    return $scope.currentSchedule;
                }
            }
        });
        modalInstance.result.then(function(response) {});

    };

	function updateCurrentSchedule() {
		console.log('updateCurrentSchedule called');
		var d = [];
		for (var i=0; i<$scope.dowlist.length; i++) {
			d[d.length] = $scope.dowlist[i].name;
		}
		console.log('dow string=' + d.toString());
		$scope.currentSchedule.DayOfWeek = d.toString();
		
		var d2 = [];
		for (var i=0; i<$scope.monthlist.length; i++) {
			d2[d2.length] = $scope.monthlist[i].name;
		}
		console.log('month string=' + d2.toString());
		$scope.currentSchedule.Month = d2.toString();

		var d3 = [];
		for (var i=0; i<$scope.hourlist.length; i++) {
			d3[d3.length] = $scope.hourlist[i].name;
		}
		console.log('hourlist string=' + d3.toString());
		$scope.currentSchedule.Hours = d3.toString();

		var d4 = [];
		for (var i=0; i<$scope.minlist.length; i++) {
			d4[d4.length] = $scope.minlist[i].name;
		}
		console.log('minlist string=' + d4.toString());
		$scope.currentSchedule.Minutes = d4.toString();

		var d5 = [];
		for (var i=0; i<$scope.domlist.length; i++) {
			d5[d5.length] = $scope.domlist[i].name;
		}
		for (var i=0; i<$scope.domlist2.length; i++) {
			d5[d5.length] = $scope.domlist2[i].name;
		}
		console.log('dom string=' + d5.toString());
		$scope.currentSchedule.DayOfMonth = d5.toString();

		console.log('enabled flag here is ' + $scope.thething.checked);
		if ($scope.thething.checked == true) {
		$scope.currentSchedule.Enabled = 'YES';
		} else {
		$scope.currentSchedule.Enabled = 'NO';
		}
		console.log('ccurrent enabled flag here is ' + $scope.currentSchedule.Enabled);
	}

 	function postit(schedule) {
		console.log('schedule posted is ' + schedule.Name);
        	var token = $cookieStore.get('cpm_token');

			clearSchedule();

     		$http.get($cookieStore.get('AdminURL') + '/backup/getschedule/' + schedule.ID + '.' + token).success(function(data, status, headers, config) {
			$scope.currentSchedule = data;
			console.log('got schedule data ' + $scope.currentSchedule.Name);
			updateScheduleOnScreen();
		}).error(function(data, status, headers, config) {
			alert('error in get schedule');
		});

		console.log('calling getallstatus with id=' + schedule.ID);
		$http.get($cookieStore.get('AdminURL') + '/backup/getallstatus/' + schedule.ID + "." + token).success(function(data2, status, headers, config) {
                    	$scope.stats = data2;
            	}).error(function(data, status, headers, config) {
                    	console.log('error:GetContainer.postit.getallstatus');
            	});

	}

	function clearSchedule() {
		console.log('clearSchedule called');

		$scope.thething.checked = false;

		for (var x=0; x<$scope.dayofweek.length; x++) {
			$scope.dayofweek[x].checked = false;
		}
		for (var x=0; x<$scope.themonths.length; x++) {
			$scope.themonths[x].checked = false;
		}
		for (var x=0; x<$scope.hours.length; x++) {
			$scope.hours[x].checked = false;
		}
		for (var x=0; x<$scope.theminutes.length; x++) {
			$scope.theminutes[x].checked = false;
		}
		for (var x=0; x<$scope.dayofmonth.length; x++) {
			$scope.dayofmonth[x].checked = false;
		}
		for (var x=0; x<$scope.dayofmonth2.length; x++) {
			$scope.dayofmonth2[x].checked = false;
		}
	}


	function updateScheduleOnScreen() {

		if ($scope.currentSchedule.Enabled == 'YES') {
			console.log('setting enabled flag to true');
		$scope.thething.checked = true;
		} else {
			console.log('setting enabled flag to false');
		$scope.thething.checked = false;
		}

		console.log('updateScheduleOnScreen enabled=[' + $scope.thething.checked + ']');
		console.log('updateScheduleOnScreen called');
		arr = $scope.currentSchedule.DayOfWeek.split(',');
		for (var i=0; i<arr.length; i++) {
			console.log(' dow = ' + arr[i]);
			for (var x=0; x<$scope.dayofweek.length; x++) {
				if ($scope.dayofweek[x].name == arr[i]) {
					$scope.dayofweek[x].checked = true;
					x = $scope.dayofweek.length;
				}
			}
		}


		arr = $scope.currentSchedule.Month.split(',');
		for (var i=0; i<arr.length; i++) {
			console.log(' month = ' + arr[i]);
			for (var x=0; x<$scope.themonths.length; x++) {
				if ($scope.themonths[x].name == arr[i]) {
					$scope.themonths[x].checked = true;
					x = $scope.themonths.length;
				}
			}
		}

		arr = $scope.currentSchedule.Hours.split(',');
		for (var i=0; i<arr.length; i++) {
			console.log(' hours = ' + arr[i]);
			for (var x=0; x<$scope.hours.length; x++) {
				if ($scope.hours[x].name == arr[i]) {
					$scope.hours[x].checked = true;
					x = $scope.hours.length;
				}
			}
		}

		arr = $scope.currentSchedule.Minutes.split(',');
		for (var i=0; i<arr.length; i++) {
			console.log(' Minutes = ' + arr[i]);
			for (var x=0; x<$scope.theminutes.length; x++) {
				if ($scope.theminutes[x].name == arr[i]) {
					$scope.theminutes[x].checked = true;
					x = $scope.theminutes.length;
				}
			}
		}

		arr = $scope.currentSchedule.DayOfMonth.split(',');
		for (var i=0; i<arr.length; i++) {
			console.log(' DayOfMonth = ' + arr[i]);
			for (var x=0; x<$scope.dayofmonth.length; x++) {
				if ($scope.dayofmonth[x].name == arr[i]) {
					$scope.dayofmonth[x].checked = true;
					x = $scope.dayofmonth.length;
				}
			}
			for (var x=0; x<$scope.dayofmonth2.length; x++) {
				if ($scope.dayofmonth2[x].name == arr[i]) {
					$scope.dayofmonth2[x].checked = true;
					x = $scope.dayofmonth2.length;
				}
			}
		}

		$scope.selectdom();
		$scope.selectdom2();
		$scope.selectdow();
		$scope.selecthour();
		$scope.selectmin();
		$scope.selectmonth();
		console.log('setting current profile name to ' + $scope.currentSchedule.ProfileName);
		console.log('test ' + $scope.currentProfileName);
		$scope.currentProfileName.name = $scope.currentSchedule.ProfileName;
		console.log('it is ' + $scope.currentProfileName.name);
	}


 	$rootScope.$on('setScheduleTarget', function(event, args) {
        	console.log('setScheduleTarget ' + args.message.ID);
        	postit(args.message);
    	});
 	$rootScope.$on('noScheduleTarget', function(event, args) {
        	console.log('no schedule event received here ');
        	$scope.currentSchedule = [];
    	});

   });      



    app.controller('GetSchedulesController', function($rootScope, $scope, $http, $modal, $cookies, $cookieStore) {
	$scope.oneAtATime = true;
	$scope.status = {
		isFirstOpen: true,
		isFirstDisabled: false
	};
    $scope.results = [];
    $scope.currentContainer = [];
    $scope.isLoading = false;
    $scope.containerid = window.containerid;

        var token = $cookieStore.get('cpm_token');

     $http.get($cookieStore.get('AdminURL') + '/node/' + window.containerid + '.' + token).success(function(data, status, headers, config) {
		$scope.currentContainer = data;
	}).error(function(data, status, headers, config) {
		alert('error in get container');
	});
    this.tab = 1;
    this.isSelected2 = function(checkTab) {
        return this.tab === checkTab;
    };
    this.selectTab2 = function(setTab) {
        console.log('setting tab to ' + setTab.ID);
        this.tab = setTab.ID;
        $scope.activeClass = setTab.ID;
        $scope.currentContainer = setTab;
        $rootScope.$emit('setSchedule', {
            message: setTab
        });
    }

    function postit() {
        console.log('in GetAllContainers postit');

        var token = $cookieStore.get('cpm_token');

        $http.get($cookieStore.get('AdminURL') + '/backup/getschedules/' + window.containerid + "." + token).
        success(function(data, status, headers, config) {
            $scope.results = data;
        }).
        error(function(data, status, headers, config) {
            console.log('error:UpdateContainerPage.http.get');
        });
    }

    var init = function() {
        console.log('GetAllContainers init called');
        var token = $cookieStore.get('cpm_token');
        $http.get($cookieStore.get('AdminURL') + '/backup/getschedules/' + window.containerid + "." + token).
        success(function(data, status, headers, config) {
            $scope.results = data;
            console.log('containers has ' + $scope.results.length);
            if ($scope.results.length > 0) {
                console.log('first container is ' + $scope.results[0].Name);
                $rootScope.$emit('setSchedule', {
                    message: $scope.results[0]
                });
                $scope.activeClass = $scope.results[0].ID;
            } else {
                $rootScope.$emit('noSchedule', {
                    message: 'hi'
                });
	    }
        }).
        error(function(data, status, headers, config) {
            console.log('error:GetAllContainers:init');
        });
    };

    $rootScope.$on('noScheduleTarget', function(event, args) {
        console.log('no schedule event received');
        $scope.results = [];
    });

    $rootScope.$on('updateSchedulePageTarget', function(event, args) {
        console.log('updating list of schedules ' + args.message);
        $scope.message = args.message;
        $scope.entryID = $scope.message.ID;
        postit(args.message);
    });

    $rootScope.$on('createScheduleTarget', function(event, args) {
        console.log('GetAllController createScheduleTarget recv ' + args.message);
  init();
    });
    $rootScope.$on('deleteScheduleTarget', function(event, args) {
        console.log('GetAllController deleteScheduleTarget received ');
        init();
    });
    if ($cookieStore.get('AdminURL')) {
        init();
    } else {
        alert('CPM AdminURL setting is NOT defined, please update on the Settings page before using CPM');
    }
});

var CreateScheduleCtrl = function($rootScope, $scope, $http, $modalInstance, $cookies, $cookieStore, value) {

    $scope.alerts = [];
    $scope.servers = null;
    $scope.ProfileName = 'pg_basebackup';
    $scope.Name = 'somename';
	$scope.value = value;
    $scope.isLoading = false;

    var token = $cookieStore.get('cpm_token');

    $http.get($cookieStore.get('AdminURL') + '/servers/' + token).
    success(function(data, status, headers, config) {
        $scope.servers = data;
        console.log('got servers len=' + data.length);
        $scope.myServer = $scope.servers[0];
    }).error(function(data, status, headers, config) {
        console.log('error in fetch of servers');
    });

    $scope.cancel = function() {
        $modalInstance.close();
    };

 	$scope.ok = function() {

		$http.post($cookieStore.get('AdminURL') + '/backup/addschedule', {
		    'Token': token,
		    'ServerID': this.myServer.ID,
		    'ContainerName': $scope.value.Name,
		    'ProfileName': this.ProfileName,
		    'Name': this.Name
		}).success(function(data, status, headers, config) {
		    $scope.results = data;
		    $rootScope.$emit('createScheduleEvent', {
			message: ""
		    });
		    $modalInstance.close();
		}).error(function(data, status, headers, config) {
		    console.log('error in Create Schedule Ctrl');
		    $scope.alerts = [{
			type: 'danger',
			msg: data.Error
		    }];
		});

	};
};

var DeleteScheduleCtrl = function($rootScope, $scope, $http, $modalInstance, $cookies, $cookieStore, value) {

    $scope.alerts = [];
	$scope.value = value;
    $scope.isLoading = false;

    var token = $cookieStore.get('cpm_token');

    $scope.cancel = function() {
        $modalInstance.close();
    };

 	$scope.ok = function() {
    		$http.get($cookieStore.get('AdminURL') + '/backup/deleteschedule/' + value.ID + "." + token).
    		success(function(data, status, headers, config) {
	    		console.log('success in delete of schedule ' + value.ID);
		    $rootScope.$emit('deleteScheduleEvent', {
			message: ""
		    });
		    $modalInstance.close();
		}).error(function(data, status, headers, config) {
		    console.log('error in Delete Schedule Ctrl');
		    $scope.alerts = [{
			type: 'danger',
			msg: data.Error
		    }];
		});

	};
};


var UpdateScheduleCtrl = function($rootScope, $scope, $http, $modalInstance, $cookies, $cookieStore, value) {

    $scope.alerts = [];
    $scope.servers = null;
    $scope.Name = 'somename';
	$scope.value = value;
    $scope.isLoading = false;

    var token = $cookieStore.get('cpm_token');

    $http.get($cookieStore.get('AdminURL') + '/servers/' + token).
    success(function(data, status, headers, config) {
        $scope.servers = data;
        console.log('got servers len=' + data.length);
        $scope.myServer = $scope.servers[0];
    }).error(function(data, status, headers, config) {
        console.log('error in fetch of servers');
    });

    $scope.cancel = function() {
        $modalInstance.close();
    };

 	$scope.ok = function() {

		var xmin = $scope.value.Minutes;
		var xhour = $scope.value.Hours;
		var xdom = $scope.value.DayOfMonth;
		var xmon = $scope.value.Month;
		var xdow = $scope.value.DayOfWeek;

		if ($scope.value.Minutes == '') {
			xmin = '*';
		}
		if ($scope.value.Hours == '') {
			xhour = '*';
		}
		if ($scope.value.DayOfMonth == '') {
			xdom = '*';
		}
		if ($scope.value.Month == '') {
			xmon = '*';
		}
		if ($scope.value.DayOfWeek == '') {
			xdow = '*';
		}

		$http.post($cookieStore.get('AdminURL') + '/backup/updateschedule', {
		    'Token': token,
		    'ID': $scope.value.ID,
		    'ServerID': this.myServer.ID,
		    'Enabled': $scope.value.Enabled,
		    'Minutes': xmin,
		    'Hours': xhour,
		    'DayOfMonth': xdom,
		    'Month': xmon,
		    'DayOfWeek': xdow,
		    'Name': $scope.value.Name
		}).success(function(data, status, headers, config) {
		    $scope.results = data;
		    $rootScope.$emit('updateSchedulePage', {
			message: ""
		    });
		    $modalInstance.close();
		}).error(function(data, status, headers, config) {
		    console.log('error in update Schedule Ctrl');
		    $scope.alerts = [{
			type: 'danger',
			msg: data.Error
		    }];
		});

	};
};

var BackupNowCtrl = function($rootScope, $scope, $http, $modalInstance, $cookies, $cookieStore, value) {

    $scope.alerts = [];
    $scope.servers = null;
    $scope.ProfileName = 'pg_basebackup';
    $scope.currentSchedule = value;
    $scope.isLoading = false;

    var token = $cookieStore.get('cpm_token');

    $http.get($cookieStore.get('AdminURL') + '/servers/' + token).
    success(function(data, status, headers, config) {
        $scope.servers = data;
        console.log('got servers len=' + data.length);
        $scope.myServer = $scope.servers[0];
    }).error(function(data, status, headers, config) {
        console.log('error in fetch of servers');
    });

    $scope.cancel = function() {
        $modalInstance.close();
    };

 	$scope.ok = function() {

		$http.post($cookieStore.get('AdminURL') + '/backup/now', {
		    'Token': token,
		    'ServerID': this.myServer.ID,
		    'ProfileName': $scope.ProfileName,
		    'ScheduleID': $scope.currentSchedule.ID
		}).success(function(data, status, headers, config) {
		    $scope.results = data;
		    $rootScope.$emit('createScheduleEvent', {
			message: ""
		    });
		    $modalInstance.close();
		}).error(function(data, status, headers, config) {
		    console.log('error in BackupNowCtrl');
		    $scope.alerts = [{
			type: 'danger',
			msg: data.Error
		    }];
		});

	};
};


})();
