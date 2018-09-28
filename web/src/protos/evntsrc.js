/**
 * @fileoverview
 * @enhanceable
 * @suppress {messageConventions} JS Compiler reports an error if a variable or
 *     field starts with 'MSG_' and isn't a translatable message.
 * @public
 */
// GENERATED CODE -- DO NOT EDIT!


goog.provide('proto.evntsrc.event.Event');

goog.require('jspb.BinaryReader');
goog.require('jspb.BinaryWriter');
goog.require('jspb.Map');
goog.require('jspb.Message');
goog.require('proto.google.protobuf.Timestamp');


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
proto.evntsrc.event.Event = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.evntsrc.event.Event, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.evntsrc.event.Event.displayName = 'proto.evntsrc.event.Event';
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
proto.evntsrc.event.Event.prototype.toObject = function(opt_includeInstance) {
  return proto.evntsrc.event.Event.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.evntsrc.event.Event} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.evntsrc.event.Event.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    stream: jspb.Message.getFieldWithDefault(msg, 2, 0),
    time: (f = msg.getTime()) && proto.google.protobuf.Timestamp.toObject(includeInstance, f),
    type: jspb.Message.getFieldWithDefault(msg, 4, ""),
    typeversion: jspb.Message.getFieldWithDefault(msg, 5, ""),
    ceversion: jspb.Message.getFieldWithDefault(msg, 6, ""),
    source: jspb.Message.getFieldWithDefault(msg, 7, ""),
    subject: jspb.Message.getFieldWithDefault(msg, 8, ""),
    acknowledged: (f = msg.getAcknowledged()) && proto.google.protobuf.Timestamp.toObject(includeInstance, f),
    metadataMap: (f = msg.getMetadataMap()) ? f.toObject(includeInstance, undefined) : [],
    contenttype: jspb.Message.getFieldWithDefault(msg, 11, ""),
    data: msg.getData_asB64()
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
 * @return {!proto.evntsrc.event.Event}
 */
proto.evntsrc.event.Event.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.evntsrc.event.Event;
  return proto.evntsrc.event.Event.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.evntsrc.event.Event} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.evntsrc.event.Event}
 */
proto.evntsrc.event.Event.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = /** @type {number} */ (reader.readInt32());
      msg.setStream(value);
      break;
    case 3:
      var value = new proto.google.protobuf.Timestamp;
      reader.readMessage(value,proto.google.protobuf.Timestamp.deserializeBinaryFromReader);
      msg.setTime(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setType(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setTypeversion(value);
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.setCeversion(value);
      break;
    case 7:
      var value = /** @type {string} */ (reader.readString());
      msg.setSource(value);
      break;
    case 8:
      var value = /** @type {string} */ (reader.readString());
      msg.setSubject(value);
      break;
    case 9:
      var value = new proto.google.protobuf.Timestamp;
      reader.readMessage(value,proto.google.protobuf.Timestamp.deserializeBinaryFromReader);
      msg.setAcknowledged(value);
      break;
    case 10:
      var value = msg.getMetadataMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readString, null, "");
         });
      break;
    case 11:
      var value = /** @type {string} */ (reader.readString());
      msg.setContenttype(value);
      break;
    case 12:
      var value = /** @type {!Uint8Array} */ (reader.readBytes());
      msg.setData(value);
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
proto.evntsrc.event.Event.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.evntsrc.event.Event.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.evntsrc.event.Event} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.evntsrc.event.Event.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getStream();
  if (f !== 0) {
    writer.writeInt32(
      2,
      f
    );
  }
  f = message.getTime();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.google.protobuf.Timestamp.serializeBinaryToWriter
    );
  }
  f = message.getType();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
  f = message.getTypeversion();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getCeversion();
  if (f.length > 0) {
    writer.writeString(
      6,
      f
    );
  }
  f = message.getSource();
  if (f.length > 0) {
    writer.writeString(
      7,
      f
    );
  }
  f = message.getSubject();
  if (f.length > 0) {
    writer.writeString(
      8,
      f
    );
  }
  f = message.getAcknowledged();
  if (f != null) {
    writer.writeMessage(
      9,
      f,
      proto.google.protobuf.Timestamp.serializeBinaryToWriter
    );
  }
  f = message.getMetadataMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(10, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeString);
  }
  f = message.getContenttype();
  if (f.length > 0) {
    writer.writeString(
      11,
      f
    );
  }
  f = message.getData_asU8();
  if (f.length > 0) {
    writer.writeBytes(
      12,
      f
    );
  }
};


/**
 * optional string ID = 1;
 * @return {string}
 */
proto.evntsrc.event.Event.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.evntsrc.event.Event.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional int32 Stream = 2;
 * @return {number}
 */
