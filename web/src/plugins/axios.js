import axiosRetry from "axios-retry";

const axios = require("axios");
window.axios = axios;

axiosRetry(axios, { retries: 3 });
