<template>
  <el-container
    :style="{ height: '100%', width: '100%', position: 'absolute' }"
  >
    <el-aside width="220px">
      <settings-side-menu></settings-side-menu>
    </el-aside>
    <el-container>
      <el-main>
        <router-view></router-view>
      </el-main>
    </el-container>
  </el-container>
</template>
<script>
import settingsSideMenu from "@/components/SettingsSideMenu.vue";

export default {
  name: "settings",
  components: {
    settingsSideMenu
  }
};
</script>
