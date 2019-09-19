<template>
  <div class="side-menu">
    <router-link to="/"><img src="@/assets/round-reverse.png"/></router-link>
    <div class="menu">
      <router-link title="Dashboard" to="/dashboard"><i class="fas fa-tachometer-alt"></i></router-link>
      <router-link title="Streams" to="/streams"><i class="fas fa-stream"></i></router-link>
    </div>
    <div class="avatar">
      <el-dropdown trigger="click" size="medium">
        <span class="el-dropdown-link">
          <Avatar size="small">{{ $root.me.name }}</Avatar>
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
  components: {
    Avatar
  }
};
</script>
<style lang="scss" scoped>
.side-menu {
  @apply absolute bg-ev-100 h-full shadow-lg w-16;

  img {
    padding: 1rem 0.75rem;
  }

  .menu {
    a {
      color: white;
      padding: 15px;
      font-weight: 100;
      font-size: 16px;
      text-decoration: none;
      display: block;

      .fas {
        padding: 0px 8px;
        display: inline-block;
      }

      &.router-link-active {
        @apply bg-ev-50 shadow-inner;
      }
    }
  }

  .avatar {
    color: white;
    position: absolute;
    bottom: 0.3rem;
    left: 0.5rem;
  }
}
</style>
