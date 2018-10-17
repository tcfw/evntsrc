<template>
	<el-dialog :visible.sync="visible" ref="ccDialog" width="30%" title="Payment Method" :show-close="false">
		<label style="margin-top: -25px; margin-bottom: 15px; display: block;">Credit or debit card</label>
		<el-card>
			<div id="stripeElement"></div>
		</el-card>
		<span slot="footer" class="dialog-footer">
			<el-button @click="visible = false" size="small">Cancel</el-button>
			<el-button type="primary" @click="submit" size="small" :loading="submitting">Attach card</el-button>
		</span>
	</el-dialog>	
</template>
<script>
export default {
  name: "CreditCardDialog",
  props: {
    visible: Boolean,
    onSubmit: Function
  },
  data() {
    return {
      submitting: false,
      stripeElement: null
    };
  },
  mounted() {
    this.$refs.ccDialog.$on("opened", this.loadStripeElement);
  },
  methods: {
    loadStripeElement() {
      if (!this.$root.stripe) {
        this.$root.loadStripe();
      }
      this.stripeElement = this.$root.stripe.elements().create("card", {});
      this.stripeElement.mount("#stripeElement");
    },
    submit() {
      this.submitting = true;
      this.visible = true;
      this.$root.stripe.createToken(this.stripeElement).then(result => {
        if (result.error) {
          this.$message({
            message: result.error.message,
            type: "error"
          });
          this.submitting = false;
        } else {
          let updateRequest = {
            cardToken: result.token.id,
            userId: this.$root.me.id
          };
          axios.post(this.$config.API + "/billing/user/" + this.$root.me.id + "/method", updateRequest)
          .then(() => {
            this.stripeElement.clear();
            this.visible = false;
            this.submitting = false;
            this.$message({
              message: "Payment method updated.",
              type: "success"
            });
            this.fetchInfo();
          })
          .catch(() => {
            this.submitting = false;
            this.$message({
              message: "Failed to update payment method. Please try again",
              type: "error"
            });
          });
        }
      });
    }
  }
};
</script>
<style lang="scss" scoped>
.plan {
  border-right: 1px solid #ddd;
  margin-left: -15px;
  margin-right: 15px;

  &:last-child {
    border-right: none;
    margin-right: 0px;
  }
}
</style>
