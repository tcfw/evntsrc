<template>
    <div style="padding: 15px;">
		<el-breadcrumb separator-class="el-icon-arrow-right">
			<el-breadcrumb-item :to="{ path: '/' }">Home</el-breadcrumb-item>
			<el-breadcrumb-item :to="{ path: '/settings' }">Settings</el-breadcrumb-item>
			<el-breadcrumb-item>Billing</el-breadcrumb-item>
		</el-breadcrumb>
		<div v-if="loading == false && plansLoading == false" style="height: calc(100% - 20px)">
			<h3>Subscription</h3>
			<div v-if="info.subscriptions && info.subscriptions.length > 0">
				<div v-for="subscription in info.subscriptions" :key="subscription.id">
					<el-card class="box-card">
						<el-row>
							<el-col :span="10" style="text-align: center; padding-right: 20px;">
								<h3 style="margin: 0">
									{{getPlan(subscription.plan).product.name}}
									<div style="color:#aaa;font-size:12px">Plan</div>
									<div v-if="getPlan(subscription.plan).product.metadata && getPlan(subscription.plan).product.metadata.description" style="margin-top: 15px; font-size: 14px; text-align: left;">
										<span v-html="getPlan(subscription.plan).product.metadata.description"></span>
									</div>
								</h3>
							</el-col>
							<el-col :span="14" style="border-left: 1px solid #ddd; padding-left: 25px;">
								${{getPlan(subscription.plan).amount | formatAmount}} billed <b>{{getPlan(subscription.plan).nickname.toLowerCase()}}</b>
								<div style="font-size:10px">Next bill on {{formatTime(subscription.currentPeriodEnds, false)}}</div>
								<div v-if="subscription.discount" style="font-size: 12px; margin-top: 5px;">
									<i class="fas fa-percentage"></i> Discount - <i>{{subscription.discount}}</i>
									<span v-if="subscription.discountEnds"> until <b>{{formatTime(subscription.discountEnds)}}</b></span>
									<b v-else-if="!subscription.discountEnds && subscription.discount.substr(0,1) != '$'"> &nbsp;forever! <i class="fas fa-award"></i></b>
									<span v-else> off</span>
								</div>
								<div v-if="subscription.status == 'trialing'" style="margin-top: 15px;">
									<el-row>
										<el-col :span="1">
											<i class="fas fa-clock" style="color: orange;"></i> 
										</el-col>
										<el-col :span="22">
											Currently on trial<br/><small>Trial ends {{timeUntil(subscription.trialEnds)}} on <b>{{formatTime(subscription.trialEnds)}}</b></small>
										</el-col>
									</el-row>
								</div>
								<div style="margin-top: 15px;"><el-button size="small" @click="showPlans = true">Change</el-button></div>
							</el-col>
						</el-row>
					</el-card>
				</div>
			</div>
			<div v-if="!info.subscriptions || info.subscriptions.length == 0">
				<i>You currently do not have any subscriptions.</i>
			</div>
			<h3>Payment Info</h3>
			<div v-if="info.sources && info.sources.length > 0">
				<el-card class="box-card" style="margin-bottom: 15px;">
					<div v-if="info.sources[0].type == 'card'">
						<el-row>
							<el-col :span="4" style="text-align: center; padding-right: 20px;">
								<b>{{info.sources[0].card.brand}}</b>
								<div style="margin-top: 15px;"><i :class="'fab fa-3x fa-cc-'+info.sources[0].card.brand.toLowerCase()"></i></div>
							</el-col>
							<el-col :span="20" style="border-left: 1px solid #ddd; padding-left: 25px;">
								<small style="color: #aaa; font-size: 12px;">Number</small> &nbsp;•••• {{info.sources[0].card.number}}<br/>
								<small style="color: #aaa; font-size: 12px;">Expiry</small> &nbsp;{{(info.sources[0].card.expMonth).toString().padStart(2,"0")}} / {{info.sources[0].card.expYear}}
								<div v-if="isExpiringSoon" style="color:red; font-size: 13px;"><i class="fas fa-exclamation"></i> &nbsp;Expiring soon</div>
								<div style="margin-top: 15px;" @click="showCCForm = true"><el-button size="small">Update</el-button></div>
							</el-col>
						</el-row>
					</div>
				</el-card>
			</div>
			<div v-if="info.accountBalance">
				<el-card class="box-card">
					<el-row>
						<el-col :span="4" style="text-align: center; padding-right: 20px;">
							Balance
						</el-col>
						<el-col :span="20" style="border-left: 1px solid #ddd; padding-left: 25px;">
							<b>$ {{Math.abs(info.accountBalance) | formatAmount}}</b> <span v-if="info.accountBalance.substr(0,1) =='-'">&nbsp;credit</span><span v-else>&nbsp;owing</span>
							<div style="font-size: 12px; color: #aaa;">This amount will be applied to your next bill</div>
						</el-col>
					</el-row>
				</el-card>
			</div>
			<div v-if="!info.sources && !info.accountBalance">
				<el-card class="box-card">
					You have no payment methods attached
				</el-card>
			</div>
		</div>
		<div v-else>
			<div style="padding-top: 50px; text-align: center">
				<i class="fas fa-sync fa-spin" style="color: #ddd"></i><br />
				<small>Loading...</small>
			</div>
		</div>
		<el-dialog :visible.sync="showCCForm" ref="ccDialog" width="30%" title="Payment Method" :show-close="false">
			<label style="margin-top: -25px; margin-bottom: 15px; display: block;">Credit or debit card</label>
			<el-card>
				<div id="stripeElement"></div>
			</el-card>
			<span slot="footer" class="dialog-footer">
				<el-button @click="showCCForm = false" size="small">Cancel</el-button>
    			<el-button type="primary" @click="getCCToken" size="small" :loading="submittingCC">Attach card</el-button>
			</span>
		</el-dialog>
		<el-dialog :visible.sync="showPlans" width="75%">
			<el-row>
				<el-col v-for="product in products()" :key="product.id" :span="Math.round(24 / Object.keys(products()).length)" class="plan">
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
				<el-button @click="showPlans = false" size="small">Cancel</el-button>
    			<el-button type="primary" @click="updateSubscription" size="small" :loading="submittingUpdateSubscription">Change</el-button>
			</span>
		</el-dialog>
	</div>
