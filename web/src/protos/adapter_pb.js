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
goog.exportSymbol('proto.evntsrc.adapter.Adapter', null, global);
goog.exportSymbol('proto.evntsrc.adapter.Adapter.EngineType', null, global);

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
proto.evntsrc.adapter.Adapter = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.evntsrc.adapter.Adapter.oneofGroups_);
};
goog.inherits(proto.evntsrc.adapter.Adapter, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.evntsrc.adapter.Adapter.displayName = 'proto.evntsrc.adapter.Adapter';
}
/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.evntsrc.adapter.Adapter.oneofGroups_ = [[5,6,7]];

/**
 * @enum {number}
 */
proto.evntsrc.adapter.Adapter.ContextCase = {
  CONTEXT_NOT_SET: 0,
  STREAM: 5,
  SUBJECT: 6,
  EVENTTYPE: 7
};

/**
 * @return {proto.evntsrc.adapter.Adapter.ContextCase}
 */
proto.evntsrc.adapter.Adapter.prototype.getContextCase = function() {
  return /** @type {proto.evntsrc.adapter.Adapter.ContextCase} */(jspb.Message.computeOneofCase(this, proto.evntsrc.adapter.Adapter.oneofGroups_[0]));
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
proto.evntsrc.adapter.Adapter.prototype.toObject = function(opt_includeInstance) {
  return proto.evntsrc.adapter.Adapter.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.evntsrc.adapter.Adapter} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.evntsrc.adapter.Adapter.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    label: jspb.Message.getFieldWithDefault(msg, 2, ""),
    engine: jspb.Message.getFieldWithDefault(msg, 3, 0),
    code: msg.getCode_asB64(),
    stream: jspb.Message.getFieldWithDefault(msg, 5, 0),
    subject: jspb.Message.getFieldWithDefault(msg, 6, ""),
    eventtype: jspb.Message.getFieldWithDefault(msg, 7, ""),
    metadataMap: (f = msg.getMetadataMap()) ? f.toObject(includeInstance, undefined) : []
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
 * @return {!proto.evntsrc.adapter.Adapter}
 */
proto.evntsrc.adapter.Adapter.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.evntsrc.adapter.Adapter;
  return proto.evntsrc.adapter.Adapter.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.evntsrc.adapter.Adapter} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.evntsrc.adapter.Adapter}
 */
proto.evntsrc.adapter.Adapter.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setLabel(value);
      break;
    case 3:
      var value = /** @type {!proto.evntsrc.adapter.Adapter.EngineType} */ (reader.readEnum());
      msg.setEngine(value);
      break;
    case 4:
      var value = /** @type {!Uint8Array} */ (reader.readBytes());
      msg.setCode(value);
      break;
    case 5:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setStream(value);
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.setSubject(value);
      break;
    case 7:
      var value = /** @type {string} */ (reader.readString());
      msg.setEventtype(value);
      break;
    case 8:
      var value = msg.getMetadataMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readString, null, "");
         });
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
proto.evntsrc.adapter.Adapter.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.evntsrc.adapter.Adapter.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.evntsrc.adapter.Adapter} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.evntsrc.adapter.Adapter.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getLabel();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getEngine();
  if (f !== 0.0) {
    writer.writeEnum(
      3,
      f
    );
  }
  f = message.getCode_asU8();
  if (f.length > 0) {
    writer.writeBytes(
      4,
      f
    );
  }
  f = /** @type {number} */ (jspb.Message.getField(message, 5));
  if (f != null) {
    writer.writeInt32(
      5,
      f
    );
  }
  f = /** @type {string} */ (jspb.Message.getField(message, 6));
  if (f != null) {
    writer.writeString(
      6,
      f
    );
  }
  f = /** @type {string} */ (jspb.Message.getField(message, 7));
  if (f != null) {
    writer.writeString(
      7,
      f
    );
  }
  f = message.getMetadataMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(8, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeString);
  }
};


/**
 * @enum {number}
 */
proto.evntsrc.adapter.Adapter.EngineType = {
  V8: 0,
  PYTHON: 1,
  LUA: 2
};

/**
 * optional string ID = 1;
 * @return {string}
 */
proto.evntsrc.adapter.Adapter.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.evntsrc.adapter.Adapter.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string label = 2;
 * @return {string}
 */
