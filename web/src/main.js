import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import "./plugins/element.js";
import "./plugins/lodash.js";
import "./plugins/axios.js";
import "./plugins/config.js";
import "./plugins/google.js";
import "./plugins/facebook.js";
import "./plugins/cookie.js";
import "./plugins/moment.js";
import "@/styles/app.scss";

router.afterEach(route => {
  if ("App" in router.app.$refs) {
    router.app.$refs.App.appClass = "pg-" + route.name;
  }
});

Vue.config.productionTip = false;
window.app = new Vue({
  router,
  render: h => h(App),
  data() {
    return {
      me: {},
      stripe: {}
    };
  },
  mounted() {
    if (!this.$root.loggedIn() && this.$route.name != "login") {
      this.$router.push("login");
    } else {
      this.fetchMe();
    }
    this.loadStripe();
  },
  methods: {
    fetchMe() {
      if (this.loggedIn()) {
        axios.get(this.$config.API + "/me").then(d => {
          this.me = d.data;
          this.$emit("me.ready", this.me);
        });
      }
    },
    loggedIn() {
      this.applySession();
      return this.$cookie.get("session") != null;
    },
    logout() {
      this.$cookie.delete("session");
      this.applySession();
      this.$message({
        message: "You have been logged out successfully",
        type: "success"
      });
      this.$router.push("/login");
    },
    applySession() {
      let session = this.$cookie.get("session");
      axios.defaults.headers.common["Authorization"] = session;
    },
    gapiCallback() {
      this.$emit("gapi.loaded");
    },
    fbCallback() {
      this.$emit("fb.loaded");
    },
    loadStripe() {
      this.stripe = Stripe(process.env.VUE_APP_STRIPE_TOKEN);
    }
  }
});
app.$mount("#app");