</template>
<script>
import billing_pb from "@/protos/billing_pb.js";
const moment = require("moment");

export default {
	name: "billing",
	data () {
		return {
			loading: true,
			plansLoading: true,
			info: false,
			plans: false,
			showCCForm: false,
			submittingCC: false,
			stripeElement: null,
			showPlans: false,
			submittingUpdateSubscription: false,
			changeSubscription: null,
		}
	},
	mounted() {
		this.fetchPlans();
		if(this.$root.me.id == undefined) {
			this.$root.$on("me.ready", this.fetchInfo);
		} else {
			this.fetchInfo();
		}
		this.$refs.ccDialog.$on('opened', this.loadStripeCC);
	},
	methods: {
		fetchInfo(me) {
			axios.get(this.$config.API+"/billing/user/"+this.$root.me.id).then(d => {
				this.loading = false;
				this.info = d.data;
				this.changeSubscription = this.info.subscriptions[0].plan
			}).catch(e => {
				this.loading = false;
				this.$message({
					message: "Failed to fetch your billing info. Please try again.",
					type: "error"
				});
			});
		},
		fetchPlans() {
			axios.get(this.$config.API+"/billing/plans").then(d => {
				this.plans = d.data.plans;
				this.plansLoading = false;
			});
		},
		getPlan(planId) {
			if(this.plansLoading) return {name: ""};
			var plan = _.find(this.plans, {'id':planId})
			return plan;
		},
		timeUntil(time) {
			return moment().to(moment.unix(time));
		},
		formatTime(time, includeTime) {
			return moment.unix(time).format("LL")
		},
		loadStripeCC() {
			if (!this.$root.stripe) {
				this.$root.loadStripe();
			}
			this.stripeElement = this.$root.stripe.elements().create('card', {});
			this.stripeElement.mount("#stripeElement")
		},
		getCCToken(done) {
			this.submittingCC = true;
			this.showCCForm = true;
			this.$root.stripe.createToken(this.stripeElement).then(result => {
				if(result.error) {
					this.$message({
						message: result.error.message,
						type: "error"
					});
					this.submittingCC = false;
				} else {
					let updateRequest = {
						cardToken: result.token.id, 
						userId: this.$root.me.id
					};
					axios.post(this.$config.API+"/billing/user/"+this.$root.me.id+"/method", updateRequest).then(d => {
						this.stripeElement.clear();
						this.showCCForm = false;
						this.submittingCC = false;
						this.$message({
							message: "Payment method updated.",
							type: "success"
						});
						this.fetchInfo();
					}).catch(e => {
						this.submittingCC = false;
						this.$message({
							message: "Failed to update payment method. Please try again",
							type: "error"
						});
					})
					
				}
			})
		},
		updateSubscription() {
			this.submittingUpdateSubscription = true;
		},
		orderByKey(arr, key) {
			return arr.sort((a,b) => {
				if (a[key] > b[key]) return 1
				if (a[key] < b[key]) return -1
				return 0
			})
		},
		products() {
			if(this.loading == true || this.plansLoading == true) return [];
			var currentPlan = _.find(this.plans, {id: this.info.subscriptions[0].plan})
			return _.sortBy(_.reduce(this.plans, (carry, plan) => {
				if((currentPlan.interval == 'year' && plan.interval == 'year') || currentPlan.interval == 'month') {
					(carry[plan.product.id] || (carry[plan.product.id] = Object.assign({plans: []},plan.product))).plans.push(plan);
				}
				return carry;
			}, {}), ['metadata.order']);
		}
	},
	filters: {
		formatAmount(value) {
			var decValue = parseFloat(value) / 100;
			return decValue.toFixed(2).replace(/\d(?=(\d{3})+\.)/g, '$&,');
		}
	},
	computed: {
		isExpiringSoon() {
			if (this.info.sources.length == 0) return false;
			var exprDate = moment(this.info.sources[0].card.expMonth + "-" +this.info.sources[0].card.expYear, ["MM-YY","MM-YYYY"])
			return moment().add(2, "months").isSameOrAfter(exprDate)
		}
	}
}
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

