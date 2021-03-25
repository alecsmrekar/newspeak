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
    <div v-if="myRoom != ''">My Room: {{myRoom}}</div>
    <div v-if="myRoom != ''"><button @click="leave">LEAVE NOW</button></div>
    <div v-if="myRoom == '' && show_map && joined"><button @click="createNewStart">CREATE NEW</button></div>
    <div v-if="newMarkerState == 1"><input v-model="newRoomName"></div>
    <div v-if="newMarkerState == 1"><button @click="createNewDone">Create</button></div>


    <div class="row row-write" v-if="joined">
      <div class="chat-text">
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
      <div class="users" style="clear: left" v-if="joined && myRoom != ''">
        Room members:
        <ul>
          <li v-for="(usr) in chat_users" :key="usr.id">{{usr}}</li>
        </ul>
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
    <div v-if="joined && show_map">
      <gmap-map
          :center="mapCenterStart"
          @center_changed="dragMap"
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
                     :label="getJoinText(item)"
            :position="getRoomGeo(item)"
            :draggable="false"
                     :clickable="true"
                     @click="clickedJoinRoom(item, key)"
        ></gmap-marker>
          <div v-if="newMarkerState==1">
          <gmap-marker class="new-room-marker"
                       :key="newMarkerId"
                       :label="newMarkerLabel"
                       :position="newMarkerPos"
                       :draggable="true"
                       :clickable="true"
                       @drag="updateNewMarkerCoords"
          ></gmap-marker></div>
        </GmapCluster>
      </gmap-map>
    </div>

  </div>
</template>

<script>

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
      mapCenter: { lat: 45.508, lng: -73.587 },
      mapCenterStart: { lat: 45.508, lng: -73.587 },
      join_error: '',
      allRooms: {},
      newMarkerId: -1,
      newMarkerLabel: 'Set Room Location',
      newMarkerPos: {lat: 0, lng: 0},
      newMarkerState: 0,
      myRoom: '',
      show_map: true,
      newRoomName: '',
      chat_users: []
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
      else if (msg.type === 'room_update') {
        self.chat_users = msg.users;
        console.log(self.chat_users);
        if (msg.message) {
          self.chatContent += msg.message + '<br/>';
        }
      }
      else {
        console.log('Unknown broadcast type');
      }
    });
  },
  methods: {
    getJoinText(room) {
      if (room.NameSuffix !== undefined) {
        return room.Name + ' (' + room.NameSuffix + ')'
      }
      return room.Name
    },
    clickedJoinRoom: function (marker) {
      this.joinRoom(marker)
    },
    joinRoom: function (room) {
      this.ws.send(
          JSON.stringify({
                room_id: room.ID,
                type: 'join_room'
              }
          ));
      this.show_map = false;
      this.myRoom = room.Name;
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
    leave: function () {
      this.myRoom = '';
      this.show_map = true;
      this.chatContent = '';
      this.ws.send(
          JSON.stringify({
                type: 'leave_room',
              }
          ));
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
    createNewStart: function () {
      this.newMarkerPos = this.mapCenter;
      this.newMarkerState = 1;
      this.allRooms = {};
    },
    createNewDone: function () {
      this.newMarkerState = 0;
      this.ws.send(
          JSON.stringify({
                room_name: this.newRoomName,
                lat: this.newMarkerPos.lat,
                lng: this.newMarkerPos.lng,
                type: 'create_room'
              }
          ));
      this.show_map = false;
      this.myRoom = this.newRoomName;
    },
    updateNewMarkerCoords: function (location) {
      this.newMarkerPos = {
        lat: location.latLng.lat(),
        lng: location.latLng.lng(),
      };
    },
    dragMap: function (center) {
      this.mapCenter = {
        lat: center.lat(),
        lng: center.lng()
      }
    }
  },
}
</script>
