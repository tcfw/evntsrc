import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import "./plugins/element.js";
import "./plugins/lodash.js";
import "./plugins/axios.js";
import "./plugins/config.js";
import "./plugins/google.js";
import "./plugins/cookie.js";
import "./plugins/moment.js";
import "@/styles/app.scss";

router.afterEach(route => {
  if ("App" in router.app.$refs) {
    router.app.$refs.App.appClass = "pg-" + route.name;
  }
});

Vue.config.productionTip = false;

const publicRoutes = [
  'login',
  'signup',
  'signup-thanks',
  'terms',
  'privacy'
];

let app = new Vue({
  router,
  render: h => h(App),
  data() {
    return {
      me: {},
      stripe: {}
    };
  },
  mounted() {
    if (!this.$root.loggedIn() && publicRoutes.indexOf(this.$route.name) < 0) {
      this.$router.push("/login");
    } else {
      this.fetchMe();
    }
    this.loadStripe();
  },
  methods: {
    fetchMe() {
      if (this.loggedIn()) {
        this.$http.get(this.$config.API + "/me").then(d => {
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
      this.$http.post(this.$config.API + "/auth/revoke");
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
      this.$http.defaults.headers.common["Authorization"] = session;
    },
    gapiCallback() {
      this.$emit("gapi.loaded");
    },
    loadStripe() {
      this.stripe = Stripe(process.env.VUE_APP_STRIPE_TOKEN);
    }
  }
});
app.$mount("#app");
