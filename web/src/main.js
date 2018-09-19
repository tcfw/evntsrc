import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import "./plugins/element.js";
import "./plugins/axios.js";
import "./plugins/config.js";

Vue.config.productionTip = false;

new Vue({
  router,
  render: h => h(App),
  data() {
    return {
      apiEndpoint: "localhost:34234"
    }
  }
}).$mount("#app");
