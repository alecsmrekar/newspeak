(function(e){function t(t){for(var o,s,i=t[0],l=t[1],c=t[2],m=0,p=[];m<i.length;m++)s=i[m],Object.prototype.hasOwnProperty.call(a,s)&&a[s]&&p.push(a[s][0]),a[s]=0;for(o in l)Object.prototype.hasOwnProperty.call(l,o)&&(e[o]=l[o]);u&&u(t);while(p.length)p.shift()();return r.push.apply(r,c||[]),n()}function n(){for(var e,t=0;t<r.length;t++){for(var n=r[t],o=!0,i=1;i<n.length;i++){var l=n[i];0!==a[l]&&(o=!1)}o&&(r.splice(t--,1),e=s(s.s=n[0]))}return e}var o={},a={app:0},r=[];function s(t){if(o[t])return o[t].exports;var n=o[t]={i:t,l:!1,exports:{}};return e[t].call(n.exports,n,n.exports,s),n.l=!0,n.exports}s.m=e,s.c=o,s.d=function(e,t,n){s.o(e,t)||Object.defineProperty(e,t,{enumerable:!0,get:n})},s.r=function(e){"undefined"!==typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(e,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(e,"__esModule",{value:!0})},s.t=function(e,t){if(1&t&&(e=s(e)),8&t)return e;if(4&t&&"object"===typeof e&&e&&e.__esModule)return e;var n=Object.create(null);if(s.r(n),Object.defineProperty(n,"default",{enumerable:!0,value:e}),2&t&&"string"!=typeof e)for(var o in e)s.d(n,o,function(t){return e[t]}.bind(null,o));return n},s.n=function(e){var t=e&&e.__esModule?function(){return e["default"]}:function(){return e};return s.d(t,"a",t),t},s.o=function(e,t){return Object.prototype.hasOwnProperty.call(e,t)},s.p="/";var i=window["webpackJsonp"]=window["webpackJsonp"]||[],l=i.push.bind(i);i.push=t,i=i.slice();for(var c=0;c<i.length;c++)t(i[c]);var u=l;r.push([0,"chunk-vendors"]),n()})({0:function(e,t,n){e.exports=n("56d7")},"56d7":function(e,t,n){"use strict";n.r(t);n("e260"),n("e6cf"),n("cca6"),n("a79d");var o=n("2b0e"),a=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("div",{attrs:{id:"app"}},[n("ChatApp")],1)},r=[],s=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("div",[n("div",{staticClass:"row row-msg"},[n("div",{staticClass:"col s12"},[n("div",{staticClass:"card horizontal"},[n("div",{staticClass:"card-content",attrs:{id:"chat-messages"},domProps:{innerHTML:e._s(e.chatContent)}})])])]),""!=e.myRoom?n("div",[e._v("My Room: "+e._s(e.myRoom))]):e._e(),""!=e.myRoom?n("div",[n("button",{on:{click:e.leave}},[e._v("LEAVE NOW")])]):e._e(),""==e.myRoom&&e.show_map&&e.joined?n("div",[n("button",{on:{click:e.createNewStart}},[e._v("CREATE NEW")])]):e._e(),1==e.newMarkerState?n("div",[n("input",{directives:[{name:"model",rawName:"v-model",value:e.newRoomName,expression:"newRoomName"}],domProps:{value:e.newRoomName},on:{input:function(t){t.target.composing||(e.newRoomName=t.target.value)}}})]):e._e(),1==e.newMarkerState?n("div",[n("button",{on:{click:e.createNewDone}},[e._v("Create")])]):e._e(),e.joined?n("div",{staticClass:"row row-write"},[n("div",{staticClass:"chat-text"},[n("div",{staticClass:"input-field col s8"},[n("input",{directives:[{name:"model",rawName:"v-model",value:e.newMsg,expression:"newMsg"}],attrs:{type:"text"},domProps:{value:e.newMsg},on:{keyup:function(t){return!t.type.indexOf("key")&&e._k(t.keyCode,"enter",13,t.key,"Enter")?null:e.send(t)},input:function(t){t.target.composing||(e.newMsg=t.target.value)}}})]),n("div",{staticClass:"input-field col s4"},[n("button",{staticClass:"waves-effect waves-light btn",on:{click:e.send}},[n("i",{staticClass:"material-icons right"},[e._v("chat")]),e._v(" Send ")])])]),e.joined&&""!=e.myRoom?n("div",{staticClass:"users",staticStyle:{clear:"left"}},[e._v(" Room members: "),n("ul",e._l(e.chat_users,(function(t){return n("li",{key:t.id},[e._v(e._s(t))])})),0)]):e._e()]):e._e(),e.joined?e._e():n("div",{staticClass:"row"},[n("div",{staticClass:"input-field col s8"},[e.join_error?n("span",[e._v(e._s(e.join_error))]):e._e(),n("input",{directives:[{name:"model",rawName:"v-model.trim",value:e.username,expression:"username",modifiers:{trim:!0}}],attrs:{type:"text",placeholder:"Username"},domProps:{value:e.username},on:{input:function(t){t.target.composing||(e.username=t.target.value.trim())},blur:function(t){return e.$forceUpdate()}}})]),n("div",{staticClass:"input-field col s4"},[n("button",{staticClass:"waves-effect waves-light btn",on:{click:function(t){return e.join()}}},[n("i",{staticClass:"material-icons right"},[e._v("done")]),e._v(" Join ")])])]),n("br"),e.joined&&e.show_map?n("div",[n("gmap-map",{staticStyle:{width:"100%",height:"400px"},attrs:{center:e.mapCenterStart,zoom:3,options:{streetViewControl:!1,fullscreenControl:!1}},on:{center_changed:e.dragMap}},[n("GmapCluster",{attrs:{position:e.center,clickable:!0,animation:2}},[e._l(e.allRooms,(function(t,o){return n("gmap-marker",{key:o,attrs:{label:e.getJoinText(t.Name),position:e.getRoomGeo(t),draggable:!1,clickable:!0},on:{click:function(n){return e.toggleInfo(t,o)}}})})),1==e.newMarkerState?n("div",[n("gmap-marker",{key:e.newMarkerId,staticClass:"new-room-marker",attrs:{label:e.newMarkerLabel,position:e.newMarkerPos,draggable:!0,clickable:!0},on:{drag:e.updateNewMarkerCoords}})],1):e._e()],2)],1)],1):e._e()])},i=[],l={components:{},name:"ChatApp",data:function(){return{ws:null,newMsg:"",chatContent:"",username:null,joined:!1,mapCenter:{lat:45.508,lng:-73.587},mapCenterStart:{lat:45.508,lng:-73.587},join_error:"",allRooms:{},newMarkerId:-1,newMarkerLabel:"Set Room Location",newMarkerPos:{lat:0,lng:0},newMarkerState:0,myRoom:"",show_map:!0,newRoomName:"",chat_users:[]}},mounted:function(){},created:function(){var e=this;this.ws=new WebSocket("ws://"+window.location.host+"/ws"),this.ws.addEventListener("message",(function(t){var n=JSON.parse(t.data);if("message"===n.type){e.chatContent+='<div class="chip">'+n.username+"</div>"+n.message+"<br/>";var o=document.getElementById("chat-messages");o.scrollTop=o.scrollHeight}else"room_list"===n.type?e.allRooms=n.room_list:"room_update"===n.type?(e.chat_users=n.users,console.log(e.chat_users),n.message&&(e.chatContent+=n.message+"<br/>")):console.log("Unknown broadcast type")}))},methods:{getJoinText:function(e){return"Join "+e},toggleInfo:function(e){this.joinRoom(e)},joinRoom:function(e){this.ws.send(JSON.stringify({room_id:e.ID,type:"join_room"})),this.show_map=!1,this.myRoom=e.Name},getRoomGeo:function(e){return{lat:e.Location.Lat,lng:e.Location.Lng}},stripHtml:function(e){var t=document.createElement("div");return t.innerHTML=e,t.textContent||t.innerText||""},send:function(){""!==this.newMsg&&(this.ws.send(JSON.stringify({message:this.stripHtml(this.newMsg),type:"message"})),this.newMsg="")},leave:function(){this.myRoom="",this.show_map=!0,this.chatContent="",this.ws.send(JSON.stringify({type:"leave_room"}))},join:function(){this.username?(this.username=this.stripHtml(this.username),this.joined=!0,this.ws.send(JSON.stringify({username:this.username,type:"register"}))):this.join_error="Please enter a username"},createNewStart:function(){this.newMarkerPos=this.mapCenter,this.newMarkerState=1,this.allRooms={}},createNewDone:function(){this.newMarkerState=0,this.ws.send(JSON.stringify({room_name:this.newRoomName,lat:this.newMarkerPos.lat,lng:this.newMarkerPos.lng,type:"create_room"})),this.show_map=!1,this.myRoom=this.newRoomName},updateNewMarkerCoords:function(e){this.newMarkerPos={lat:e.latLng.lat(),lng:e.latLng.lng()}},dragMap:function(e){this.mapCenter={lat:e.lat(),lng:e.lng()}}}},c=l,u=n("2877"),m=Object(u["a"])(c,s,i,!1,null,null,null),p=m.exports,d=n("ae66"),h=n.n(d),f=n("755e"),g="AIzaSyDQ6FHKm4ZWkcQDrBl3m7B7PSTlui4Fs5c";o["a"].component("GmapCluster",h.a),o["a"].use(h.a),o["a"].use(f,{load:{key:g}});var w={name:"App",components:{ChatApp:p}},v=w,_=Object(u["a"])(v,a,r,!1,null,null,null),y=_.exports;o["a"].config.productionTip=!1,new o["a"]({render:function(e){return e(y)}}).$mount("#app")}});
//# sourceMappingURL=app.f9aad359.js.map