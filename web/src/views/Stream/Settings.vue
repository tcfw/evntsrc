<template>
	<div class="panel">
		<button @click="showDelete" class="bg-red-700 hover:bg-red-600 text-white py-2 px-5 text-sm font-bold rounded m-3">Delete Stream</button>
		<modal size="sm" ref="deleteModal">
			<div slot="title">Confirm</div>
			<div slot="body">
				<div>Are you sure you want to delete <i>"{{$parent.stream.name}}"</i>?<br/><b>This cannot be undone!</b></div>
				<div class="text-right mt-3">
					<button @click="deleteStream" class="bg-gray-100 hover:bg-red-700 hover:text-white py-2 px-3 rounded mr-3">Yes, Delete!</button>
					<button @click="hideDelete" class="bg-ev-100 text-white hover:bg-green-700 py-2 px-3 rounded">No</button>
				</div>
			</div>
		</modal>
	</div>
</template>
<script>
import modal from '@/components/Modal.vue'
export default {
	name:'stream-settings',
	components: {
		modal
	},
	methods: {
		showDelete() {
			this.$refs.deleteModal.show()
		},
		hideDelete() {
			this.$refs.deleteModal.close()
		},
		deleteStream() {
			this.$http.delete(this.$config.API + '/stream/'+this.$parent.stream.id).then(d => {
				this.hideDelete();
				this.$message.success("Deleted "+this.$parent.stream.name+" successfully")
				this.$nextTick(() => {
					this.$parent.$parent.reload();
					this.$router.push('/streams')
				})
			});
		}
	}
}
</script>
