<template>
  <div id="app">
    <ChatAppV2 />
  </div>
</template>

<script>
import ChatApp from "@/components/ChatApp";
import ChatAppV2 from "@/components/ChatAppV2";
import Vue from "vue";
import * as VueGoogleMaps from "vue2-google-maps";
import { maps_key } from "./keys";

Vue.use(VueGoogleMaps, {
  load: {
    key: maps_key,
  }
});


export default {
  name: 'App',
  components: {
    ChatApp, ChatAppV2
  }
}
</script>