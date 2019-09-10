import Vue from "vue";
import axiosRetry from "axios-retry";

const axios = require("axios");

Vue.prototype.$http = axios;

axiosRetry(axios, { retries: 3 });
