<template>
  <div class="panel">
    <div class="d-head rounded-t clearfix h-12">
      <div class="float-right" style="margin-top:2px" v-if="keys.length > 0">
        <button @click="showCreate">Create key</button>
      </div>
    </div>
    <table class="d-table">
      <thead>
        <tr>
          <th>Description</th>
          <th>Permissions</th>
          <th>Restriction</th>
          <th>Last used</th>
          <th>Key</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="key in keys" :key="key.id">
          <td class="w-1/4">{{ key.label }}</td>
          <td class="text-xs">
            {{formatPerms(key)}}
          </td>
          <td class="text-xs" v-if="key.restrictionFilter">{{key.restrictionFilter}}</td>
          <td class="text-xs" v-else><i>None</i></td>
          <td class="text-sm"><i>Never used</i></td>
          <td>{{ key.key }}:•••••••••••••</td>
          <td class="w-1 actions">
            <button @click="editKey(key)">Edit</button>
            <button @click="confirmDelete(key)">Remove</button>
          </td>
        </tr>
        <tr v-if="keys.length == 0">
          <td colspan="999" class="text-center">
            <div class="text-gray-500 text-sm my-6">No keys exist for this stream. Click 'Create key' to create one...</div>
            <button @click="showCreate" class="mb-4 bg-ev-100 rounded text-white py-2 px-4 hover:bg-ev-700">Create key</button>
          </td>
        </tr>
      </tbody>
    </table>
    <modal size="md" closable ref="keyModal">
      <div slot="title">
        <template v-if="modalForm.id == 0">Create a new key</template>
        <template v-else>Edit {{modalForm.label}}</template>
      </div>
      <div slot="body">
        <div class="mx-auto w-4/5 p-2">
          <div class="mb-4">
            <label for="name">Description *</label>
            <input type="text" name="label" class="input-text" v-model="modalForm.label" ref="stream_name" :class="{error: errors.label}"/>
            <div v-if="errors.label" class="error-text">{{errors.label}}</div>
          </div>
          <div class="mb-6">
            <label><b>Permissions</b></label>
            <div class="pl-3 mt-1">
              <div class="mb-1">
                <label class="input-check">
                  <input type="checkbox" name="perms_pub" v-model="modalForm.permissions.publish" />
                  <span>Publish <i class="text-xs text-gray-500">&nbsp; Allows key user to publish events</i></span>
                </label>
              </div>
              <div class="mb-1">
                <label class="input-check">
                  <input type="checkbox" name="perms_pub" v-model="modalForm.permissions.subscribe" />
                  <span>Subscribe <i class="text-xs text-gray-500">&nbsp; Allows key user to subscribe to subjects</i></span>
                </label>
              </div>
              <div class="mb-1">
                <label class="input-check">
                  <input type="checkbox" name="perms_pub" v-model="modalForm.permissions.replay" />
                  <span>Replay  <i class="text-xs text-gray-500">&nbsp; Allows key user to replay events to all/other subscribers</i></span>
                </label>
              </div>
            </div>
          </div>
          <div class="mb-4">
            <label for="name"><b>Restriction</b> &nbsp;<i class="text-xs text-gray-500">RegExp filter applied to the subject of the event</i></label>
            <input type="text" placeholder="example: /user\.[created|deleted]/" name="restrictionFilter" class="input-text" v-model="modalForm.restrictionFilter" ref="stream_name" :class="{error: errors.restrictionFilter}"/>
            <div v-if="errors.restrictionFilter" class="error-text">{{errors.restrictionFilter}}</div>
          </div>
          <div class="text-right mb-3">
            <button @click="submitModal" class="bg-ev-100 text-white hover:bg-ev-700 py-2 px-3 rounded">
              <template v-if="modalForm.id == 0">Create</template>
              <template v-else>Save</template>
            </button>
          </div>
        </div>
      </div>
    </modal>
    <key-created-modal :response="keyCreateResponse" @hidden="clearKeyCreateResponse" ref="keyCreatedModal"></key-created-modal>
    <modal size="sm" ref="deleteModal">
			<div slot="title">Confirm</div>
			<div slot="body">
				<div>Are you sure you want to delete <i>"{{findById(keyToDelete).label}}"</i>?<br/><b>This cannot be undone!</b></div>
				<div class="text-right mt-3">
					<button @click="deleteKey" class="bg-gray-100 hover:bg-red-700 hover:text-white py-2 px-3 rounded mr-3">Yes, Delete!</button>
					<button @click="hideDelete" class="bg-ev-100 text-white hover:bg-green-700 py-2 px-3 rounded">No</button>
				</div>
			</div>
		</modal>
  </div>
