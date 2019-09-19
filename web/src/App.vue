<template>
  <div id="app" ref="App" :class="appClass">
    <div class="flex h-full">
      <div class="w-16" width="" v-if="$root.loggedIn()">
        <side-menu></side-menu>
      </div>
      <div>
        <div :style="
            $root.loggedIn()
              ? {
                  position: 'absolute',
                  top: '0px',
                  width: 'calc(100% - 4rem)',
                  height: '100%'
                }
              : { zIndex: 1 }
          ">
          <router-view></router-view>
        </div>
      </div>
    </div>
    <div class="c-info">
      &copy; {{year}} EvntSrc.io &nbsp;&nbsp;
        <router-link to="/terms">Terms</router-link> &nbsp;|&nbsp;
        <router-link to="/privacy">Privacy</router-link> &nbsp;|&nbsp;
        <router-link to="/help">Help</router-link>
    </div>
  </div>
</template>
<script>
import pageHeader from "@/components/Header.vue";
import sideMenu from "@/components/SideMenu.vue";

export default {
  name: "app",
  components: {
    pageHeader,
    sideMenu
  },
  data() {
    return {
      appClass: ""
    };
  },
  created() {
    this.$root.$refs.App = this;
  },
  computed: {
    year() {
      return (new Date).getFullYear();
    }
  }
};
</script>
<style lang="scss">
.c-info {
  @apply absolute bottom-0 right-0 mr-4 mb-4 text-xs;
  a {
    &:hover {
      text-decoration: underline;
    }
  }
}
</style>
