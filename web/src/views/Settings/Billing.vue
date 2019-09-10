<template>
  <div style="padding: 15px;">
    <el-breadcrumb separator-class="el-icon-arrow-right">
      <el-breadcrumb-item :to="{ path: '/' }">Home</el-breadcrumb-item>
      <el-breadcrumb-item :to="{ path: '/settings' }"
        >Settings</el-breadcrumb-item
      >
      <el-breadcrumb-item>Billing</el-breadcrumb-item>
    </el-breadcrumb>
    <div
      v-if="loading == false && plansLoading == false"
      style="height: calc(100% - 20px)"
    >
      <h3>Subscription</h3>
      <div v-if="info.subscriptions && info.subscriptions.length > 0">
        <div v-for="subscription in info.subscriptions" :key="subscription.id">
          <el-card class="box-card">
            <el-row>
              <el-col
                :span="10"
                style="text-align: center; padding-right: 20px;"
              >
                <h3 style="margin: 0">
                  {{ getPlan(subscription.plan).product.name }}
                  <div style="color:#aaa;font-size:12px">Plan</div>
                  <div
                    v-if="
                      getPlan(subscription.plan).product.metadata &&
                        getPlan(subscription.plan).product.metadata.description
                    "
                    style="margin-top: 15px; font-size: 14px; text-align: left;"
                  >
                    <span
                      v-html="
                        getPlan(subscription.plan).product.metadata.description
                      "
                    ></span>
                  </div>
                </h3>
              </el-col>
              <el-col
                :span="14"
                style="border-left: 1px solid #ddd; padding-left: 25px;"
              >
                ${{ getPlan(subscription.plan).amount | formatAmount }} billed
                <b>{{ getPlan(subscription.plan).nickname.toLowerCase() }}</b>
                <div style="font-size:10px">
                  Next bill on
                  {{ formatTime(subscription.currentPeriodEnds, false) }}
                </div>
                <div
                  v-if="subscription.discount"
                  style="font-size: 12px; margin-top: 5px;"
                >
                  <i class="fas fa-percentage"></i> Discount -
                  <i>{{ subscription.discount }}</i>
                  <span v-if="subscription.discountEnds">
                    until
                    <b>{{ formatTime(subscription.discountEnds) }}</b></span
                  >
                  <b
                    v-else-if="
                      !subscription.discountEnds &&
                        subscription.discount.substr(0, 1) != '$'
                    "
                  >
                    &nbsp;forever! <i class="fas fa-award"></i
                  ></b>
                  <span v-else> off</span>
                </div>
                <div
                  v-if="subscription.status == 'trialing'"
                  style="margin-top: 15px;"
                >
                  <el-row>
                    <el-col :span="1">
                      <i class="fas fa-clock" style="color: orange;"></i>
                    </el-col>
                    <el-col :span="22">
                      Currently on trial<br /><small
                        >Trial ends {{ timeUntil(subscription.trialEnds) }} on
                        <b>{{ formatTime(subscription.trialEnds) }}</b></small
                      >
                    </el-col>
                  </el-row>
                </div>
                <div style="margin-top: 15px;">
                  <el-button size="small" @click="showPlans = true"
                    >Change</el-button
                  >
                </div>
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
              <el-col
                :span="4"
                style="text-align: center; padding-right: 20px;"
              >
                <b>{{ info.sources[0].card.brand }}</b>
                <div style="margin-top: 15px;">
                  <i
                    :class="
                      'fab fa-3x fa-cc-' +
                        info.sources[0].card.brand.toLowerCase()
                    "
                  ></i>
                </div>
              </el-col>
              <el-col
                :span="20"
                style="border-left: 1px solid #ddd; padding-left: 25px;"
              >
                <small style="color: #aaa; font-size: 12px;">Number</small>
                &nbsp;•••• {{ info.sources[0].card.number }}<br />
                <small style="color: #aaa; font-size: 12px;">Expiry</small>
                &nbsp;{{
                  info.sources[0].card.expMonth.toString().padStart(2, "0")
                }}
                / {{ info.sources[0].card.expYear }}
                <div v-if="isExpiringSoon" style="color:red; font-size: 13px;">
                  <i class="fas fa-exclamation"></i> &nbsp;Expiring soon
                </div>
                <div style="margin-top: 15px;" @click="showCCForm = true">
                  <el-button size="small">Update</el-button>
                </div>
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
            <el-col
              :span="20"
              style="border-left: 1px solid #ddd; padding-left: 25px;"
            >
              <b>$ {{ Math.abs(info.accountBalance) | formatAmount }}</b>
              <span v-if="info.accountBalance.substr(0, 1) == '-'"
                >&nbsp;credit</span
              ><span v-else>&nbsp;owing</span>
              <div style="font-size: 12px; color: #aaa;">
                This amount will be applied to your next bill
              </div>
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
    <CreditCardDialog :visible="showCCForm"></CreditCardDialog>
    <PlanDialog :visible="showPlans" :products="products()"></PlanDialog>
  </div>
