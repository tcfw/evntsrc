<template>
	<modal ref="modal" size="md" closable position="top" v-on:shown="focus">
		<div slot="title">Create a Stream...</div>
		<div slot="body">
			<div class="mx-auto w-4/5 p-2">
				<div class="mb-2">
					<label for="name">Stream name</label>
					<input type="text" name="name" class="input-text" v-model="form.name" ref="stream_name" :class="{error: errors.name}"/>
					<div v-if="errors.name" class="error-text">{{errors.name}}</div>
				</div>
				<div class="mb-2">
					<label for="name">Data location</label>
					<input type="text" name="location" class="input-text" v-model="form.location" :class="{error: errors.location}"/>
				<div v-if="errors.location" class="error-text">{{errors.location}}</div>
				</div>
				<div class="my-4 text">
					<label class="input-check">
						<input type="checkbox" name="create_key" v-model="form.create_key" />
						<span>Create an API key</span>
					</label>
				</div>
				<button class="input-button-huge my-6" @click="submit">Create Stream</button>
			</div>
		</div>
	</modal>
</template>
<script>
import modal from '@/components/Modal.vue'

export default {
	name: 'create-modal',
	components: {
		modal
	},
	data() {
		return {
			form: {
				name: "",
				location: "",
				create_key: false,
			},
			errors: {
				name: null,
				location: null,
			}
		}
	},
	methods: {
		show() {
			this.$refs.modal.show()
		},
		close() {
			this.$refs.modal.close()
		},
		focus() {
			this.$nextTick(() => {
				this.$refs.stream_name.focus()
			})
		},
		clearForm() {
			this.form.name = this.form.location = "";
		},
		clearErrors() {
			this.errors.name = this.errors.location = null;
		},
		hasErrors() {
			return this.errors.name || this.errors.location;
		},
		validate() {
			let required = "This field is required"

			this.clearErrors();

			if (!this.form.name) {
				this.errors.name = required
			}
			if (!this.form.location) {
				this.errors.location = required
			}

			return this.hasErrors()
		},
		submit() {
			if (this.validate()) {
				return
			}

			this.$http.post(
				this.$config.API + "/stream",
				{name: this.form.name, cluster: this.form.location}
			).then(d => {
				this.$parent.load();
				if (this.form.create_key) {
					this.createKey(d.data.id);
				}
				this.clearForm();
				this.close();
				this.$router.push("/streams/"+d.data.id);
			})
		},
		createKey(id) {
			this.$http.post(
				this.$config.API + "/stream/"+id+"/key",
				{label: "root", permissions: {
					publish: true, subscribe: true, replay: true
				}}
			).then(d => {
				this.$parent.keyCreateResponse = d.data
				this.$parent.$refs.keyCreatedModal.show()
			})
		}
	}
}
</script>