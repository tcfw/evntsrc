/**
 * @fileoverview
 * @enhanceable
 * @suppress {messageConventions} JS Compiler reports an error if a variable or
 *     field starts with 'MSG_' and isn't a translatable message.
 * @public
 */
// GENERATED CODE -- DO NOT EDIT!

var jspb = require('google-protobuf');
var goog = jspb;
var global = Function('return this')();

var github_com_gogo_protobuf_gogoproto_gogo_pb = require('./github.com/gogo/protobuf/gogoproto/gogo_pb.js');
goog.object.extend(proto, github_com_gogo_protobuf_gogoproto_gogo_pb);
var google_api_annotations_pb = require('./google/api/annotations_pb.js');
goog.object.extend(proto, google_api_annotations_pb);
goog.exportSymbol('proto.evntsrc.stsmetrics.MetricCount', null, global);
goog.exportSymbol('proto.evntsrc.stsmetrics.MetricCount.Timestamp', null, global);
goog.exportSymbol('proto.evntsrc.stsmetrics.MetricsCountRequest', null, global);
goog.exportSymbol('proto.evntsrc.stsmetrics.MetricsCountRequest.Interval', null, global);
goog.exportSymbol('proto.evntsrc.stsmetrics.MetricsCountResponse', null, global);
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.evntsrc.stsmetrics.MetricsCountRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.evntsrc.stsmetrics.MetricsCountRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.evntsrc.stsmetrics.MetricsCountRequest.displayName = 'proto.evntsrc.stsmetrics.MetricsCountRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.evntsrc.stsmetrics.MetricCount = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.evntsrc.stsmetrics.MetricCount, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.evntsrc.stsmetrics.MetricCount.displayName = 'proto.evntsrc.stsmetrics.MetricCount';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.evntsrc.stsmetrics.MetricCount.Timestamp = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.evntsrc.stsmetrics.MetricCount.Timestamp, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.evntsrc.stsmetrics.MetricCount.Timestamp.displayName = 'proto.evntsrc.stsmetrics.MetricCount.Timestamp';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.evntsrc.stsmetrics.MetricsCountResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.evntsrc.stsmetrics.MetricsCountResponse.repeatedFields_, null);
};
goog.inherits(proto.evntsrc.stsmetrics.MetricsCountResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.evntsrc.stsmetrics.MetricsCountResponse.displayName = 'proto.evntsrc.stsmetrics.MetricsCountResponse';
}



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.evntsrc.stsmetrics.MetricsCountRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.evntsrc.stsmetrics.MetricsCountRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.evntsrc.stsmetrics.MetricsCountRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.evntsrc.stsmetrics.MetricsCountRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    stream: jspb.Message.getFieldWithDefault(msg, 1, 0),
    interval: jspb.Message.getFieldWithDefault(msg, 2, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.evntsrc.stsmetrics.MetricsCountRequest}
 */
proto.evntsrc.stsmetrics.MetricsCountRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.evntsrc.stsmetrics.MetricsCountRequest;
  return proto.evntsrc.stsmetrics.MetricsCountRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.evntsrc.stsmetrics.MetricsCountRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.evntsrc.stsmetrics.MetricsCountRequest}
 */
proto.evntsrc.stsmetrics.MetricsCountRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setStream(value);
      break;
    case 2:
      var value = /** @type {!proto.evntsrc.stsmetrics.MetricsCountRequest.Interval} */ (reader.readEnum());
      msg.setInterval(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.evntsrc.stsmetrics.MetricsCountRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.evntsrc.stsmetrics.MetricsCountRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.evntsrc.stsmetrics.MetricsCountRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.evntsrc.stsmetrics.MetricsCountRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getStream();
  if (f !== 0) {
    writer.writeInt32(
      1,
      f
    );
  }
  f = message.getInterval();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
};


/**
 * @enum {number}
 */
proto.evntsrc.stsmetrics.MetricsCountRequest.Interval = {
  MIN10: 0,
  MIN30: 1,
  HOUR: 2,
  HOUR12: 3,
  DAY: 4,
  WEEK: 5,
  MONTH: 6
};

/**
 * optional int32 stream = 1;
 * @return {number}
 */
