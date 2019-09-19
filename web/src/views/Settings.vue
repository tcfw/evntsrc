<template>
  <div class="w-full h-full absolute flex">
    <settings-side-menu style="width: 220px"></settings-side-menu>
    <div class="flex w-full h-full">
      <div class="w-full">
        <router-view></router-view>
      </div>
    </div>
  </div>
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
