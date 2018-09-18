<template>
	<div class="side-menu">
		<router-link to="/"><img src="@/assets/logo.png"/></router-link>
		<div class="menu">
			<router-link to="/dashboard"><i class="fas fa-tachometer-alt"></i> Dashboard</router-link>
			<router-link to="/streams"><i class="fas fa-stream"></i> Streams</router-link>
		</div>
	</div>
</template>
<script>
export default {
	
}
</script>
<style lang="scss" scoped>
.side-menu {
	background: #1A1E30;
	height: 100%;
	box-shadow: 2px 2px 4px 0 rgba(0,0,0,0.10);

	img {
		height: 18px;
		padding: 25px 15px 20px;
	}

	.menu {
		a {
			color: white;
			padding: 15px;
			font-weight: 100;
			font-size: 12px;
			text-decoration: none;
			display: block;

			.fas {
				padding: 0px 10px;
				display: inline-block;
			}

			&.router-link-active {
				background: rgba(255,255,255,0.12);
				border-left: 4px solid #636D97;
				padding-left: 11px;
			}
		}
	}
}
</style>