proto.evntsrc.event.Event.prototype.getStream = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/** @param {number} value */
proto.evntsrc.event.Event.prototype.setStream = function(value) {
  jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional google.protobuf.Timestamp Time = 3;
 * @return {?proto.google.protobuf.Timestamp}
 */
proto.evntsrc.event.Event.prototype.getTime = function() {
  return /** @type{?proto.google.protobuf.Timestamp} */ (
    jspb.Message.getWrapperField(this, proto.google.protobuf.Timestamp, 3));
};


/** @param {?proto.google.protobuf.Timestamp|undefined} value */
proto.evntsrc.event.Event.prototype.setTime = function(value) {
  jspb.Message.setWrapperField(this, 3, value);
};


proto.evntsrc.event.Event.prototype.clearTime = function() {
  this.setTime(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.evntsrc.event.Event.prototype.hasTime = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional string Type = 4;
 * @return {string}
 */
proto.evntsrc.event.Event.prototype.getType = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/** @param {string} value */
proto.evntsrc.event.Event.prototype.setType = function(value) {
  jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * optional string TypeVersion = 5;
 * @return {string}
 */
proto.evntsrc.event.Event.prototype.getTypeversion = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/** @param {string} value */
proto.evntsrc.event.Event.prototype.setTypeversion = function(value) {
  jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * optional string CEVersion = 6;
 * @return {string}
 */
proto.evntsrc.event.Event.prototype.getCeversion = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 6, ""));
};


/** @param {string} value */
proto.evntsrc.event.Event.prototype.setCeversion = function(value) {
  jspb.Message.setProto3StringField(this, 6, value);
};


/**
 * optional string Source = 7;
 * @return {string}
 */
proto.evntsrc.event.Event.prototype.getSource = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 7, ""));
};


/** @param {string} value */
proto.evntsrc.event.Event.prototype.setSource = function(value) {
  jspb.Message.setProto3StringField(this, 7, value);
};


/**
 * optional string Subject = 8;
 * @return {string}
 */
proto.evntsrc.event.Event.prototype.getSubject = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 8, ""));
};


/** @param {string} value */
proto.evntsrc.event.Event.prototype.setSubject = function(value) {
  jspb.Message.setProto3StringField(this, 8, value);
};


/**
 * optional google.protobuf.Timestamp Acknowledged = 9;
 * @return {?proto.google.protobuf.Timestamp}
 */
proto.evntsrc.event.Event.prototype.getAcknowledged = function() {
  return /** @type{?proto.google.protobuf.Timestamp} */ (
    jspb.Message.getWrapperField(this, proto.google.protobuf.Timestamp, 9));
};


/** @param {?proto.google.protobuf.Timestamp|undefined} value */
proto.evntsrc.event.Event.prototype.setAcknowledged = function(value) {
  jspb.Message.setWrapperField(this, 9, value);
};


proto.evntsrc.event.Event.prototype.clearAcknowledged = function() {
  this.setAcknowledged(undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.evntsrc.event.Event.prototype.hasAcknowledged = function() {
  return jspb.Message.getField(this, 9) != null;
};


/**
 * map<string, string> Metadata = 10;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,string>}
 */
proto.evntsrc.event.Event.prototype.getMetadataMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,string>} */ (
      jspb.Message.getMapField(this, 10, opt_noLazyCreate,
      null));
};


proto.evntsrc.event.Event.prototype.clearMetadataMap = function() {
  this.getMetadataMap().clear();
};


/**
 * optional string ContentType = 11;
 * @return {string}
 */
proto.evntsrc.event.Event.prototype.getContenttype = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 11, ""));
};


/** @param {string} value */
proto.evntsrc.event.Event.prototype.setContenttype = function(value) {
  jspb.Message.setProto3StringField(this, 11, value);
};


/**
 * optional bytes Data = 12;
 * @return {!(string|Uint8Array)}
 */
proto.evntsrc.event.Event.prototype.getData = function() {
  return /** @type {!(string|Uint8Array)} */ (jspb.Message.getFieldWithDefault(this, 12, ""));
};


/**
 * optional bytes Data = 12;
 * This is a type-conversion wrapper around `getData()`
 * @return {string}
 */
proto.evntsrc.event.Event.prototype.getData_asB64 = function() {
  return /** @type {string} */ (jspb.Message.bytesAsB64(
      this.getData()));
};


/**
 * optional bytes Data = 12;
 * Note that Uint8Array is not supported on all browsers.
 * @see http://caniuse.com/Uint8Array
 * This is a type-conversion wrapper around `getData()`
 * @return {!Uint8Array}
 */
proto.evntsrc.event.Event.prototype.getData_asU8 = function() {
  return /** @type {!Uint8Array} */ (jspb.Message.bytesAsU8(
      this.getData()));
};


/** @param {!(string|Uint8Array)} value */
proto.evntsrc.event.Event.prototype.setData = function(value) {
  jspb.Message.setProto3BytesField(this, 12, value);
};

