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
    <div>
    <gmap-map
        :center="{ lat: 45.508, lng: -73.587 }"
        :zoom="3"
        style="width:100%;  height: 400px;"
        :options="{
          streetViewControl: false,
          fullscreenControl: false,
        }"
    >
      <gmap-marker
        :position="startingCenter"
        :draggable="true"
        @dragend="mapDrag"
        ></gmap-marker>
      <gmap-circle
          :strokeOpacity= "0.8"
          :strokeWeight= "2"
          :fillOpacity= "0.5"
          :center="center"
          :radius="radius"
      ></gmap-circle>
    </gmap-map>
  </div>
    <div class="radius">
      <vue-range-slider @change="radiusChange"
                        style="width: 100%"
                        :min="1000"
                        :max="500000"
          ref="slider" v-model="slider_radius"></vue-range-slider>
    </div>

  </div>
</template>

<script>

import VueRangeSlider from "vue-range-slider";

export default {
  components: { VueRangeSlider },
  name: 'ChatApp',
  data: function() {
    return {
      ws: null, // Our websocket
      newMsg: '', // Holds new messages to be sent to the server
      chatContent: '', // A running list of chat messages displayed on the screen
      username: null, // Our username
      joined: false, // True if email and username have been filled in
      center: { lat: 45.508, lng: -73.587 },
      startingCenter: { lat: 45.508, lng: -73.587 },
      radius: 200000,
      max_radius:500000,
      slider_radius: 200000,
      join_error: '',
    }
  },
  mounted() {
  },
  created: function() {
    var self = this;
    this.ws = new WebSocket('ws://' + window.location.host + '/ws');

    // Handle incoming messages
    this.ws.addEventListener('message', function(e) {
      var msg = JSON.parse(e.data);
      self.chatContent += '<div class="chip">'
          + msg.username
          + '</div>'
          + msg.message + '<br/>';

      var element = document.getElementById('chat-messages');
      element.scrollTop = element.scrollHeight; // Auto scroll to the bottom
    });
  },
  methods: {
    radiusChange: function () {
      this.radius = this.slider_radius;

      // Update the users radius as he changes it
      this.ws.send(
          JSON.stringify({
                radius: this.radius,
                radiusUpdate: true,
              }
          ));
    },
    mapDrag: function (arg) {
      this.center = arg.latLng;

      // Update the users location as he changes it
      this.ws.send(
          JSON.stringify({
                lat: this.center.lat,
                lng: this.center.lng,
                locationUpdate: true,
              }
          ));
    },
    stripHtml: function (html){
      var temporalDivElement = document.createElement("div");
      temporalDivElement.innerHTML = html;
      return temporalDivElement.textContent || temporalDivElement.innerText || "";
    },
    send: function () {
      if (this.newMsg != '') {
        this.ws.send(
            JSON.stringify({
                  username: this.username,
                  message: this.stripHtml(this.newMsg)
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

      // As soon as the user joins, initialize his initial geo data
      this.ws.send(
          JSON.stringify({
                radius: this.radius,
                radiusUpdate: true,
                lat: this.center.lat,
                lng: this.center.lng,
                locationUpdate: true,
              }
          ));
    },
  },
}
</script>
