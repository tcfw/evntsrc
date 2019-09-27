<template>
	<div class="h-screen bg-white shadow-lg p-6 pb-2 xl:ml-64 xl:w-5/6 xl:p-20 xl:pb-6 md:ml-12 md:w-2/3 md:p-16 lg:w-3/6 lg:ml-12 xxl:w-1/3 xxl:ml-64 xxl:px-32 xxl:pt-32 overflow-auto" style="max-width: 600px;">
		<div class="logo text-center">
			<img src="@/assets/logo.png" class="h-20 -ml-2 sm:ml-0" />
		</div>
		<div v-if="!verified && !error" class="mt-32 text-center">
			<div class="text-2xl">Verifying your account...</div>
			<div class="mt-6">
				<i class="fas fa-sync fa-2x fa-spin"></i>
			</div>
		</div>
		<div v-if="error" class="mt-32 text-center">
			Something went wrong while trying to verify your account.<br/><br/>
			Please try again by refreshing the page or<br/>
			Contact us via <a href="mailto:hello@evntsrc.io" class="text-ev-100 font-bold">hello@evntsrc.io</a>

			<div class="mt-12"><a href="javascript:location.reload()" class="text-ev-100 font-bold">Refresh</a></div>
		</div>
		<div class="mt-16 text-center" v-if="verified">
			<img src="@/assets/ready.png" class="inline" alt="Set up complete"/>
			<div class="mt-12">Your account is set to go! Login to set up your first stream...</div>
			<div class="mt-12">
				<router-link to="/login" class="input-button-huge">Login</router-link>
			</div>
		</div>
	</div>
</template>
<script>
export default {
	name: "verify",
	data() {
		return {
			verified: false,
			email: "",
			error: null,
			token: null,
		}
	},
	mounted() {
		//This is usually done by the router, except for landing pages
		this.$root.$refs.App.appClass = "pg-verify";
		this.email = this.$route.query.e
		this.token = this.$route.params.token
		if (this.email == "" || this.email == null) {
			this.error = "No email found in request"
		}
		this.email = atob(this.email)
		this.$http.post(
			this.$config.API + "/auth/validate_account",
			{
				email: this.email,
				token: this.token
			}
		).then(d => {
			this.verified = true
		})
		.catch(d => {
			this.error = "Something went wrong"
		})
	}
}
</script>