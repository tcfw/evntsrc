<template>
	<Row type="flex" justify="center" align="middle" :style='{height: "100vh"}'>
		<Col :md="10" :lg="7" :xs="24" :xlg="5">
			<Card :style='{padding: "25px"}' id="login-wrapper">
				<img src="../assets/logo.png" />
				<Form ref="loginForm" :model="loginForm" :rules="loginFormValidationRules">
					<div v-if="hasPriorKnowledge">
						<Row>
							<Col span=5 offset=5>
								<Avatar :src="profileKnowledge.photo" size="large">{{profileKnowledge.name}}</Avatar>
							</Col>
							<Col span=12 offset=1>
								<div id="welcome-back">Welcome Back</div>
								<div id="profile-knowledge">
									<i class="fab google" v-if="profileKnowledge.provider == 'google'" :style='{marginRight: "5px"}'></i>
									<i class="fab facebook" v-if="profileKnowledge.provider == 'facebook'" :style='{marginRight: "5px"}'></i>
									{{profileKnowledge.email}}
								</div>
								<div id="not-you" @click="clearKnowledge">Not you?</div>
							</Col>
						</Row>
					</div>
					<FormItem prop="email" v-if="!hasPriorKnowledge">
						<Input type="text" v-model="loginForm.email" placeholder="Email" />
					</FormItem>
					<FormItem prop="password" v-if="!hasPriorKnowledge || profileKnowledge.provider=='storage'">
						<Input type="password" v-model="loginForm.password" placeholder="Password" />
					</FormItem>
					<FormItem>
						<Button :loading="submitting" ref="loginSubmitBtn" size="large" type="primary" @click="login()" id="login-btn">Log in</Button>
						<router-link to="/forgot" id="forgot-btn" v-if="!hasPriorKnowledge || profileKnowledge.provider=='storage'">Forgot your password?</router-link>
					</FormItem>
				</Form>
				<div v-show="profileKnowledge.provider == 'storage' || !hasPriorKnowledge">
					<div class="login-divider"></div>
					<div id="social-btns">
						<div class="fb-wrapper">
							<div class="fb-login-button" data-scope="public_profile,email" data-width="220px" data-max-rows="1" data-size="medium" data-button-type="continue_with" data-show-faces="false" data-auto-logout-link="false" data-onlogin="app.$route.matched[0].instances.default.fbClickCallback()"></div>
						</div>
						<div class="gapi-wrapper" @click="googleClick">
							<div id="gapi-signin2"></div>
						</div>

						<div id="create-btn"><router-link to="/signup">Create an account...</router-link></div>
					</div>
				</div>
				<div v-if="loading">
					<Spin size="large" fix>
						<Icon type="load-c" size=18 :style='{animation: "ani-demo-spin 1s linear infinite"}'></Icon>
						<div>Loading...</div>
					</Spin>
				</div>
			</Card>
		</Col>
	</Row>
</template>
<script>
import passport from '@/protos/passport_pb.js';
import errorReader from '@/protos/error.js';

