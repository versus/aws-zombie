angular.module('app')
    .controller('index', [
        '$location',
        '$rootScope',
        function($location, $rootScope) {
            if (window.localStorage.getItem("token")) {
                $location.path("/contacts");
                return;
            }

            $rootScope.showNav = false;
        }
    ]);