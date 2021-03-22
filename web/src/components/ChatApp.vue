<template>
  <div>
    <div class="row row-msg">
      <div class="col s12">
        <div class="card horizontal">
          <div id="chat-messages" class="card-content" v-html="chatContent">
          </div>
        </div>
      </div>
    </div>


    <div class="row row-write" v-if="joined">
      <div class="input-field col s8">
        <input type="text" v-model="newMsg" @keyup.enter="send">
      </div>
      <div class="input-field col s4">
        <button class="waves-effect waves-light btn" @click="send">
          <i class="material-icons right">chat</i>
          Send
        </button>
      </div>
    </div>

    <div class="row" v-if="!joined">
      <div class="input-field col s8">
        <span v-if="join_error">{{join_error}}</span>
        <input type="text" v-model.trim="username" placeholder="Username">
      </div>
      <div class="input-field col s4">
        <button class="waves-effect waves-light btn" @click="join()">
          <i class="material-icons right">done</i>
          Join
        </button>
      </div>
    </div>


    <br>
    <div v-if="joined">
      <gmap-map
          :center="startingCenter"
          :zoom="3"
          style="width:100%;  height: 400px;"
          :options="{
streetViewControl: false,
fullscreenControl: false,
}"
      >
        <GmapCluster
            :position="center" :clickable="true" :animation="2"
        >
        <gmap-marker v-for="(item, key) in allRooms"
                     :key="key"
                     :label="getJoinText(item.Name)"
            :position="getRoomGeo(item)"
            :draggable="false"
                     :clickable="true"
                     @click="toggleInfo(item, key)"
        ></gmap-marker>
        </GmapCluster>
      </gmap-map>
    </div>

  </div>
</template>

<script>
/*
TODO

- there should not be too many markers on the map
- the markers should show a join action and the number of active users
- add a link for creating a new room:
  1. it removes all other marker on map
  2. it places a new marker which is draggable by the user
  3. it shows a field for the room name
  4. it shows a Save button
  5. on submit, validate the name
  6. after validation, send room to backend and open the chat
  7. open the regular map
 */


export default {
  components: {},
  name: 'ChatApp',
  data: function() {
    return {
      ws: null, // Our websocket
      newMsg: '', // Holds new messages to be sent to the server
      chatContent: '', // A running list of chat messages displayed on the screen
      username: null, // Our username
      joined: false, // True if email and username have been filled in
      startingCenter: { lat: 45.508, lng: -73.587 },
      join_error: '',
      allRooms: {},
      myRoom: null,
    }
  },
  mounted() {
  },
  created: function() {
    var self = this;
    this.ws = new WebSocket('ws://' + window.location.host + '/ws');
    this.ws.addEventListener('message', function(e) {
      var msg = JSON.parse(e.data);
      if (msg.type === 'message') {
        self.chatContent += '<div class="chip">'
            + msg.username
            + '</div>'
            + msg.message + '<br/>';
        var element = document.getElementById('chat-messages');
        element.scrollTop = element.scrollHeight; // Auto scroll to the bottom
      }
      else if (msg.type === 'room_list') {
        self.allRooms = msg.room_list;
      }
      else if (msg.type === 'room_joined') {
        self.myRoom = msg.room
      }
      else {
        console.log('Unknown broadcast type');
      }
    });
  },
  methods: {
    getJoinText(roomName) {
      return 'Join ' + roomName
    },
    toggleInfo: function (marker) {
      this.joinRoom(marker)
    },
    joinRoom: function (room) {
      console.log('Implement room join', room);
    },
    getRoomGeo: function (room) {
      return {
        lat: room.Location.Lat,
        lng: room.Location.Lng
      }
    },
    stripHtml: function (html){
      var temporalDivElement = document.createElement("div");
      temporalDivElement.innerHTML = html;
      return temporalDivElement.textContent || temporalDivElement.innerText || "";
    },
    send: function () {
      if (this.newMsg !== '') {
        this.ws.send(
            JSON.stringify({
                  message: this.stripHtml(this.newMsg),
                  type: 'message',
                }
            ));
        this.newMsg = '';
      }
    },
    join: function () {
      if (!this.username) {
        this.join_error = 'Please enter a username';
        return
      }
      this.username = this.stripHtml(this.username);
      this.joined = true;
      this.ws.send(
          JSON.stringify({
                username: this.username,
                type: 'register'
              }
          ));
    },
  },
}
</script>