</template>
<script>
// import billing_pb from "@/protos/billing_pb.js";
import CreditCardDialog from "./Billing/CreditCardDialog.vue";
import PlanDialog from "./Billing/PlanDialog.vue";

const moment = require("moment");

export default {
  name: "billing",
  components: {
    CreditCardDialog,
    PlanDialog
  },
  data() {
    return {
      loading: true,
      plansLoading: true,
      info: false,
      plans: false,
      showCCForm: false,
      submittingCC: false,
      showPlans: false,
      changeSubscription: null
    };
  },
  mounted() {
    this.fetchPlans();
    if (this.$root.me.id == undefined) {
      this.$root.$on("me.ready", this.fetchInfo);
    } else {
      this.fetchInfo();
    }
  },
  methods: {
    fetchInfo() {
      this.$http
        .get(this.$config.API + "/billing/user/" + this.$root.me.id)
        .then(d => {
          this.loading = false;
          this.info = d.data;
          this.changeSubscription = this.info.subscriptions[0].plan;
        })
        .catch(() => {
          this.loading = false;
          this.$message({
            message: "Failed to fetch your billing info. Please try again.",
            type: "error"
          });
        });
    },
    fetchPlans() {
      this.$http.get(this.$config.API + "/billing/plans").then(d => {
        this.plans = d.data.plans;
        this.plansLoading = false;
      });
    },
    getPlan(planId) {
      if (this.plansLoading)
        return {
          name: ""
        };
      var plan = this._.find(this.plans, {
        id: planId
      });
      return plan;
    },
    timeUntil(time) {
      return this.$moment().to(moment.unix(time));
    },
    formatTime(time) {
      return this.$moment.unix(time).format("LL");
    },
    products() {
      if (this.loading == true || this.plansLoading == true) return [];
      var currentPlan = this._.find(this.plans, {
        id: this.info.subscriptions[0].plan
      });
      return this._.sortBy(
        this._.reduce(
          this.plans,
          (carry, plan) => {
            if (
              (currentPlan.interval == "year" && plan.interval == "year") ||
              currentPlan.interval == "month"
            ) {
              (
                carry[plan.product.id] ||
                (carry[plan.product.id] = Object.assign(
                  {
                    plans: []
                  },
                  plan.product
                ))
              ).plans.push(plan);
            }
            return carry;
          },
          {}
        ),
        ["metadata.order"]
      );
    }
  },
  filters: {
    formatAmount(value) {
      var decValue = parseFloat(value) / 100;
      return decValue.toFixed(2).replace(/\d(?=(\d{3})+\.)/g, "$&,");
    }
  },
  computed: {
    isExpiringSoon() {
      if (this.info.sources.length == 0) return false;
      var exprDate = moment(
        this.info.sources[0].card.expMonth +
          "-" +
          this.info.sources[0].card.expYear,
        ["MM-YY", "MM-YYYY"]
      );
      return moment()
        .add(2, "months")
        .isSameOrAfter(exprDate);
    }
  }
};
</script>
