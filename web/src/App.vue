<template>
  <div id="app" ref="App" :class="appClass">
    <el-container :style="{height: '100%'}">
      <el-aside width="220px" v-if="$root.loggedIn()"><side-menu></side-menu></el-aside>
      <el-container>
        <el-header v-if="$root.loggedIn()" :style="{zIndex: 2}">
          <page-header></page-header>
        </el-header>
        <el-main :style="$root.loggedIn() ? {position: 'absolute', top: '60px', width: 'calc(100% - 220px)', height: 'calc(100% - 59px)'} : {zIndex: 1}"><router-view></router-view></el-main>
      </el-container>
    </el-container>

  </div>
</template>
<script>
import pageHeader from "@/components/Header.vue"
import sideMenu from "@/components/SideMenu.vue"


export default {
  name: 'app',
  components: {
    pageHeader,
    sideMenu
  },
  data() {
    return {
      appClass: ""
    }
  },
  created() {
    this.$root.$refs.App = this;
  }
}
</script>
<style lang="scss">
a {
  color: #4A90E2;
  text-decoration: none;
}
</style>