</template>
<script>
import modal from '@/components/Modal.vue';
import keyCreatedModal from '@/components/Streams/KeyCreatedModal.vue';

export default {
  name: "stream-auth",
  components: {
    modal,
    keyCreatedModal,
  },
  data() {
    return {
      loading: true,
      keys: [],
      errors: {
        label: "",
      },
      keyToDelete: 0,
      keyCreateResponse: {},
      modalForm: {
        id: 0,
        label: "",
        restrictionFilter: "",
        permissions: {
          publish: true,
          subscribe: true,
          replay: true,
        }
      }
    };
  },
  created() {
    this.load();
  },
  methods: {
    load() {
      this.$http.get(this.$config.API + "/stream/" + this.$route.params.id + "/keys")
      .then(d => {
        this.loading = false;
        this.keys = d.data.keys || [];
      });
    },
    findById(id) {
      var f = {};
      this.keys.forEach(v => {
        if (v.id == id) {
          f = v;
          return false;
        }
      })

      return f
    },
    clearKeyCreateResponse() {
      this.keyCreateResponse = {}
    },
    formatPerms(key) {
      var perms = [];

      if(!key.permissions) {
        perms.push("Pub");
        perms.push("Sub");
        perms.push("Replay");
      } else {
        if (key.permissions.publish) {
          perms.push("Pub");
        }
        if (key.permissions.subscribe) {
          perms.push("Sub");
        }
        if (key.permissions.replay) {
          perms.push("Replay");
        }
      }

      return perms.join(", ")
    },
    clearModal() {
      this.modalForm.id = 0;
      this.modalForm.label = "";
      this.modalForm.restrictionFilter = "";
      this.modalForm.permissions = {
        publish: true,
        subscribe: true,
        replay: true,
      };
    },
    showCreate() {
      this.clearModal();
      this.$refs.keyModal.show();
    },
    editKey(key) {
      Object.assign(this.modalForm, key);
      this.$nextTick(() => {
        this.$refs.keyModal.show();
      });
    },
    submitModal() {
      if (this.modalForm.id == 0) {
        this.submitCreate();
      } else {
        this.submitEdit();
      }
    },
    submitCreate() {
      this.$http.post(this.$config.API + "/stream/" + this.$parent.stream.id + "/key", {
        label: this.modalForm.label,
        permissions: this.modalForm.permissions,
        restrictionFilter: this.modalForm.restrictionFilter,
      }).then(d => {
        this.loading = true;
        this.load();
        this.keyCreateResponse = d.data;
        this.$refs.keyModal.close();
        this.$refs.keyCreatedModal.show();
      }).catch(err => {
        this.$message.error("Failed to create your key. Please try again")
      })
    },
    submitEdit() {
      this.$http.patch(this.$config.API + "/stream/"+this.$parent.stream.id+"/key/" + this.modalForm.id, {
        label: this.modalForm.label,
        permissions: this.modalForm.permissions,
        restrictionFilter: this.modalForm.restrictionFilter,
      }).then(d => {
        this.loading = true;
        this.load();
        this.$refs.keyModal.close();
        this.clearModal();
      }).catch(err => {
        this.$message.error("Failed to update your key. Please try again")
      })
    },
    confirmDelete(key) {
      this.keyToDelete = key.id;
      this.$refs.deleteModal.show();
    },
    deleteKey() {
      this.$http.delete(this.$config.API + "/stream/"+this.$parent.stream.id+"/key/"+this.keyToDelete).then(d => {
        this.hideDelete();
        this.load();
        this.$message.success("Successfully deleted key");
      }).catch(err => {
        this.$message.error("Failed to delete your key. Please try again")
      })
    },
    hideDelete() {
      this.keyToDelete = 0;
      this.$refs.deleteModal.close();
    }
  }
};
</script>