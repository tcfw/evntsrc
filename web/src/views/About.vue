<template>
  <div class="about">
    <el-breadcrumb separator-class="el-icon-arrow-right">
      <el-breadcrumb-item :to="{ path: '/' }">Home</el-breadcrumb-item>
      <el-breadcrumb-item>About</el-breadcrumb-item>
    </el-breadcrumb>
    <h1>This is an about page</h1>
  </div>
</template>
<style lang="scss" scoped>
.about {
  padding: 15px;
}
</style>
