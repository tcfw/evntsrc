<template>
	<div class="stream-side-menu">
		<div class="search-container">
			<input name="search" type="text" placeholder="Search..." v-model="searchInput" @keyup.enter="load" @keyup.esc="clearSearch"/>
			<i class="fas fa-search" @click="load"></i>
			<i class="fas fa-times clear-search" v-if="isSearchResults" @click="clearSearch"></i>
		</div>
		<div class="stream-list" v-if="!loading">
			<router-link :to="'/streams/'+stream.ID" v-for="stream in filteredStreams" :key="stream.ID" class="stream">
				<div class="icon" :style="iconStyling(stream)">
					<i v-if="'Icon' in stream && stream.Icon==''" :class="'fas fa-'+stream.Icon"></i>
					<i v-else class="fas fa-bolt"></i>
				</div>
				<div class="info">
					<div class="name">{{stream.Name}}</div>
					<div class="cluster">{{stream.Cluster}}</div>
				</div>
				<router-link :to="'/streams/'+stream.ID+'/settings'"><i class="fas fa-cog config"></i></router-link>
			</router-link>
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
			isSearchResults: false,
			streams: []
		}
	},
	methods: {
		clearSearch() {
			this.searchInput = "";
			this.load();
		},
		iconStyling(stream) {
			var styling = {}
			if ("Color" in stream) {
				styling.background = stream.Color
			}
			return styling
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
				this.error = "Failed to "+((this.searchInput)?"search":"load")+" streams :(";
			}).finally(() => {
				if (this.searchInput != "") {
					this.isSearchResults = true;
				}
				this.loading = false;
			})
		},
		buildList(streams) {
			this.streams = streams.Streams;
		},
		retry() {
			if(this.searchInput != "") {
				this.search()
			} else {
				this.load()
			}
		}
	},
	computed: {
		filteredStreams() {
			if (this.searchInput == "") {
				return this.streams;
			}
			return _.filter(this.streams, stream => stream.Name.search(new RegExp(this.searchInput, "gi")) >= 0);
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
	z-index: 2;

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

.stream-list {
	padding-top: 70px;
	position: relative;
	width: 100%;

	.stream {
		padding: 10px 15px;
		width: calc(100% - 30px);
		position: relative;
		cursor: pointer;
		display: block;
		
		&:hover, &.router-link-active {
			background: -moz-linear-gradient(left, rgba(255,255,255,0.15) 0%, rgba(255,255,255,0.06) 100%);
			background: -webkit-linear-gradient(left, rgba(255,255,255,0.15) 0%,rgba(255,255,255,0.06) 100%);
			background: linear-gradient(to right, rgba(255,255,255,0.15) 0%,rgba(255,255,255,0.06) 100%);
			filter: progid:DXImageTransform.Microsoft.gradient( startColorstr='#26ffffff', endColorstr='#0fffffff',GradientType=1 );
		}

		.icon, .info {
			display: inline-block;
			vertical-align: middle;
		}

		.icon {
			height: 42px;
			width: 42px;
			background: #F4F8FB;
			border-radius: 3px;

			.fas {
				height: 100%;
				width: 100%;
				text-align: center;
				line-height: 42px;
				font-size: 16px;
				color: #1A1E30;
			}
		}

		.info {
			margin-left: 10px;

			.name {
				color: white;
				font-size: 12px;
				font-weight: 300;
			}

			.cluster {
				color: #C2C2C2;
				font-size: 9px;
				margin-top: -1px;
				text-transform: uppercase;
			}
		}
		
		.fas.config {
			position: absolute;
			top: 50%;
			transform: translateY(-50%);
			right: 15px;
			color: rgba(255,255,255,0.3);
			cursor: pointer;

			&:hover {
				color: white;
			}
		}
	}
}
</style>
