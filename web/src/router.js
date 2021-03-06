import Vue from "vue";
import Router from "vue-router";

import Login from "./views/auth/Login.vue";
import Signup from "./views/auth/Signup.vue";
import SignupComplete from "./views/auth/SignupComplete.vue";
import ForgotPassword from "./views/auth/ForgotPassword.vue";
import VerifyEmail from "./views/auth/VerifyEmail.vue";

import Dashboard from "./views/Dashboard.vue";

import Streams from "./views/Streams.vue";
import Stream from "./views/Stream.vue";
import StreamHistory from "./views/Stream/History.vue";
import StreamKeys from "./views/Stream/Keys.vue";
import StreamSettings from "./views/Stream/Settings.vue";

import Settings from "./views/Settings.vue";
import SettingsAccount from "./views/Settings/Account.vue";
import SettingsSecurity from "./views/Settings/Security.vue";
import SettingsSessions from "./views/Settings/Sessions.vue";

import e404 from "./views/errors/404.vue";

Vue.use(Router);

export default new Router({
  mode: "history",
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
      path: "/signup",
      name: "signup",
      component: Signup
    },
    {
      path: "/signup/thanks",
      name: "signup-thanks",
      component: SignupComplete
    },
    {
      path: "/verify/:token",
      name: "verify",
      component: VerifyEmail,
    },
    {
      path: "/forgot",
      name: "forgot",
      component: ForgotPassword,
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
              name: "stream-settings",
              component: StreamSettings
            },
            {
              path: "history",
              name: "stream-history",
              component: StreamHistory
            },
            {
              path: "keys",
              name: "stream-auth",
              component: StreamKeys
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
          component: () =>
            import(/* webpackChunkName: "billing" */ "./views/Settings/Billing.vue")
        },
        {
          path: "sessions",
          name: "sessions",
          component: SettingsSessions
        }
      ]
    },
    {
      path: '/404',
      alias: '*',
      name: '404',
      component: e404
    }
  ]
});
