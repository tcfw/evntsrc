<template>
	<div v-loading.lock="loading" :style="{height: '100%', width: '100%'}">
		<el-tabs class="stream-tabs">
			<el-tab-pane label="STATS" name="stats">Stats</el-tab-pane>
			<el-tab-pane label="HISTORY" name="history">History</el-tab-pane>
			<el-tab-pane label="AUTH" name="auth">Auth</el-tab-pane>
			<el-tab-pane label="INGRESS" name="ingress">Ingress</el-tab-pane>
			<el-tab-pane label="WEBHOOKS" name="webhooks">Webhooks</el-tab-pane>
			<el-tab-pane label="SETTINGS" name="settings">Settings</el-tab-pane>
		</el-tabs>
	</div>
</template>
<script>
export default {
	name: "stream",
	data() {
		return {
			loading: true,
			error: "",
			stream: null,
		}
	},
	methods: {
		load() {
			// this.loading = false;
		}
	},
	mounted() {
		this.load();
	}
}
</script>
<style lang="scss">
.stream-tabs {
	.el-tabs__nav-wrap {
		padding: 6px 0px 0px 15px;
		background: white;
		box-shadow: 0 2px 3px 0 rgba(0,0,0,0.10);
		
		.el-tabs__item {
			font-size: 13px;
			font-weight: 300;
			color: #50566F;
			padding-bottom: 45px;

			&.is-active {
				font-weight: 400;
			}
		}

		&::after {
			display: none;
		}
	}
}
</style>
