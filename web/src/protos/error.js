var jspb = require('google-protobuf');

export default {
	readError(data) {

		data = btoa(
			new Uint8Array(data)
			.reduce((data, byte) => data + String.fromCharCode(byte), '')
		)

		var reader = new jspb.BinaryReader(data);

		var code = 0;
		var message = "";

		while (reader.nextField()) {
			if (reader.isEndGroup()) {
				break;
			}
			var field = reader.getFieldNumber();
			switch (field) {
				case 1:
					message = reader.readString();
					break;
				case 2:
					code = reader.readInt32();
					break;
				default:
					reader.skipField();
					break;
			}
		}

		return {
			code,
			message
		};
	}
}