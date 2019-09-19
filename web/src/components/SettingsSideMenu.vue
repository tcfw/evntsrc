<template>
  <div class="settings-side-menu">
    <div class="header">
      Account Settings
    </div>
    <div class="links">
      <router-link :to="{ name: 'account' }">Account</router-link>
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
  @apply h-full relative shadow-lg bg-ev-999;

  .header {
    @apply text-text-100 p-3 text-xs;
  }

  .links {
    margin-top: 15px;

    a {
      @apply block text-text-100 block text-sm p-4;

      &.router-link-active {
        @apply text-white;
        background: lighten(#50566f, 10%);
      }

      &:hover {
        @apply text-white;
        background: lighten(#50566f, 20%);
      }
    }
  }
}
</style>
