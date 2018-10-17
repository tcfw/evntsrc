<template>
    <div class="settings-side-menu">
        <div class="header">
            Account Settings
        </div>
        <div class="links">
            <router-link :to="{name: 'account'}">Account</router-link>
            <router-link to="/settings/billing">Billing</router-link>
            <router-link to="/settings/security">Security</router-link>
            <router-link to="/settings/sessions">Sessions</router-link>
        </div>
    </div>
</template>
<script>
export default {
  name: "settings-side-menu"
};
</script>
<style lang="scss" scoped>
.settings-side-menu {
  background: white;
  box-shadow: 2px 0 4px 0 rgba(0, 0, 0, 0.1);
  position: relative;
  height: 100%;
  background: #50566f;

  .header {
    box-shadow: inset 0 1px 3px 0 rgba(0, 0, 0, 0.1);
    width: calc(100% - 20px);
    color: white;
    padding: 10px;
    margin-top: -1px;
    font-size: 12px;
  }

  .links {
    margin-top: 15px;

    a {
      display: block;
      font-size: 13px;
      padding: 15px;
      color: #ddd;

      &.router-link-active {
        background: lighten(#50566f, 10%);
        color: white;
        font-weight: bold;
      }

      &:hover {
        background: lighten(#50566f, 20%);
        color: white;
      }
    }
  }
}
</style>