proto.evntsrc.stsmetrics.MetricsCountRequest.prototype.getStream = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/** @param {number} value */
proto.evntsrc.stsmetrics.MetricsCountRequest.prototype.setStream = function(value) {
  jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional Interval interval = 2;
 * @return {!proto.evntsrc.stsmetrics.MetricsCountRequest.Interval}
 */
proto.evntsrc.stsmetrics.MetricsCountRequest.prototype.getInterval = function() {
  return /** @type {!proto.evntsrc.stsmetrics.MetricsCountRequest.Interval} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/** @param {!proto.evntsrc.stsmetrics.MetricsCountRequest.Interval} value */
proto.evntsrc.stsmetrics.MetricsCountRequest.prototype.setInterval = function(value) {
  jspb.Message.setProto3EnumField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.evntsrc.stsmetrics.MetricCount.prototype.toObject = function(opt_includeInstance) {
  return proto.evntsrc.stsmetrics.MetricCount.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.evntsrc.stsmetrics.MetricCount} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.evntsrc.stsmetrics.MetricCount.toObject = function(includeInstance, msg) {
  var f, obj = {
    count: +jspb.Message.getFieldWithDefault(msg, 1, 0.0),
    timestamp: (f = msg.getTimestamp()) && proto.evntsrc.stsmetrics.MetricCount.Timestamp.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.evntsrc.stsmetrics.MetricCount}
 */
proto.evntsrc.stsmetrics.MetricCount.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.evntsrc.stsmetrics.MetricCount;
  return proto.evntsrc.stsmetrics.MetricCount.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.evntsrc.stsmetrics.MetricCount} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.evntsrc.stsmetrics.MetricCount}
 */
proto.evntsrc.stsmetrics.MetricCount.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readFloat());
      msg.setCount(value);
      break;
    case 2:
      var value = new proto.evntsrc.stsmetrics.MetricCount.Timestamp;
      reader.readMessage(value,proto.evntsrc.stsmetrics.MetricCount.Timestamp.deserializeBinaryFromReader);
      msg.setTimestamp(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.evntsrc.stsmetrics.MetricCount.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.evntsrc.stsmetrics.MetricCount.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.evntsrc.stsmetrics.MetricCount} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.evntsrc.stsmetrics.MetricCount.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getCount();
  if (f !== 0.0) {
    writer.writeFloat(
      1,
      f
    );
  }
  f = message.getTimestamp();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.evntsrc.stsmetrics.MetricCount.Timestamp.serializeBinaryToWriter
    );
  }
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.evntsrc.stsmetrics.MetricCount.Timestamp.prototype.toObject = function(opt_includeInstance) {
  return proto.evntsrc.stsmetrics.MetricCount.Timestamp.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.evntsrc.stsmetrics.MetricCount.Timestamp} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.evntsrc.stsmetrics.MetricCount.Timestamp.toObject = function(includeInstance, msg) {
  var f, obj = {
    seconds: jspb.Message.getFieldWithDefault(msg, 1, 0),
    nanos: jspb.Message.getFieldWithDefault(msg, 2, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.evntsrc.stsmetrics.MetricCount.Timestamp}
 */
proto.evntsrc.stsmetrics.MetricCount.Timestamp.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.evntsrc.stsmetrics.MetricCount.Timestamp;
  return proto.evntsrc.stsmetrics.MetricCount.Timestamp.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.evntsrc.stsmetrics.MetricCount.Timestamp} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.evntsrc.stsmetrics.MetricCount.Timestamp}
 */
proto.evntsrc.stsmetrics.MetricCount.Timestamp.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setSeconds(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setNanos(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.evntsrc.stsmetrics.MetricCount.Timestamp.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.evntsrc.stsmetrics.MetricCount.Timestamp.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.evntsrc.stsmetrics.MetricCount.Timestamp} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.evntsrc.stsmetrics.MetricCount.Timestamp.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSeconds();
  if (f !== 0) {
    writer.writeInt64(
      1,
      f
    );
  }
  f = message.getNanos();
  if (f !== 0) {
    writer.writeInt32(
      2,
      f
    );
  }
};


/**
 * optional int64 seconds = 1;
 * @return {number}
 */
