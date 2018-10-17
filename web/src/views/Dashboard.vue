<template>
  <div class="dashboard">
    <el-breadcrumb separator-class="el-icon-arrow-right">
      <el-breadcrumb-item :to="{ path: '/' }">Home</el-breadcrumb-item>
      <el-breadcrumb-item>Dashboard</el-breadcrumb-item>
    </el-breadcrumb>
  </div>
</template>

<script>
export default {
  name: "home",
  components: {}
};
</script>
<style lang="scss">
.dashboard {
  padding: 15px;
}
</style>
