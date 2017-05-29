angular.module('app', [
    'ngRoute',
    'ngResource',
    'ui.bootstrap',
    'angular-loading-bar'
    // 'ngAnimate'
])
    .run(['$rootScope', '$http', function ($rootScope, $http) {

        $rootScope.hasToken = function(){
            return window.localStorage.getItem('token') != null;
        };

        if (window.localStorage.getItem('token') != null) {
            $rootScope.token = window.localStorage.getItem('token');
            $http.defaults.headers.common["Authorization"] = "Bearer " + $rootScope.token;
        }

        $rootScope.$on('$routeChangeSuccess', function (event, current, previous) {
            if(current.$$route.hasOwnProperty('token') && current.$$route.token && !$rootScope.hasToken()) {
                $location.path('/');
            }

            var navMain = $("#bs-example-navbar-collapse-1");
            navMain.collapse('hide');
        });

    }])
    .constant('host', 'https://api.eu.zombiegram.tk')
    // .constant('shost', 'socket.zombiegram.tk')
    .config([
        '$routeProvider', '$httpProvider', '$locationProvider', function ($routeProvider, $httpProvider, $locationProvider) {
            $routeProvider
                .when('/', {
                    templateUrl: './templates/index/index.html',
                    title: 'AWS Zombie',
                    controller: 'index'
                })
                .when('/signup', {
                    templateUrl: './templates/index/signup.html',
                    title: '',
                    controller: 'sign-up'
                })
                .when('/profile-create', {
                    templateUrl: './templates/index/profile-create.html',
                    title: '',
                    controller: 'profile-create'
                })
                .when('/main/:phone', {
                    templateUrl: './templates/messenger/main.html',
                    title: '',
                    controller: 'main'
                })
                .when('/chats', {
                    templateUrl: './templates/messenger/chats.html',
                    title: '',
                    controller: 'chats'
                })
                .when('/contacts', {
                    templateUrl: './templates/messenger/contacts.html',
                    title: '',
                    controller: 'contacts'
                })
                .when('/logout', {
                    templateUrl: './templates/index/logout.html',
                    title: '',
                    controller: 'logout'
                })
                .otherwise({
                    redirectTo: '/'
                });


            // $locationProvider.html5Mode(true);
            // $locationProvider.html5Mode({
            //     enabled: true,
            //     rewriteLinks: false
            // });

            // console.log($httpProvider);

            // var appendTransform = function(defaults, transform) {
            //
            //     // We can't guarantee that the default transformation is an array
            //     defaults = angular.isArray(defaults) ? defaults : [defaults];
            //
            //     // Append the new transformation to the defaults
            //     return defaults.concat(transform);
            // };
            //
            // $httpProvider.defaults.transformRequest = appendTransform($httpProvider.defaults.transformRequest, function(data) {
            //     //do whatever you want
            //     if (data) {
            //         data = JSON.parse(data);
            //
            //         data = JSON.stringify({
            //             body: JSON.stringify(data),
            //             httpMethod: data.httpMethod
            //         })
            //     }
            //
            //     return data;
            // });

        }]);