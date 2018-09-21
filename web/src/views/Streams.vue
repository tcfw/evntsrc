<template>
	<el-container :style="{height: 'calc(100% - 60px)', position: 'absolute', width: 'calc(100% - 220px)'}">
		<el-aside width="220px">
			<stream-side-menu></stream-side-menu>
		</el-aside>
		<el-container>
			<el-main>
				<router-view></router-view>
			</el-main>
		</el-container>
	</el-container>
</template>
<script>
import streamSideMenu from "@/components/StreamSideMenu.vue";

export default {
	name: "streams",
	components: {
		streamSideMenu
	}
}
</script>
