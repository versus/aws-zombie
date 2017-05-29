angular.module('app')
    .factory('app.message', function ($resource, host) {
        return $resource(host, {}, {
            send: {
                method: 'POST',
                url: host + '/message/send'
            },
            all: {
                method: 'GET',
                url: host + '/user/list',
                isArray: true
            },
            getAll: {
                method: 'POST',
                url: host + '/message/list',
                isArray: true
            }
        });
    });