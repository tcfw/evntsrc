<template>
  <div class="dashboard">
    <el-breadcrumb separator-class="el-icon-arrow-right">
      <el-breadcrumb-item :to="{ path: '/' }">Home</el-breadcrumb-item>
      <el-breadcrumb-item>Dashboard</el-breadcrumb-item>
    </el-breadcrumb>
    <modal :closable="true" width="w-2/3">
      <div slot="title">Header</div>
      <div slot="body">Body</div>
    </modal>
  </div>
</template>

<script>
import modal from '@/components/Modal.vue'

export default {
  name: "home",
  components: {
    modal
  }
};
</script>
<style lang="scss">
.dashboard {
  padding: 15px;
}
</style>
