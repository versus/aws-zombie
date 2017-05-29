angular.module('app')
    .factory('app.contact', function ($resource, host) {
        return $resource(host, {}, {
            add: {
                method: 'POST',
                url: host + '/user/:phone/contacts/add'
            },
            get: {

            }
        });
    });