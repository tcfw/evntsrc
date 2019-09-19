<template>
	<div class="absolute text-center" style="top: 50%; left: 50%; transform: translate(-50%, -50%);">
		<img src="@/assets/404.png" class="w-1/2 inline" /><br/><br/>
		Page not found<br/><br/>
		You look lost :(<br/> Time to go <router-link to="/" class="text-ev-100">home</router-link> or go <a @click="$router.go(-1)" class="text-ev-100 cursor-pointer">back</a> the way you came...
	</div>
</template>