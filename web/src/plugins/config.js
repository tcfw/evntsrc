import Vue from "vue";
const vueConfig = require("vue-config");
const configs = {
  API: process.env.VUE_APP_API_ENDPOINT
};

Vue.use(vueConfig, configs);
