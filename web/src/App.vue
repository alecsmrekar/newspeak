<template>
  <div id="app">
    <ChatApp />
  </div>
</template>

<script>
import ChatApp from "@/components/ChatApp";
import Vue from "vue";
import * as VueGoogleMaps from "vue2-google-maps";

Vue.use(VueGoogleMaps, {
  load: {
    key: "AIzaSyA1YSCBIT7bedXGz11kEPSnhz_XGHQrjs8",
    libraries: "places" // necessary for places input
  }
});

export default {
  name: 'App',
  components: {
    ChatApp
  }
}
</script>