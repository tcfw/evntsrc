<template>
	<div class="h-screen bg-white shadow-lg p-6 pb-2 xl:ml-64 xl:w-5/6 xl:p-20 xl:pb-6 md:ml-12 md:w-2/3 md:p-16 lg:w-3/6 lg:ml-12 xxl:w-1/3 xxl:ml-64 xxl:px-32 xxl:pt-32 overflow-auto" style="max-width: 600px;">
		<div class="logo text-center">
			<img src="@/assets/logo.png" class="h-20 -ml-2 sm:ml-0" />
		</div>
		<router-link to="/login" class="float-right mt-10 text-sm text-text-800">Already have an account?</router-link>
		<h2 class="font-bold text-xl text-text-100 pt-8 pb-6">Create an account</h2>
		<form action="#" @submit.prevent.stop="signup">
			<div>
				<label class="block">Email *</label>
				<input type="email" tabindex="1" name="email" v-model="signUpForm.email" class="input-text" />
			</div>
			<div class="mt-3">
				<label class="block">Password *</label>
				<input type="password" tabindex="2" name="password" v-model="signUpForm.password" autocomplete="new-password" class="input-text" />
			</div>
			<div class="mt-3">
				<label class="block">Name *</label>
				<input type="text" tabindex="3" name="name" v-model="signUpForm.name" class="input-text" />
			</div>
			<!-- <div class="mt-3">
				<label class="block">Company</label>
				<input type="text" tabindex="4" name="company" v-model="signUpForm.company" class="input-text" />
			</div>
			<div class="mt-3">
				<label class="block">Country *</label>
				<input type="text" tabindex="5" name="country" v-model="signUpForm.country" class="input-text" placeholder="Select a country..." />
			</div> -->
			<div class="mt-6">By clicking "Create account" below, you agree to our <a href="" class="font-bold text-ev-100">Terms of Service</a> and <a href="" class="font-bold text-ev-100">Privacy Statement</a></div>
			<div>
				<button tabindex="6" :disabled="submitting" @click.stop.prevent="signup()" class="input-button-huge mt-6">
					{{btn}}
				</button>
			</div>
		</form>
	</div>
</template>
<script>
export default {
	name: 'signup',
	data() {
		return {
			signUpForm: {
				email: null,
				password: null,
				name: null,
				company: null,
				country: null,
				time: null,
			},
			btn: "Create account",
			submitting: false,
		}
	},
	methods: {
		signup() {
			this.btn = "..."
			this.submitting = true;
			this.$http.post(
				this.$config.API + "/auth/register",
				this.signUpForm
			).then(d => {
				this.submitting = false;
				this.btn = "Create account"
				if (d.status == 200) {
					this.$router.push('/signup/thanks')
				}
			}).catch(d => {
				this.submitting = false;
				this.btn = "Create account"
				this.$message.error("Something went wrong whilst trying to register. Please try again")
			})
		}
	},
	mounted() {
		 //This is usually done by the router, except for landing pages
		this.$root.$refs.App.appClass = "pg-signup";
	}
}
</script>