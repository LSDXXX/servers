import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import ElementUI from "element-ui";
import "element-ui/lib/theme-chalk/index.css";
import socketio from "socket.io-client";
// import VueSocketio from "vue-socket.io";
import axios from "axios";

Vue.use(ElementUI);
//Vue.use(
//  new VueSocketio({
//    debug: true,
//    connection: socketio.connect("http://192.168.0.159:2120", {
//      path: "",
//      transports: ["websocket", "xhr-polling", "jsonp-polling"],
//    }),
//  })
//);

Vue.prototype.$socketio = socketio;
Vue.prototype.$axios = axios;
Vue.config.productionTip = false;

new Vue({
  router,
  store,
  render: (h) => h(App),
}).$mount("#app");
