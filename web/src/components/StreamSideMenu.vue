<template>
	<div class="stream-side-menu">
		<div class="search-container">
			<input name="search" type="text" placeholder="Search..." v-model="searchInput" @keyup.enter="load" @keyup.esc="clearSearch"/>
			<i class="fas fa-search" @click="load"></i>
			<i class="fas fa-times clear-search" v-if="isSearchResults" @click="clearSearch"></i>
		</div>
		<div class="stream-list" v-if="!loading">

		</div>
		<div class="searching-loader" v-if="loading">
			<i class="fas fa-sync fa-spin"></i>
		</div>
		<div class="loading-error" v-if="error">
			<i class="fas fa-exclamation-triangle" :style="{color: '#f56c6c'}"></i>
			<p>{{error}}</p>
			<el-button type="danger" size="mini" @click="retry">Retry</el-button>
		</div>
	</div>
</template>
<script>
export default {
	name: "stream-side-menu",
	data () {
		return {
			searchInput: "",
			loading: false,
			error: null,
			isSearchResults: false
		}
	},
	methods: {
		clearSearch() {
			this.searchInput = "";
			this.load();
		},
		load() {
			this.loading = true;
			this.error = null;
			this.isSearchResults = false;
			let url = this.$config.API+"/streams";
			if (this.searchInput != "") {
				url += "/search?query="+this.searchInput
			}
			axios.get(url).then(d => {
				this.buildList(d.data)
			}).catch(e => {
				this.error = "Failed to "+((this.searchInput)?"search":"load")+" streams...";
			}).finally(() => {
				if (this.searchInput != "") {
					this.isSearchResults = true;
				}
				this.loading = false;
			})
		},
		buildList(streams) {

		},
		retry() {
			if(this.searchInput != "") {
				this.search()
			} else {
				this.load()
			}
		}
	},
	mounted() {
		this.load();
	},
}
</script>

<style lang="scss" scoped>
.stream-side-menu {
	background: #50566F;
	width: 100%;
	height: 100%;
	box-shadow: 2px 0 4px 0 rgba(0,0,0,0.10);
	position: relative;
}

.search-container {
	margin: 15px;
	width: calc(100% - 30px);;
	height: 36px;
	border-radius: 3px;
	background: #787E99;
	position: absolute;	
	box-shadow: inset 0 1px 3px 0 rgba(0,0,0,0.1);

	input[type="text"] {
		background: none;
		border: none;
		color: white;
		height: calc(100% - 2px);
		margin-left: 10px;
		outline: none;
		font-weight: 400;
		font-size: 12px;

		&::placeholder {
			font-weight: 100;
			color: #BBB;
		}

	}
	.fas {
		position: absolute;
		right: 15px;
		top: 12px;
		font-size: 10px;
		color: white;
		cursor: pointer;
	}
	.clear-search {
		right: 35px;
	}
}
.searching-loader {
	position: absolute;
	left: 50%;
	top: 80px;
	transform: translateX(-10px);
	color: white;
	font-size: 14px;
}
.loading-error {
	color: white;
	font-size: 12px;
	position: absolute;
	top: 80px;
	text-align: center;
	width: 100%;
}
</style>
