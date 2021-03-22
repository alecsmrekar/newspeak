<template>
  <div id="app">
    <ChatApp />
  </div>
</template>

<script>
import ChatApp from "@/components/ChatApp";
import Vue from "vue";
import GmapCluster from 'vue2-google-maps/dist/components/cluster';
import * as VueGoogleMaps from "vue2-google-maps";
import { maps_key } from "./keys";

Vue.component('GmapCluster', GmapCluster)
Vue.use(GmapCluster)
Vue.use(VueGoogleMaps, {
  load: {
    key: maps_key,
  }
});


export default {
  name: 'App',
  components: {
    ChatApp
  }
}
</script>