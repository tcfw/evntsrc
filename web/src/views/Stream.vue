<template>
  <div class="w-full h-full">
    <template v-if="loaded">
      <div class="stream-tabs">
        <router-link disabled :to="'/streams/'+id+'/metrics'">Metrics</router-link>
        <router-link :to="'/streams/'+id+'/explore'">Explore</router-link>
        <router-link disabled :to="'/streams/'+id+'/webhooks'">Webhooks</router-link>
        <router-link :to="'/streams/'+id+'/settings'">Settings</router-link>
        <router-link :to="'/streams/'+id+'/keys'">API Keys</router-link>
      </div>
      <router-view></router-view>
    </template>
    <err404 v-if="notfound && !loaded"></err404>
  </div>
</template>
<script>
import err404 from '@/views/errors/404.vue'

export default {
  name: "stream",
  components: {
    err404
  },
  props: {
    id: String
  },
  data() {
    return {
      stream: null,
      loaded: false,
      notfound: false,
    }
  },
  watch: {
    id() {
      this.load();
    }
  },
  methods: {
    load() {
      this.loaded = false;
      this.notfound = false;
      this.$http.get(this.$config.API + '/stream/'+this.id).then(d => {
        this.stream = d.data
        this.loaded = true;
      }).catch(err => {
        if (err.response.status == 404) {
          this.notfound = true;
        }
      })
    }
  },
  mounted() {
    this.load();
  }
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
