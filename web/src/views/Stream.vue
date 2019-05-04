<template>
  <div v-loading.lock="loading" :style="{ height: '100%', width: '100%' }">
    <el-menu
      ref="menu"
      class="stream-tabs"
      mode="horizontal"
      :default-active="activeIndex"
      menu-trigger="click"
      @select="menuSelect"
    >
      <el-menu-item disabled index="stats">Stats</el-menu-item>
      <el-menu-item index="history">History</el-menu-item>
      <el-menu-item index="auth">API Keys</el-menu-item>
      <el-menu-item disabled index="ingress">Ingress</el-menu-item>
      <el-menu-item disabled index="webhooks">Webhooks</el-menu-item>
      <el-menu-item index="settings">Settings</el-menu-item>
    </el-menu>
    <router-view></router-view>
  </div>
</template>
<script>
export default {
  name: "stream",
  props: {
    id: String
  },
  data() {
    return {
      loading: false,
      error: "",
      stream: null,
      activeIndex: null
    };
  },
  watch: {
    $route: "setActiveMenu"
  },
  methods: {
    load() {
      // this.loading = false;
    },
    setActiveMenu() {
      this.activeIndex =
        this.$route.matched.length >= 3
          ? this.$route.matched[2].name.substr("stream-".length)
          : "";
      this.$refs.menu.activeIndex = this.activeIndex;
    },
    menuSelect(key) {
      this.$router.push({ name: "stream-" + key, params: { id: this.id } });
    }
  },
  mounted() {
    this.load();
    this.setActiveMenu();
  }
};
</script>
<style lang="scss">
.stream-tabs {
  padding: 6px 0px 0px 15px;
  background: white;
  box-shadow: 0 2px 3px 0 rgba(0, 0, 0, 0.1);

  .el-menu-item {
    font-size: 13px;
    font-weight: 300;
    color: #50566f;
    padding-bottom: 45px;

    &.is-active {
      font-weight: 400;
      padding-right: 18px;
    }
  }
}
</style>
