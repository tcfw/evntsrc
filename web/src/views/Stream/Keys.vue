<template>
  <div class="panel">
    <div class="d-head rounded-t clearfix">
      <div class="float-right">
        <button>Create key</button>
      </div>
    </div>
    <table class="d-table">
      <thead>
        <tr>
          <th>Label</th>
          <th>Key</th>
          <th>Secret</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="key in keys" :key="key.id">
          <td class="w-2/3">{{ key.label }}</td>
          <td>{{ key.key }}</td>
          <td>•••••••••••••</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
<script>
export default {
  name: "stream-auth",
  data() {
    return {
      loading: true,
      keys: []
    };
  },
  created() {
    this.$http.get(this.$config.API + "/stream/" + this.$route.params.id + "/keys")
      .then(d => {
        this.loading = false;
        this.keys = d.data.keys;
      });
  }
};
</script>