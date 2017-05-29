/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
// var globalSocket;

var app = {
    // Application Constructor
    initialize: function() {
        document.addEventListener('deviceready', this.onDeviceReady.bind(this), false);



        document.socket = io.connect("http://socket.zombiegram.tk:8080");

        document.socket.on('foo', function() {
            console.log("fo!");
            if (window.localStorage.getItem("phone")) {
                document.socket.emit('phone', window.localStorage.getItem("phone"));
            }
        });

        document.socket.on('message', function(data){
            //console.log(data);
            data = JSON.parse(data);
            if (document.location.hash.search(data.sender.phone) == -1) {
                document.location.hash = '#!/main/'+data.sender.phone;
            }

            if (!document.appActive) {
                cordova.plugins.notification.local.schedule({
                    // id: 1,
                    text: 'Chat message!'
                    // sound: isAndroid ? 'file://sound.mp3' : 'file://beep.caf',
                    // every: 'day',
                    // firstAt: next_monday,
                    // data: { key:'value' }
                })
            }
            //if (document.location.hash.)
        });
    },

    // deviceready Event Handler
    //
    // Bind any cordova events here. Common events are:
    // 'pause', 'resume', etc.
    onDeviceReady: function() {
        this.receivedEvent('deviceready');

        document.appActive = true;

        document.addEventListener("pause", function(){
            document.appActive = false;
        }, false);

        document.addEventListener("resume", function(){
            document.appActive = true;
        }, false);
    },

    // Update DOM on a Received Event
    receivedEvent: function(id) {
        /*var parentElement = document.getElementById(id);
        var listeningElement = parentElement.querySelector('.listening');
        var receivedElement = parentElement.querySelector('.received');

        listeningElement.setAttribute('style', 'display:none;');
        receivedElement.setAttribute('style', 'display:block;');*/

        console.log('Received Event: ' + id);






        // console.log(typeof(PushNotification));
        var push = PushNotification.init({
            "android": {
                "senderID": "944405234963"
            }
            // "ios": {"alert": "true", "badge": "true", "sound": "true"},
            // "windows": {}
        });

        push.on('registration', function(data) {
            console.log("registration event");
            //here is your registration id
            if (data.hasOwnProperty("registrationId")) {
                window.localStorage.setItem("pushId", data.registrationId)
            }
        });
    }
};

app.initialize();
