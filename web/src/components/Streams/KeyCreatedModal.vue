<template>
	<modal ref="modal" size="md" @hidden="hide" closable>
		<div slot="title">Your API key</div>
		<div slot="body">
			<div class="rounded p-4 border-2 border-orange-200 bg-orange-100 text-orange-900">
				<div class="flex">
					<div class="flex-initial pl-3"><i class="fas fa-exclamation-circle fa-3x"></i></div>
					<div class="flex-initial pl-5">
						The secret will <b>not</b> be displayed again.<br/>
						Please make sure you keep a <i>secure</i> copy.
					</div>
				</div>
			</div>
			<div class="m-4">
				<b>Key</b> <input type="text" :value="ks" class="input-text" readonly/>
			</div>
			<button @click="downloadJSON" class="mr-4 mb-6 float-right bg-ev-100 text-white hover:bg-ev-700 py-2 px-3 rounded">Download as JSON</button>
		</div>
	</modal>
</template>
<script>
import modal from '@/components/Modal.vue'

export default {
	name: 'key-created-modal',
	components: {
		modal
	},
	props: {
		response: Object,
	},
	computed: {
		ks() {
			if (this.response) {
				return this.response.key+":"+this.response.secret
			}
			return ""
		}
	},
	methods: {
		hide() {
			this.$emit("hidden")
		},
		show() {
			this.$refs.modal.show()
		},
		downloadJSON() {
			var file = new Blob(["{\n\t\"key\":\""+this.ks+"\"\n}"], {type: "application/json"})
			var a = document.createElement("a"), url = URL.createObjectURL(file);
			a.href=url;
			a.download = "key.json"
			document.body.appendChild(a);
			a.click()
			this.$nextTick(() => {
				document.body.removeChild(a);
				URL.revokeObjectURL(url);
			})
		}
	}
}
</script>