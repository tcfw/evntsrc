<template>
  <div class="flex h-full">
    <div class="w-64">
      <stream-side-menu ref="sidemenu"></stream-side-menu>
    </div>
    <div class="w-full">
      <router-view></router-view>
    </div>
  </div>
</template>
<script>
import streamSideMenu from "@/components/StreamSideMenu.vue";

export default {
  name: "streams",
  components: {
    streamSideMenu
  },
  methods: {
    reload() {
      this.$refs.sidemenu.load();
    }
  }
};
</script>
