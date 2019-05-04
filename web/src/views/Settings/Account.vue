<template>
  <div style="padding: 15px;">
    <el-row>
      <el-col :span="16">
        <el-breadcrumb separator-class="el-icon-arrow-right">
          <el-breadcrumb-item :to="{ path: '/' }">Home</el-breadcrumb-item>
          <el-breadcrumb-item :to="{ path: '/settings' }"
            >Settings</el-breadcrumb-item
          >
          <el-breadcrumb-item>Account</el-breadcrumb-item>
        </el-breadcrumb>
        <h3>Profile</h3>
        <el-card class="box-card">
          <el-form
            label-position="left"
            :rules="editProfileRules"
            ref="editProfileForm"
            label-width="200px"
            :model="editProfileForm"
          >
            <el-form-item label="Name" prop="name">
              <el-input v-model="editProfileForm.name"></el-input>
            </el-form-item>
            <el-form-item label="Email" prop="email">
              <el-input v-model="editProfileForm.email"></el-input>
            </el-form-item>
            <el-form-item label="Company" prop="company">
              <el-input v-model="editProfileForm.company"></el-input>
            </el-form-item>
            <br />
            <el-form-item label="Current password" prop="current_password">
              <el-input
                autocomplete="off"
                type="password"
                v-model="editProfileForm.current_password"
                placeholder="Leave blank to keep the current password"
              ></el-input>
            </el-form-item>
            <el-form-item label="New password" prop="new_password">
              <el-input
                autocomplete="off"
                type="password"
                v-model="editProfileForm.new_password"
              ></el-input>
            </el-form-item>
            <el-form-item label="Confirm new password" prop="c_new_password">
              <el-input
                autocomplete="off"
                type="password"
                v-model="editProfileForm.c_new_password"
              ></el-input>
            </el-form-item>
            <el-button
              type="success"
              size="small"
              :loading="submittingEditProfileForm"
              @click="saveProfile"
              >Update profile</el-button
            >
          </el-form>
        </el-card>
        <hr />
        <h5>Delete account</h5>
        <small
          >If you delete your account, there is no going back. This will be
          perminant.<br />This will also include all stored events.</small
        >
        <el-button
          type="danger"
          style="margin-top: 15px;"
          size="small"
          icon="fas fa-trash"
        >
          &nbsp;Delete account</el-button
        >
      </el-col>
      <el-col :span="8"> </el-col>
    </el-row>
  </div>
</template>
<script>
import users_pb from "@/protos/users_pb.js";

export default {
  data() {
    return {
      submittingEditProfileForm: false,
      editProfileForm: {
        name: "",
        email: "",
        company: "",
        current_password: "",
        new_password: "",
        c_new_password: ""
      }
    };
  },
  computed: {
    editProfileRules() {
      return {
        name: [{ required: true, trigger: "blur" }],
        email: [
          { required: true, trigger: "blur" },
          { type: "email", trigger: "blur" }
        ],
        new_password: [
          {
            required: this.editProfileForm.current_password != "",
            trigger: "blur"
          }
        ],
        c_new_password: [
          {
            required: this.editProfileForm.current_password != "",
            trigger: "blur"
          }
        ]
      };
    }
  },
  mounted() {
    axios.get(this.$config.API + "/me").then(d => {
      this.editProfileForm = d.data;
      this.editProfileForm.current_password = "";
      this.editProfileForm.new_password = "";
      this.editProfileForm.c_new_password = "";
    });
  },
  methods: {
    saveProfile() {
      this.$refs["editProfileForm"].validate(valid => {
        if (valid) {
          this.submittingEditProfileForm = true;
          var user = new proto.evntsrc.users.User();
          user.setEmail(this.editProfileForm.email);
          user.setName(this.editProfileForm.name);
          user.setCompany(this.editProfileForm.company);

          var updateRequest = new proto.evntsrc.users.UserUpdateRequest();
          updateRequest.setUser(user);
          axios
            .post(this.$config.API + "/me", updateRequest.toObject())
            .then(resp => {
              this.$root.me = resp.data;
              this.$message({
                message: "Profile updated successfully",
                type: "success"
              });

              if (
                this.editProfileForm.current_password != "" &&
                this.editProfileForm.new_password != "" &&
                this.editProfileForm.c_new_password != "" &&
                this.editProfileForm.new_password ==
                  this.editProfileForm.c_new_password
              ) {
                var passRequest = new proto.evntsrc.users.PasswordUpdateRequest();
                passRequest.setCurrentPassword(
                  this.editProfileForm.current_password
                );
                passRequest.setPassword(this.editProfileForm.new_password);
                axios
                  .post(
                    this.$config.API + "/me/password",
                    passRequest.toObject()
                  )
                  .then(d => {
                    this.submittingEditProfileForm = false;
                    this.editProfileForm.current_password = "";
                    this.editProfileForm.new_password = "";
                    this.editProfileForm.c_new_password = "";
                    this.$message({
                      message: "Profile and password updated successfully",
                      type: "success"
                    });
                  })
                  .catch(err => {
                    this.submittingEditProfileForm = false;
                    const h = this.$createElement;
                    this.$message({
                      message: h("div", { class: "el-message__content" }, [
                        h(
                          "div",
                          null,
                          "Failed to update profile. Please try again."
                        ),
                        h("small", null, err.response.data.error)
                      ]),
                      type: "error"
                    });
                  });
              } else {
                this.submittingEditProfileForm = false;
              }
            })
            .catch(err => {
              this.submittingEditProfileForm = false;
              var errMsg = "";
              if (err != undefined) {
                errMsg += err.response.data.error;
              }
              const h = this.$createElement;
              this.$message({
                message: h("div", { class: "el-message__content" }, [
                  h("div", null, "Failed to update profile. Please try again."),
                  h("small", null, errMsg)
                ]),
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
small {
  font-size: 11px;
  color: #aaa;
}

h5 ~ small {
  margin-top: -20px;
  display: block;
}
</style>