proto.evntsrc.stsmetrics.MetricCount.Timestamp.prototype.getSeconds = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/** @param {number} value */
proto.evntsrc.stsmetrics.MetricCount.Timestamp.prototype.setSeconds = function(value) {
  jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional int32 nanos = 2;
 * @return {number}
 */
proto.evntsrc.stsmetrics.MetricCount.Timestamp.prototype.getNanos = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/** @param {number} value */
proto.evntsrc.stsmetrics.MetricCount.Timestamp.prototype.setNanos = function(value) {
  jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional float count = 1;
 * @return {number}
 */
proto.evntsrc.stsmetrics.MetricCount.prototype.getCount = function() {
  return /** @type {number} */ (+jspb.Message.getFieldWithDefault(this, 1, 0.0));
};


/** @param {number} value */
proto.evntsrc.stsmetrics.MetricCount.prototype.setCount = function(value) {
  jspb.Message.setProto3FloatField(this, 1, value);
};


/**
 * optional Timestamp timestamp = 2;
 * @return {?proto.evntsrc.stsmetrics.MetricCount.Timestamp}
 */
proto.evntsrc.stsmetrics.MetricCount.prototype.getTimestamp = function() {
  return /** @type{?proto.evntsrc.stsmetrics.MetricCount.Timestamp} */ (
    jspb.Message.getWrapperField(this, proto.evntsrc.stsmetrics.MetricCount.Timestamp, 2));
};


/** @param {?proto.evntsrc.stsmetrics.MetricCount.Timestamp|undefined} value */
proto.evntsrc.stsmetrics.MetricCount.prototype.setTimestamp = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.evntsrc.stsmetrics.MetricCount.prototype.clearTimestamp = function() {
  this.setTimestamp(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.evntsrc.stsmetrics.MetricCount.prototype.hasTimestamp = function() {
  return jspb.Message.getField(this, 2) != null;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.evntsrc.stsmetrics.MetricsCountResponse.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.evntsrc.stsmetrics.MetricsCountResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.evntsrc.stsmetrics.MetricsCountResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.evntsrc.stsmetrics.MetricsCountResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.evntsrc.stsmetrics.MetricsCountResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    metricsList: jspb.Message.toObjectList(msg.getMetricsList(),
    proto.evntsrc.stsmetrics.MetricCount.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.evntsrc.stsmetrics.MetricsCountResponse}
 */
proto.evntsrc.stsmetrics.MetricsCountResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.evntsrc.stsmetrics.MetricsCountResponse;
  return proto.evntsrc.stsmetrics.MetricsCountResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.evntsrc.stsmetrics.MetricsCountResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.evntsrc.stsmetrics.MetricsCountResponse}
 */
proto.evntsrc.stsmetrics.MetricsCountResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.evntsrc.stsmetrics.MetricCount;
      reader.readMessage(value,proto.evntsrc.stsmetrics.MetricCount.deserializeBinaryFromReader);
      msg.addMetrics(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.evntsrc.stsmetrics.MetricsCountResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.evntsrc.stsmetrics.MetricsCountResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.evntsrc.stsmetrics.MetricsCountResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.evntsrc.stsmetrics.MetricsCountResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getMetricsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto.evntsrc.stsmetrics.MetricCount.serializeBinaryToWriter
    );
  }
};


/**
 * repeated MetricCount metrics = 1;
 * @return {!Array<!proto.evntsrc.stsmetrics.MetricCount>}
 */
proto.evntsrc.stsmetrics.MetricsCountResponse.prototype.getMetricsList = function() {
  return /** @type{!Array<!proto.evntsrc.stsmetrics.MetricCount>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.evntsrc.stsmetrics.MetricCount, 1));
};


/** @param {!Array<!proto.evntsrc.stsmetrics.MetricCount>} value */
proto.evntsrc.stsmetrics.MetricsCountResponse.prototype.setMetricsList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.evntsrc.stsmetrics.MetricCount=} opt_value
 * @param {number=} opt_index
 * @return {!proto.evntsrc.stsmetrics.MetricCount}
 */
proto.evntsrc.stsmetrics.MetricsCountResponse.prototype.addMetrics = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.evntsrc.stsmetrics.MetricCount, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.evntsrc.stsmetrics.MetricsCountResponse.prototype.clearMetricsList = function() {
  this.setMetricsList([]);
};


goog.object.extend(exports, proto.evntsrc.stsmetrics);