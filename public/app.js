new Vue({
    el: '#app',

    data: {
        ws: null, // Our websocket
        newMsg: '', // Holds new messages to be sent to the server
        chatContent: '', // A running list of chat messages displayed on the screen
        room: "@room0", // Room id
        username: "b", // Our username
        joined: false // True if room and username have been filled in
    },

    created: function() {
        var self = this;
        this.ws = new WebSocket('ws://' + window.location.host + '/ws');
        this.ws.addEventListener('message', function(e) {
            var msg = JSON.parse(e.data);
            self.chatContent += '<div class="chip">'
                    + '<img src="' + self.avatar(msg) + '">' // Avatar
                    //+ '<img src="2.png">'
                    + msg.username
                + '</div>'
                + emojione.toImage(msg.message) + '<br/>'; // Parse emojis
            //Delay the scroller to refresh 2022/11
            setTimeout(function() {
                var element = document.getElementById('chat-messages');
                element.scrollTop = element.scrollHeight; // Auto scroll to the bottom
            }, 60);    
        });
    },

    methods: {
        scroll: function() {
            var element = document.getElementById('chat-messages');
            element.scrollTop = element.scrollHeight + element.clientHeight; // Auto scroll to the bottom
        },
        send: function () {
            //if (this.newMsg != 'zxswertgvb') {
                this.ws.send(
                    JSON.stringify({
                        room: this.room,
                        username: this.username,
                        avatar: '',  //server generating the avatar
                        message: $('<p>').html(this.newMsg).text() // Strip out html
                    }
                ));
                this.newMsg = ''; // Reset newMsg
            //}
        },

        join: function () {
            if (!this.room) {
                Materialize.toast('You must enter a valid room', 2000);
                return
            }
            if (!this.username) {
                Materialize.toast('You must choose a username', 2000);
                return
            }
            this.room = $('<p>').html(this.room).text();
            this.username = $('<p>').html(this.username).text();
            document.getElementById('room-title').text = this.room;
            this.joined = true;
            this.send();
        },

        avatar: function(msg) {
            if (msg.avatar != '') {
                avatar = msg.avatar;
            } else {
                avatar = ((username.charCodeAt(0) % 5) + 1).toString() + '.png';
            }
            return avatar;
        },

        gravatarURL: function(email) {
            return 'http://www.gravatar.com/avatar/' + CryptoJS.MD5(email);
        }
    }
});
