<template>
	<el-dialog :visible.sync="visible" width="75%">
		<el-row>
			<el-col v-for="product in products" :key="product.id" :span="Math.round(24 / Object.keys(products).length)" class="plan">
				<div style="text-align: center">
					<h3 style="margin: 0">
						{{product.name}}
						<div style="color:#aaa;font-size:12px">Plan</div>
					</h3>
				</div>
				<div v-html="product.metadata.description"></div>
				<div v-for="plan in orderByKey(product.plans, 'nickname')" :key="plan.id" style="padding-left: 20px;">
					<el-radio v-model="changeSubscription" :label="plan.id">
						${{plan.amount | formatAmount}} <small>{{plan.nickname}}</small>
					</el-radio>
				</div>
			</el-col>
		</el-row>
		<span slot="footer" class="dialog-footer">
			<el-button @click="visible = false" size="small">Cancel</el-button>
			<el-button type="primary" @click="update" size="small" :loading="submitting">Change</el-button>
		</span>
	</el-dialog>	
</template>
<script>
export default {
  name: "PlanDialog",
  props: {
    products: Object,
    visible: Boolean
  },
  data() {
    return {
      submitting: false
    };
  },
  methods: {
    update() {
      this.submitting = true;
    },
    orderByKey(arr, key) {
      return arr.sort((a, b) => {
        if (a[key] > b[key]) return 1;
        if (a[key] < b[key]) return -1;
        return 0;
      });
    }
  }
};
</script>