export default {
	name: 'Login',
	data () {
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
				tokens: {},
			},
			loginForm: {
				email: '',
				password: ''
			},
			loginFormValidationRules: {
				email: [
					{ required: true, message: 'Please fill in your email', trigger: 'blur' }
				],
				password: [
					{ required: true, message: 'Please fill in your password.', trigger: 'blur' },
					{ type: 'string', min: 6, message: 'The password length cannot be less than 6 characters', trigger: 'blur' }
				]
			}
		}
	},
	mounted() {
		this.$root.$refs.App.bodyClass="pg-login";
		this.$root.$on('gapi.loaded', () => {
			gapi.signin2.render('gapi-signin2', {
				scope: 'profile email',
				width: 220,
				height: 28,
				longtitle: true,
				theme: 'light',
				onSuccess: this.googleLoginCallback,
				onError: this.googleLoginFailed,
			});
		});
		this.$root.$on('fb.loaded', () => {
			FB.XFBML.parse();
			FB.getLoginStatus(r => {
				if (r.status == "connected") {
					this.profileKnowledge.tokens = r.authResponse;
					this.fbLoginCallback();
				}
			});
		});
		this.checkForLocalProfileKnowledge();

		if("gapi" in window) {
			this.$root.$emit("gapi.loaded");
		}
		if("FB" in window) {
			this.$root.$emit("fb.loaded");
		}
	},
	computed: {
		logoColor() {
			if(window.innerWidth < 768) {
				return 'white';	
			}

			return 'black';
		}
	},
	methods: {
		clearKnowledge() {
			this.hasPriorKnowledge = false;

			if (this.profileKnowledge.provider == "google") {
				gapi.auth2.getAuthInstance().disconnect();
			}
			if (this.profileKnowledge.provider == "facebook") {
				FB.logout();
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
				tokens: {},
			};

		},
		checkForLocalProfileKnowledge() {
			if (localStorage) {
				if (localStorage.getItem("prokno") !== null 
					&& localStorage.getItem("prokno-e") !== null
					&& localStorage.getItem("prokno-n") !== null) {
					this.profileKnowledge.provider = "storage";
					this.loginForm.email = this.profileKnowledge.email = localStorage.getItem("prokno-e");
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

			if(this.didClickThrough == true) {
				this.login();
			}
		},
		googleLoginFailed() {
			this.$Message.error("Failed to login using Google");
		},
		fbLoginCallback(callback) {
			FB.getLoginStatus(r => {
				FB.api('/me?fields=name,email,picture', r => {
					this.profileKnowledge.provider = "facebook";
					this.profileKnowledge.name = r.name;
					this.profileKnowledge.email = r.email;
					this.profileKnowledge.photo = r.picture.data.url;
					this.profileKnowledge.tokens = FB.getAuthResponse();

					this.hasPriorKnowledge = true;

					if (callback) {
						callback();
					}
				});
			});
		},
		fbClickCallback() {
			this.didClickThrough = true;
			this.fbLoginCallback(this.login);
		},
		socialLogin() {
			var socialTokens = new passport.Tokens;

			switch(this.profileKnowledge.provider) {
				case "facebook":
					socialTokens.setToken(FB.getAccessToken());
					break;
				case "google":
					socialTokens.setToken(gapi.auth2.getAuthInstance().currentUser.get().getAuthResponse().id_token)
					break;
			}

			var socialRequest = new passport.SocialRequest;
			socialRequest.setProvider(this.profileKnowledge.provider);
			socialRequest.setIdptokens(socialTokens);

			axios.post(this.$root.rootEndpoint+"auth/social", socialRequest.serializeBinary(),{
				headers:{'Content-Type':'application/protobuf'},
				transformResponse: [function (data) {
					return data;
				}],
				responseType: 'arraybuffer'
			}).then(d => {
				this.loading = false;
				if ("data" in d) {
					this.readinAuthResponse(d, false);
				}
			}).catch(e => {
				this.loading = false;
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
			var authResponse = passport.AuthResponse.deserializeBinary(response.data);

			if(!authResponse.getSuccess()) {
				this.$Message.error({
					content: "Unable to log you in. Please try again.", 
					duration: 10
				});
				return
			}

			var expires = new Date(0)
			expires.setUTCSeconds(authResponse.getTokens().getTokenexpire().getSeconds())

			this.$cookie.set('session', authResponse.getTokens().getToken(), {expires: expires})
			this.$root.applySession()

			if (setProfileKnowledge) {
				this.applyProfileKnowledgeFromStandard();
			}

			this.$root.fetchMe();
			this.$router.push('/');
		},
		login() {
			if (this.hasPriorKnowledge && this.profileKnowledge.provider != "storage") {
				this.loading = true;
				this.socialLogin();
			} else {
				this.$refs['loginForm'].validate(valid => {
					if (!valid) {
						this.$Message.error('Failed to login');
					} else {
						this.submitting = true;

						var userCreds = new passport.UserCreds;
						userCreds.setPassword(this.loginForm.password);
						userCreds.setUsername(this.loginForm.email);

						var authRequest = new passport.AuthRequest;
						authRequest.setUsercreds(userCreds);

						axios.post(this.$root.rootEndpoint+"auth/login", authRequest.serializeBinary(),{
							headers:{'Content-Type':'application/protobuf'},
							transformResponse: [function (data) {
								return data;
							}],
							responseType: 'arraybuffer'
						}).then(d => {
							this.submitting = false;
							this.readinAuthResponse(d, true)
						}).catch(de => {
							this.submitting = false;
							var msg = "Incorrect login details";
							if ("response" in de) { 
								var error = JSON.parse(new TextDecoder("utf-8").decode(de.response.data))
								if(de.response.status == 429) {
									app.$Message.destroy()
									this.$Message.error({
										content: error.error,
										duration: 30
									});
								} else {
									if (error.error.includes("no reachable servers")) {
										msg = "Unable to log you in. Please try again.";
									}
									app.$Message.destroy()
									this.$Message.error({
										content: msg, 
										duration: 10
									});
									throw error.message;
								}
							}
						})
					}
				})
			}
		}
	}
}
</script>
<style lang="scss">
.pg-login {
	&:after { //bg
		content: "";
		display: block;
		z-index: -1;
		position: fixed;
		top: -10px;
		left: -10px;
		height: calc(100% + 20px);
		width: calc(100% + 20px);
		background:
			linear-gradient(
				rgba(255, 255,255,0.4), 
				rgba(255, 255,255,0.4)
			),
			url("https://images.unsplash.com/photo-1451187580459-43490279c0fa?ixlib=rb-0.3.5&ixid=eyJhcHBfaWQiOjEyMDd9&s=3f6e5055d9ad1d603fd364c11823d026&auto=format&fit=crop&w=1952&q=80") no-repeat center;
		background-size: cover;
		filter: blur(5px);
	}

	#login-btn {
		float: right;

		@media(max-width: 768px) {
			width: 100%;
		}
	}

	#forgot-btn {
		float: left;

		@media(max-width: 768px) {
			padding-top: 75px;
			display: block;
			text-align: center;
			float: inherit;
			color: white;
		}
	}

	#create-btn {
		margin-top: 20px;

		@media(max-width: 768px) {
			margin-top: 35px;

			> a {
				color: white;
			}
		}
	}

	#welcome-back {
		font-size: 18px;
		color: #515151;
		font-weight: 500;
	}

	#profile-knowledge {
		margin-top: -5px;
	}

	#profile-knowledge, #not-you {
		font-size: 12px;
		color: #515151;
		font-weight: 100;
	}

	#not-you {
		margin-top: 5px;
		margin-bottom: 25px;
		cursor: pointer;
	}

	#social-btns {
		text-align: center;

		@media(max-width: 768px) {
			margin-top: 35px;
		}

		.fb-wrapper {
			display: inline-block;
			background: #4267b2;
			border-radius: 5px;
			color: white;
			height: 28px;
			overflow: hidden;
			text-align: center;
			width: 220px;
			margin: 0px;
		}

		.gapi-wrapper {
			display: block;

			#gapi-signin2 {
				display: inline-block;
				margin-left: -5px;
				margin-top: 15px;
				height: 28px;
			}
		}
	}

	.login-divider {
		background-image: linear-gradient(to right, transparent 50%, #E5E5E5 50%);
		background-size: 9px 100%;
		height: 1px;
		position: relative;
		margin-bottom: 30px;

		@media(max-width: 768px) {
			display: none;
		}

		&:after {
			content: "or";
			display: block;
			position: absolute;
			top: -9px;
			left: 50%;
			background: white;
			padding: 0px 14px;
			color: #BBBBBB;
			font-size: 12px;
			transform: translateX(-50%);

		}
	}

	.ivu-avatar-large {
		width: 60px;
		height: 60px;
		border-radius: 60px;
		box-shadow: 0px 2px 2px rgba(0,0,0,0.2)
	}

	@media(max-width: 768px) {
		.ivu-card {
			background: transparent;
			border: none;
			box-shadow: none;
		}

		.ivu-input {
			box-shadow: none;
		}
	}
}
</style>