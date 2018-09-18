<template>
  <div id="app">
    <el-container :style="{height: '100%'}">
      <el-aside width="220px"><side-menu></side-menu></el-aside>
      <el-container>
        <el-header>
          <page-header></page-header>
        </el-header>
        <el-main><router-view></router-view></el-main>
      </el-container>
    </el-container>

  </div>
</template>
<script>
import "@/styles/app.scss";

import pageHeader from "@/components/Header.vue"
import sideMenu from "@/components/SideMenu.vue"


export default {
  name: 'app',
  components: {
    pageHeader,
    sideMenu
  }
}
</script>
