<template>
	<div :style="{position: 'relative'}" :class="containerClass">
		<div id="injector">
			<small class="injector-header">&nbsp;</small>
			<el-form :model="injector" label-width="120px">
				<el-form-item label="Channel">
					<el-input v-model="injector.channel"></el-input>
				</el-form-item>
				<el-form-item label="Event type">
					<el-input v-model="injector.type"></el-input>
				</el-form-item>
				<el-form-item label="Data">
					<el-input v-model="injector.data"></el-input>
				</el-form-item>
				<div :style="{textAlign: 'right'}">
					<el-button type="primary" size="small">Inject</el-button>
				</div>
			</el-form>
		</div>
		<el-button id="inject-btn" type="info" @click="toggleInjector" size="small"><span v-show="!injectorOpen">Injector</span> <i class="fas fa-times" v-show="injectorOpen"></i></el-button>
	</div>
</template>
<script>
export default {
  name: "stream-history",
  data() {
    return {
      injectorOpen: false,
      injector: {
        channel: "",
        type: "",
        data: ""
      }
    };
  },
  computed: {
    containerClass() {
      return {
        "injector-open": this.injectorOpen
      };
    }
  },
  methods: {
    toggleInjector() {
      this.injectorOpen = !this.injectorOpen;
    }
  }
};
</script>
<style lang="scss" scoped>
.injector-header {
  color: #888;
  margin-top: 5px;
  display: block;
}
#inject-btn {
  position: absolute;
  top: 0px;
  right: 15px;
  border-top-left-radius: 0px;
  border-top-right-radius: 0px;
}

#injector {
  overflow: hidden;
  height: 0px;
  width: calc(100% - 60px);
  margin: 0 15px;
  padding: 0 15px;
  background: white;
  box-shadow: 0px 2px 2px rgba(0, 0, 0, 0.1);
  border-bottom-left-radius: 3px;
}

$openHeight: 260px;

.injector-open {
  #inject-btn {
    top: $openHeight;
  }

  #injector {
    height: $openHeight;
  }
}
</style>