proto.evntsrc.adapter.Adapter.prototype.getLabel = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.evntsrc.adapter.Adapter.prototype.setLabel = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional EngineType engine = 3;
 * @return {!proto.evntsrc.adapter.Adapter.EngineType}
 */
proto.evntsrc.adapter.Adapter.prototype.getEngine = function() {
  return /** @type {!proto.evntsrc.adapter.Adapter.EngineType} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/** @param {!proto.evntsrc.adapter.Adapter.EngineType} value */
proto.evntsrc.adapter.Adapter.prototype.setEngine = function(value) {
  jspb.Message.setProto3EnumField(this, 3, value);
};


/**
 * optional bytes Code = 4;
 * @return {!(string|Uint8Array)}
 */
proto.evntsrc.adapter.Adapter.prototype.getCode = function() {
  return /** @type {!(string|Uint8Array)} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * optional bytes Code = 4;
 * This is a type-conversion wrapper around `getCode()`
 * @return {string}
 */
proto.evntsrc.adapter.Adapter.prototype.getCode_asB64 = function() {
  return /** @type {string} */ (jspb.Message.bytesAsB64(
      this.getCode()));
};


/**
 * optional bytes Code = 4;
 * Note that Uint8Array is not supported on all browsers.
 * @see http://caniuse.com/Uint8Array
 * This is a type-conversion wrapper around `getCode()`
 * @return {!Uint8Array}
 */
proto.evntsrc.adapter.Adapter.prototype.getCode_asU8 = function() {
  return /** @type {!Uint8Array} */ (jspb.Message.bytesAsU8(
      this.getCode()));
};


/** @param {!(string|Uint8Array)} value */
proto.evntsrc.adapter.Adapter.prototype.setCode = function(value) {
  jspb.Message.setProto3BytesField(this, 4, value);
};


/**
 * optional int32 stream = 5;
 * @return {number}
 */
proto.evntsrc.adapter.Adapter.prototype.getStream = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 5, 0));
};


/** @param {number} value */
proto.evntsrc.adapter.Adapter.prototype.setStream = function(value) {
  jspb.Message.setOneofField(this, 5, proto.evntsrc.adapter.Adapter.oneofGroups_[0], value);
};


proto.evntsrc.adapter.Adapter.prototype.clearStream = function() {
  jspb.Message.setOneofField(this, 5, proto.evntsrc.adapter.Adapter.oneofGroups_[0], undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.evntsrc.adapter.Adapter.prototype.hasStream = function() {
  return jspb.Message.getField(this, 5) != null;
};


/**
 * optional string subject = 6;
 * @return {string}
 */
proto.evntsrc.adapter.Adapter.prototype.getSubject = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 6, ""));
};


/** @param {string} value */
proto.evntsrc.adapter.Adapter.prototype.setSubject = function(value) {
  jspb.Message.setOneofField(this, 6, proto.evntsrc.adapter.Adapter.oneofGroups_[0], value);
};


proto.evntsrc.adapter.Adapter.prototype.clearSubject = function() {
  jspb.Message.setOneofField(this, 6, proto.evntsrc.adapter.Adapter.oneofGroups_[0], undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.evntsrc.adapter.Adapter.prototype.hasSubject = function() {
  return jspb.Message.getField(this, 6) != null;
};


/**
 * optional string eventType = 7;
 * @return {string}
 */
proto.evntsrc.adapter.Adapter.prototype.getEventtype = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 7, ""));
};


/** @param {string} value */
proto.evntsrc.adapter.Adapter.prototype.setEventtype = function(value) {
  jspb.Message.setOneofField(this, 7, proto.evntsrc.adapter.Adapter.oneofGroups_[0], value);
};


proto.evntsrc.adapter.Adapter.prototype.clearEventtype = function() {
  jspb.Message.setOneofField(this, 7, proto.evntsrc.adapter.Adapter.oneofGroups_[0], undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.evntsrc.adapter.Adapter.prototype.hasEventtype = function() {
  return jspb.Message.getField(this, 7) != null;
};


/**
 * map<string, string> metadata = 8;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,string>}
 */
proto.evntsrc.adapter.Adapter.prototype.getMetadataMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,string>} */ (
      jspb.Message.getMapField(this, 8, opt_noLazyCreate,
      null));
};


proto.evntsrc.adapter.Adapter.prototype.clearMetadataMap = function() {
  this.getMetadataMap().clear();
};


goog.object.extend(exports, proto.evntsrc.adapter);
