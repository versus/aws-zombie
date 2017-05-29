angular.module('app')
    .factory('app.user', function ($resource, host) {
        return $resource(host, {}, {
            get: {
                method: 'GET',
                url: host + '/user/:phone'
            },
            all: {
                method: 'GET',
                url: host + '/user/list',
                isArray: true
            },
            update: {
                method: 'PUT',
                url: host + '/user/:phone'
            }
        });
    });