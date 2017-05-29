angular.module('app')
    .controller('main', [
        '$rootScope',
        '$scope',
        '$routeParams',
        'app.user',
        'app.message',
        'app.contact',
        function($rootScope, $scope, $routeParams, user, message, contact) {
            $rootScope.showNav = true;

            $scope.messages = [];

            user.get({phone: $routeParams.phone}, function(cData){
                console.log(cData);
                contact.add({phone: window.localStorage.getItem("phone")}, {
                    contact_phone: $routeParams.phone
                });


                message.getAll({
                    phone_to: $routeParams.phone,
                    phone_from: window.localStorage.getItem("phone")
                }, function(data1){
                    for (var i=0;i<data1.length;i++) {
                        $scope.loadedMessages.push(data1[i]);
                    }

                    message.getAll({
                        phone_to: window.localStorage.getItem("phone"),
                        phone_from: $routeParams.phone
                    }, function(data2){
                        for (var i=0;i<data2.length;i++) {
                            $scope.loadedMessages.push(data2[i]);
                        }

                        for (var i=0;i<$scope.loadedMessages.length;i++) {
                            $scope.loadedMessages[i].message = $scope.loadedMessages[i].text;
                            if ($scope.loadedMessages[i].from == window.localStorage.getItem("phone")) {
                                $scope.loadedMessages[i].cl = "you";
                                $scope.loadedMessages[i].sender = {
                                    first_name: "You",
                                    phone: window.localStorage.getItem("phone"),
                                };
                            } else {
                                $scope.loadedMessages[i].sender = {
                                    first_name: cData.first_name,
                                    last_name: cData.last_name
                                };
                            }
                        }

                        console.log($scope.loadedMessages);
                        $scope.loadedMessages.sort(function(a, b){
                            if (parseInt(a.timestamp) < parseInt(b.timestamp)) {
                                return -1;
                            }
                            if (parseInt(a.timestamp) > parseInt(b.timestamp)) {
                                return 1;
                            }

                            return 0;
                        });
                        $scope.messages = $scope.loadedMessages;
                    });
                });
            });



            //
            $scope.message = '';
            $scope.loadedMessages = [];




            $scope.send = function(){
                console.log($scope.message);
                message.send({
                    from_phone: window.localStorage.getItem("phone"),
                    to: $routeParams.phone,
                    text: $scope.message
                });
                $scope.messages.push({
                    message: $scope.message,
                    cl: "you",
                    sender: {
                        first_name: "You",
                        phone: $scope.phone
                    }
                });
                $scope.message = '';
                window.scrollTo(0,document.body.scrollHeight);
            };

            // console.log(globalSocket);
            $scope.phone = window.localStorage.getItem("phone");
            document.socket.on('message', function(data){
                $scope.$apply(function () {
                    //$scope.message = "Timeout called!";
                    $scope.messages.push(JSON.parse(data));
                });

            });
        }
    ]);