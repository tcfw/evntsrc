<template>
	<div class="h-screen bg-white shadow-lg p-6 pb-2 xl:ml-64 xl:w-5/6 xl:p-20 xl:pb-6 md:ml-12 md:w-2/3 md:p-16 lg:w-3/6 lg:ml-12 xxl:w-1/3 xxl:ml-64 xxl:px-32 xxl:pt-32 overflow-auto" style="max-width: 600px;">
		<div class="logo text-center">
			<img src="@/assets/logo.png" class="h-20 -ml-2 sm:ml-0" />
		</div>
		<router-link to="/login" class="float-right mt-10 text-sm text-text-800">Login</router-link>
		<h2 class="font-bold text-xl text-text-100 pt-8 pb-6">Verify your email</h2>
		<div class="mt-6">
			<p>We've sent you an email to validate your account. If you don't receive the email within a few minutes, double check your Junk/Spam folder.</p>
		</div>
		<div class="text-center w-full mt-16">
			<img class="inline max-w-sm" src="@/assets/email_sent.png"/>
		</div>
		<!-- TODO(tcfw) <a href="#" class="mt-16 text-xs text-ev-100 block font-bold">I still didn't receive my email</a> -->
	</div>
</template>
<script>
export default {
	name: 'signup-thanks',
	beforeRouteEnter (to, from, next) {
		if (from.path != "/signup") {
			next('/signup')
		} else {
			next()
		}
	}
}
</script>