<template>
  <div class="page-header">
    <div id="profile">
      <el-dropdown trigger="click" size="medium">
        <span class="el-dropdown-link">
          <Avatar size="small">{{ $root.me.name }}</Avatar> &nbsp;<i
            class="el-icon-arrow-down el-icon--right"
          ></i>
        </span>
        <el-dropdown-menu slot="dropdown">
          <el-dropdown-item disabled>{{ $root.me.name }}</el-dropdown-item>
          <router-link :to="{ name: 'settings' }"
            ><el-dropdown-item>Settings</el-dropdown-item></router-link
          >
          <el-dropdown-item @click.native="$root.logout()" divided
            >Logout</el-dropdown-item
          >
        </el-dropdown-menu>
      </el-dropdown>
    </div>
  </div>
</template>
<script>
import Avatar from "@/components/Avatar.vue";

export default {
  name: "page-header",
  components: {
    Avatar
  }
};
</script>
<style lang="scss" scoped>
.page-header {
  position: relative;
  height: 60px;

  #profile {
    position: absolute;
    top: 50%;
    right: 10px;
    transform: translateY(-50%);
  }
}
</style>
