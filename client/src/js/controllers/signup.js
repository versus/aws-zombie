angular.module('app')
    .controller('sign-up', [
        '$rootScope',
        '$scope',
        'app.login',
        '$location',
        'app.user',
        function($rootScope, $scope, auth, $location, user) {

            if (window.localStorage.getItem("token")) {
                $location.path("/contacts");
                return;
            }

            $scope.signup = {};
            $scope.submit = function(){
                if (!$scope.signup.phone || !$scope.signup.password) {
                    $scope.error = 'All fields are required.';
                    return
                }

                if (!$scope.signup.phone.match(/^[0-9]{12}$/)) {
                    $scope.error = 'Wrong phone number.';
                    return
                }

                auth.login({
                    phone: $scope.signup.phone,
                    password: $scope.signup.password}, function(data){
                    if (data.hasOwnProperty("token")) {
                        window.localStorage.setItem("token", data.token);
                        window.localStorage.setItem("phone", $scope.signup.phone);


                        if (window.localStorage.getItem("pushId")) {
                            user.update({phone: $scope.signup.phone}, {
                                push_id: window.localStorage.getItem("pushId")
                            });
                        }

                        $location.path('/contacts');
                    }

                }, function(err){
                    if (err.hasOwnProperty('data')) {
                        if (err.data.message == "Error: Wrong password") {
                            $scope.error = err.data.message;
                        } else if(err.data.message == "Error: User does not exists") {
                            $rootScope.profile = {
                                phone: $scope.signup.phone,
                                password: $scope.signup.password
                            };

                            $location.path('/profile-create');
                        }
                    }
                });
            };
        }
    ])
    .controller('logout', [
        '$location',
        function($location){
            window.localStorage.removeItem("token");
            window.localStorage.removeItem("phone");
            $location.path("/");
        }
    ])
    .controller('profile-create', [
        '$rootScope',
        '$scope',
        'app.login',
        '$location',
        function($rootScope, $scope, login, $location) {

            if (window.localStorage.getItem("token")) {
                $location.path("/contacts");
                return;
            }

            if (!$rootScope.profile) {
                $location.path('/');
            }

            $scope.profile = {};
            $scope.submit = function(){
                console.log($scope.profile);

                if (!$scope.profile.first_name || !$scope.profile.last_name || !$scope.profile.skills) {
                    $scope.error = 'All fields are required.';
                    return;
                }

                var params = {
                    phone: $rootScope.profile.phone,
                    first_name: $scope.profile.first_name,
                    last_name: $scope.profile.last_name,
                    password: $rootScope.profile.password,
                    specs: $scope.profile.skills
                };

                if (window.localStorage.getItem("pushId")) {
                    params.push_id = window.localStorage.getItem("pushId");
                }

                login.signUp(params, function(){
                    login.login({
                        phone: $rootScope.profile.phone,
                        password: $rootScope.profile.password}, function(data){
                        if (data.hasOwnProperty("token")) {
                            window.localStorage.setItem("token", data.token);
                            window.localStorage.setItem("phone", $rootScope.profile.phone);

                            $location.path('/contacts');
                        }}, function(){

                        $location.path("/signup");
                    });

                    }, function(err){
                }, function() {
                    $location.path("/signup");
                })

            };


        }
    ]);

