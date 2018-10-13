import Vue from "vue";
import Router from "vue-router";
import Login from "./views/Login.vue";
import Dashboard from "./views/Dashboard.vue";
import Streams from "./views/Streams.vue";
import Stream from "./views/Stream.vue";
import StreamHistory from "./views/Stream/History.vue";
import Settings from "./views/Settings.vue";
import SettingsAccount from "./views/Settings/Account.vue";
import SettingsSecurity from "./views/Settings/Security.vue";
import SettingsBilling from "./views/Settings/Billing.vue";

Vue.use(Router);

export default new Router({
  routes: [
    {
      path: "/",
      name: "home",
      redirect: "/dashboard"
    },
    {
      path: "/login",
      name: "login",
      component: Login
    },
    {
      path: "/dashboard",
      name: "dashboard",
      component: Dashboard
    },
    {
      path: "/streams",
      name: "streams",
      component: Streams,
      children: [
        {
          path: ":id",
          name: "stream",
          component: Stream,
          props: true,
          children: [
            {
              path: "settings",
              name: "stream-settings"
            },
            {
              path: "history",
              name: "stream-history",
              component: StreamHistory
            }
          ]
        }
      ]
    },
    {
      path: "/about",
      name: "about",
      // route level code-splitting
      // this generates a separate chunk (about.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () =>
        import(/* webpackChunkName: "about" */ "./views/About.vue")
    },
    {
      path: "/settings",
      name: "settings",
      redirect: "/settings/account",
      component: Settings,
      children: [
        {
          path: "account",
          name: "account",
          component: SettingsAccount
        },
        {
          path: "security",
          name: "security",
          component: SettingsSecurity
        },
        {
          path: "billing",
          name: "billing",
          component: SettingsBilling
        }
      ]
    }
  ]
});
