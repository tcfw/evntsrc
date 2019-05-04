<template>
  <div class="side-menu">
    <router-link to="/"><img src="@/assets/logo.png"/></router-link>
    <div class="menu">
      <router-link to="/dashboard"
        ><i class="fas fa-tachometer-alt"></i> Dashboard</router-link
      >
      <router-link to="/streams"
        ><i class="fas fa-stream"></i> Streams</router-link
      >
    </div>
    <div class="info">
      &copy; 2018 EvntSrc.io
      <p>
        <router-link to="/about">About</router-link> |
        <router-link to="/terms">Terms</router-link> |
        <router-link to="/privacy">Privacy</router-link> |
        <router-link to="/help">Help</router-link>
      </p>
    </div>
  </div>
</template>
<script>
export default {};
</script>
<style lang="scss" scoped>
.side-menu {
  background: #1a1e30;
  height: 100%;
  box-shadow: 2px 2px 4px 0 rgba(0, 0, 0, 0.1);
  z-index: 100;
  position: absolute;
  width: 220px;

  img {
    height: 18px;
    padding: 25px 15px 20px;
  }

  .menu {
    a {
      color: white;
      padding: 15px;
      font-weight: 100;
      font-size: 12px;
      text-decoration: none;
      display: block;

      .fas {
        padding: 0px 10px;
        display: inline-block;
      }

      &.router-link-active {
        background: rgba(255, 255, 255, 0.12);
        border-left: 4px solid #636d97;
        padding-left: 11px;
      }
    }
  }

  .info {
    color: white;
    position: absolute;
    bottom: 5px;
    left: 15px;
    font-size: 10px;

    a {
      color: #eee;

      &:hover {
        text-decoration: underline;
      }
    }
  }
}
</style>
