<template>
  <div class="h-screen bg-white shadow-lg p-6 xl:ml-64 xl:w-5/6 xl:p-20 md:ml-12 md:w-2/3 md:p-16 lg:w-3/6 lg:ml-12 xxl:w-1/3 xxl:ml-64 xxl:px-32 xxl:pt-32" style="max-width: 700px;">
    <div class="logo text-center">
        <img src="../assets/logo.png" class="h-20 -ml-2 sm:ml-0" />
      </div>
      <form action="#" @submit.prevent.stop="login">
        <h2 class="font-bold text-xl text-text-100 pt-8 pb-6">Login to your account</h2>
        <div v-if="hasPriorKnowledge">
            <div class="flex mt-1 -mb-2">
              <div class="w-20">
                <Avatar :src="profileKnowledge.photo" size="large">
                  {{profileKnowledge.name}}
                </Avatar>
              </div>
              <div>
                <div class="font-bold">Welcome Back</div>
                <div id="profile-knowledge">
                  <i class="fab fa-google" v-if="profileKnowledge.provider == 'google'" :style="{ marginRight: '5px' }"></i>
                  {{ profileKnowledge.email }}
                </div>
                <div id="not-you" @click="clearKnowledge">Not you?</div>
            </div>
            </div>
        </div>
        <div v-if="!hasPriorKnowledge">
          <label class="block">Email</label>
          <input type="text" tabindex="1" v-model="loginForm.email" @keydown.enter.stop.prevent.native="login" class="input-text" />
        </div>
        <div v-if="!hasPriorKnowledge || profileKnowledge.provider == 'storage'" class="mt-3">
          <router-link tabindex="4" class="float-right text-sm pt-1 text-text-800" to="/forgot" v-if="!hasPriorKnowledge || profileKnowledge.provider == 'storage'">
            Forgot your password?
          </router-link>
          <label class="block">Password</label>
          <input type="password" tabindex="2" v-model="loginForm.password" @keydown.enter.stop.prevent.native="login" class="input-text" />
        </div>
        <div>
          <button tabindex="3" ref="loginSubmitBtn" :disabled="submitting" @click="login()" id="login-btn" class="input-button-huge mt-10 mb-6">
            {{loginBtn}}
          </button>
        </div>
      </form>
      <div v-show="profileKnowledge.provider == 'storage' || !hasPriorKnowledge">
        <div id="social-btns">
          Or login using
          <span class="social-acc" @click="googleClick">Google</span>
          <span class="social-acc">Github</span>
        </div>
      </div>
      <div class="mt-16 lg:mt-32 xxl:mt-48">
        <div class="mb-3 text-text-300 text-sm">Don't have an account yet?</div>
        <router-link to="/signup" class="input-button-huge clear ">Create an account</router-link>
      </div>
      <div v-if="loading">
        <Spin size="large" fix>
          <Icon
            type="load-c"
            size="18"
            :style="{ animation: 'ani-demo-spin 1s linear infinite' }"
          ></Icon>
          <div>Loading...</div>
        </Spin>
      </div>
  </div>
</template>
<script>
import errorReader from "@/protos/error.js";
import Avatar from "@/components/Avatar";

const proto = require('@/protos/passport_pb.js')

