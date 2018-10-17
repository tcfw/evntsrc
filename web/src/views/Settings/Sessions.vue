<template>
    <div style="padding: 15px;">
		<el-breadcrumb separator-class="el-icon-arrow-right">
			<el-breadcrumb-item :to="{ path: '/' }">Home</el-breadcrumb-item>
			<el-breadcrumb-item :to="{ path: '/settings' }">Settings</el-breadcrumb-item>
			<el-breadcrumb-item>Sessions</el-breadcrumb-item>
		</el-breadcrumb>
	</div>
</template>
