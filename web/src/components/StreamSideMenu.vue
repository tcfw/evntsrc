<template>
	<div class="stream-side-menu">
		<div class="search-container">
			<input name="search" type="text" placeholder="Search..."/>
		</div>
	</div>
</template>
<style lang="scss" scoped>
.stream-side-menu {
	background: #50566F;
	width: 100%;
	height: 100%;
	box-shadow: 2px 0 4px 0 rgba(0,0,0,0.10);
}

.search-container {
	margin: 15px;
	width: calc(100% - 30px);;
	height: 36px;
	border-radius: 3px;
	background: #787E99;
	position: absolute;	
	box-shadow: inset 0 1px 3px 0 rgba(0,0,0,0.20);

	input[type="text"] {
		background: none;
		border: none;
		color: white;
		height: calc(100% - 2px);
		margin-left: 10px;
		outline: none;

		&::placeholder {
			color: #C7C7C7;
		}
	}
}
</style>
