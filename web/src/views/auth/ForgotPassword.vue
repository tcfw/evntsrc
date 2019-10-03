<template>
	<div class="h-screen bg-white shadow-lg p-6 pb-2 xl:ml-64 xl:w-5/6 xl:p-20 xl:pb-6 md:ml-12 md:w-2/3 md:p-16 lg:w-3/6 lg:ml-12 xxl:w-1/3 xxl:ml-64 xxl:px-32 xxl:pt-32 overflow-auto" style="max-width: 600px;">
		<div class="logo text-center">
			<img src="@/assets/logo.png" class="h-20 -ml-2 sm:ml-0" />
		</div>
		<router-link to="/login" class="float-right mt-10 text-sm text-text-800">Login</router-link>
		<h2 class="font-bold text-xl text-text-100 pt-8 pb-6">Account Recovery</h2>
		<!-- TODO(tcfw) -->
	</div>
</template>
<script>
export default {
	name: 'forgot',
	mounted() {
		 //This is usually done by the router, except for landing pages
		this.$root.$refs.App.appClass = "pg-forgot";
	}
}
</script>