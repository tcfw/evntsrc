<template>
	<span :class="classes">
        <img :src="src" v-if="src">
        <i :class="customIcon" v-else-if="icon || customIcon"></i>
        <span ref="children" :class="[prefixCls + '-string']" :style="childrenStyle" v-else><slot></slot></span>
    </span>
</template>
<script>
const prefixCls = 'avatar';

export default {
	name: "Avatar",
	props: {
		shape: {
			validator (value) {
				return ['circle', 'square'].includes(value);
			},
			default: 'circle'
		},
		size: {
			validator (value) {
				return ['small', 'large', 'default'].includes(value);
			},
			default: 'default'
		},
		src: {
			type: String
		},
		icon: {
			type: String
		},
		customIcon: {
			type: String,
			default: ''
		},
	},
	data () {
		return {
			prefixCls: prefixCls,
			scale: 1,
			childrenWidth: 0,
			isSlotShow: false
		};
	},
	computed: {
		classes () {
			return [
				`${prefixCls}`,
				`${prefixCls}-${this.shape}`,
				`${prefixCls}-${this.size}`,
				{
					[`${prefixCls}-image`]: !!this.src,
					[`${prefixCls}-icon`]: !!this.icon || !!this.customIcon
				}
			];
		},
		childrenStyle () {
			let style = {};
			if (this.isSlotShow) {
				style = {
					msTransform: `scale(${this.scale})`,
					WebkitTransform: `scale(${this.scale})`,
					transform: `scale(${this.scale})`,
					position: 'absolute',
					display: 'inline-block',
					left: `calc(50% - ${Math.round(this.childrenWidth / 2)}px)`
				};
			}
			return style;
		}
	},
	methods: {
		setScale () {
			this.isSlotShow = !this.src && !this.icon;
			if (this.$refs.children) {
				// set children width again to make slot centered
				this.childrenWidth = this.$refs.children.offsetWidth;
				const avatarWidth = this.$el.getBoundingClientRect().width;
				// add 4px gap for each side to get better performance
				if (avatarWidth - 8 < this.childrenWidth) {
					this.scale = (avatarWidth - 8) / this.childrenWidth;
				} else {
					this.scale = 1;
				}
			}
		}
	},
	mounted () {
		this.setScale();
	},
	updated () {
		this.setScale();
	}
}
</script>
<style lang="scss" scoped>
.avatar {
	display: inline-block;
    text-align: center;
    background: #888;
    color: #fff;
    white-space: nowrap;
    position: relative;
    overflow: hidden;
	vertical-align: middle;
	height: 55px;
	width: 55px;
	line-height: 55px;
	border-radius: 28px;
	& > * {
        line-height: 55px;
    }
	
	&-image{
        background: transparent;
	}
	
	& > img {
        width: 100%;
        height: 100%;
	}
	
	&-small {
		height: 35px;
		width: 35px;
		line-height: 35px;
	}
}
</style>