export default {
  name: "Login",
  components: {
    Avatar
  },
  data() {
    return {
      submitting: false,
      loading: false,
      hasPriorKnowledge: false,
      didClickThrough: false,
      profileKnowledge: {
        photo: null,
        name: null,
        email: null,
        provider: "storage",
        tokens: {}
      },
      loginForm: {
        email: "",
        password: ""
      },
      loginBtn: "Login",
      loginFormValidationRules: {
        email: [
          {
            required: true,
            message: "Please fill in your email",
            trigger: "blur"
          }
        ],
        password: [
          {
            required: true,
            message: "Please fill in your password.",
            trigger: "blur"
          },
          {
            type: "string",
            min: 6,
            message: "The password length cannot be less than 6 characters",
            trigger: "blur"
          }
        ]
      }
    };
  },
  mounted() {
    this.$root.$refs.App.appClass = "pg-login";
    this.$root.$on("gapi.loaded", () => {
      gapi.signin2.render("gapi-signin2", {
        scope: "profile email",
        width: 220,
        height: 28,
        longtitle: true,
        theme: "light",
        onSuccess: this.googleLoginCallback,
        onError: this.googleLoginFailed
      });
    });
    this.checkForLocalProfileKnowledge();

    if ("gapi" in window) {
      this.$root.$emit("gapi.loaded");
    }
    if ("FB" in window) {
      this.$root.$emit("fb.loaded");
    }

    if (this.$root.loggedIn()) {
      this.$router.push("/");
    }
  },
  computed: {
    logoColor() {
      if (window.innerWidth < 768) {
        return "white";
      }

      return "black";
    }
  },
  methods: {
    clearKnowledge() {
      this.hasPriorKnowledge = false;

      if (this.profileKnowledge.provider == "google") {
        gapi.auth2.getAuthInstance().disconnect();
      }

      if (this.profileKnowledge.provider == "storage") {
        localStorage.removeItem("prokno");
        localStorage.removeItem("prokno-e");
        localStorage.removeItem("prokno-n");
        this.loginForm.email = "";
      }

      this.profileKnowledge = {
        photo: null,
        name: null,
        email: null,
        provider: "storage",
        tokens: {}
      };
    },
    checkForLocalProfileKnowledge() {
      if (localStorage) {
        if (
          localStorage.getItem("prokno") !== null &&
          localStorage.getItem("prokno-e") !== null &&
          localStorage.getItem("prokno-n") !== null
        ) {
          this.profileKnowledge.provider = "storage";
          this.profileKnowledge.photo = this.loginForm.email = this.profileKnowledge.email = localStorage.getItem(
            "prokno-e"
          );
          this.profileKnowledge.name = localStorage.getItem("prokno-n");
          this.hasPriorKnowledge = true;
        } else {
        }
      }
    },
    googleClick() {
      this.didClickThrough = true;
    },
    googleLoginCallback(r) {
      if (!r) return;
      let basicProfile = r.getBasicProfile();
      this.profileKnowledge.provider = "google";
      this.profileKnowledge.name = basicProfile.getName();
      this.profileKnowledge.email = basicProfile.getEmail();
      this.profileKnowledge.photo = basicProfile.getImageUrl();
      this.profileKnowledge.tokens = r.getAuthResponse();

      this.hasPriorKnowledge = true;

      if (this.didClickThrough == true) {
        this.login();
      }
    },
    googleLoginFailed() {
      this.$Message.error("Failed to login using Google");
    },
    socialLogin() {
      var socialTokens = new proto.Tokens();

      switch (this.profileKnowledge.provider) {
        case "google":
          socialTokens.setToken(
            gapi.auth2
              .getAuthInstance()
              .currentUser.get()
              .getAuthResponse().id_token
          );
          break;
      }

      var socialRequest = new proto.SocialRequest();
      socialRequest.setProvider(this.profileKnowledge.provider);
      socialRequest.setIdptokens(socialTokens);

      this.$http.post(
          this.$config.API + "/auth/social",
          socialRequest.serializeBinary(),
          {
            headers: { "Content-Type": "application/protobuf" },
            transformResponse: [
              function(data) {
                return data;
              }
            ],
            responseType: "arraybuffer"
          }
        )
        .then(d => {
          this.loading = false;
          if ("data" in d) {
            this.readinAuthResponse(d, false);
          }
        })
        .catch(e => {
          this.loading = false;
          this.loginBtn = "Login";
          this.$Message.error({
            content: "Unable to log you in. Please try again.",
            duration: 10
          });
        });
    },
    applyProfileKnowledgeFromStandard() {
      localStorage.setItem("prokno", true);
      localStorage.setItem("prokno-e", this.loginForm.email);
      localStorage.setItem("prokno-n", this.loginForm.email);
    },
    readinAuthResponse(response, setProfileKnowledge) {
      var authResponse = proto.AuthResponse.deserializeBinary(response.data);

    if (!authResponse.getSuccess()) {
        this.$message.error({
          content: "Unable to log you in. Please try again.",
          duration: 10
        });
        return;
      }

      this.$message.closeAll();

      var expires = new Date(0);
      expires.setUTCSeconds(
        authResponse
          .getTokens()
          .getTokenexpire()
          .getSeconds()
      );

      var token = authResponse.getTokens().getToken()

      var cookieSettings = {
        expires: expires,
        secure: process.env.VUE_APP_SEC_COOKIE == "true",
      }

      if (process.env.VUE_APP_SEC_DOMAIN != "") {
        cookieSettings = process.env.VUE_APP_SEC_DOMAIN
      }

      this.$cookie.set("session", token, cookieSettings);
      this.$root.applySession();

      if (setProfileKnowledge) {
        this.applyProfileKnowledgeFromStandard();
      }

      this.$root.fetchMe();
      this.$router.push("/");
    },
    login() {
      if (
        this.hasPriorKnowledge &&
        this.profileKnowledge.provider != "storage"
      ) {
        this.loading = true;
        this.loginBtn = "...";
        this.socialLogin();
      } else {
        // this.$refs["loginForm"].validate(valid => {
        //   if (!valid) {
        //     this.$Message.error("Failed to login");
        //   } else {
            this.submitting = true;
            this.loginBtn = "...";
            var userCreds = new proto.UserCreds();
            userCreds.setPassword(this.loginForm.password);
            userCreds.setUsername(this.loginForm.email);

            var authRequest = new proto.AuthRequest();
            authRequest.setUsercreds(userCreds);

            this.$http.post(
                this.$config.API + "/auth/login",
                authRequest.serializeBinary(),
                {
                  headers: { "Content-Type": "application/protobuf" },
                  transformResponse: [
                    function(data) {
                      return data;
                    }
                  ],
                  responseType: "arraybuffer"
                }
              )
              .then(d => {
                this.submitting = false;
                this.loginBtn = "Login";
                this.readinAuthResponse(d, true);
              })
              .catch(de => {
                this.submitting = false;
                this.loginBtn = "Login";
                var msg = "Incorrect login details";
                if ("response" in de) {
                  var error = JSON.parse(
                    new TextDecoder("utf-8").decode(de.response.data)
                  );
                  if (de.response.status == 429) {
                    this.$message({
                      message: error.error,
                      duration: 30000,
                      type: "error"
                    });
                  } else {
                    if (error.error.includes("no reachable servers")) {
                      msg = "Unable to log you in. Please try again.";
                    }
                    this.$message({
                      message: msg,
                      duration: 30000,
                      type: "error"
                    });
                    throw error.message;
                  }
                }
              });
          // }
        // });
      }
    }
  }
};
</script>