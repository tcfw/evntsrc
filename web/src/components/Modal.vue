<template>
	<div class="modal" :class="{block: visible, sm : size == 'sm', md: size == 'md', lg: size =='lg'}" @keyup.esc="close" tabindex=0>
		<div class="container" :style="{width: width}" :class="containerClass">
			<div class="title">
				<div v-if="closable" class="close" @click="close"><i class="fas fa-times"></i></div>
				<slot name="title"></slot>
			</div>
			<div class="body"><slot name="body"></slot></div>
		</div>
	</div>
</template>
<script>
export default {
	name: "modal",
	props: {
		closable: Boolean,
		width: String,
		size: String,
		position: String,
	},
	data() {
		return {
			visible: false,
		}
	},
	computed: {
		widthClass() {
			if (this.width && this.width.indexOf("-") >=0) {
				return this.width;
			}

			return "";
		},
		positionClass() {
			if (this.position) {
				return "pos-" + this.position
			}
			return "" 
		},
		containerClass() {
			return [this.widthClass, this.positionClass];
		}
	},
	methods: {
		close() {
			this.visible = false;
			this.$emit("hidden");
		},
		show() {
			this.visible = true;
			this.$emit("shown");
		}
	},
	mounted() {
		document.body.appendChild(this.$el);
	}
}
</script>
<style lang="scss" scoped>
.modal {
	@apply absolute hidden rounded top-0 left-0 w-full h-full;

	&.sm .container {
		width: 24rem;
	}

	&.md .container {
		width: 34rem;
	}

	&.lg .container {
		width: 56rem;
	}

	.container {
		@apply absolute rounded shadow-xl bg-white max-w-full;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%);
		z-index: 9999;

		.title {
			@apply p-3 bg-ev-100 text-white rounded-t relative;

			.close {
				@apply absolute text-white cursor-pointer font-bold;
				top: 50%;
    			right: 0.75rem;
    			transform: translate(-50%, -50%);
			}
		}

		.body {
			@apply p-3 rounded-b;
		}

		&.pos {
			&-top {
				top: 20%;
				transform: translateX(-50%);
			}
		}
	}

	&:after {
		content: "";
		@apply block absolute top-0 left-0 w-full h-full;
		background: rgba(0,0,0,0.1);
		z-index: 9998;
	}

	&.block {
		display: block !important;
	}
}
</style>
