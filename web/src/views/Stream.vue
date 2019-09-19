<template>
  <div class="w-full h-full">
    <div class="stream-tabs">
      <router-link disabled :to="'/streams/'+id+'/metrics'">Metrics</router-link>
      <router-link :to="'/streams/'+id+'/explore'">Explore</router-link>
      <router-link disabled :to="'/streams/'+id+'/webhooks'">Webhooks</router-link>
      <router-link :to="'/streams/'+id+'/settings'">Settings</router-link>
      <router-link :to="'/streams/'+id+'/keys'">API Keys</router-link>
    </div>
    <router-view></router-view>
  </div>
</template>
<script>
export default {
  name: "stream",
  props: {
    id: String
  },
};
</script>
<style lang="scss">
.stream-tabs {
  @apply py-5 px-3 text-sm;
  font-weight: 300;
  color: #50566f;

  a {
    @apply px-4 py-2 rounded;

    &.router-link-active {
      @apply text-white bg-ev-100;
    }
  }

  .is-active {
    font-weight: 400;
    padding-right: 18px;
  }
}
</style>
