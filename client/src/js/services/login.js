angular.module('app')
    .factory('app.login', function ($resource, host) {
        return $resource(host, {}, {
            signUp: {
                method: 'POST',
                url: host + '/user'
            },
            login: {
                method: 'POST',
                url: host + '/auth'
            },
        });
    });