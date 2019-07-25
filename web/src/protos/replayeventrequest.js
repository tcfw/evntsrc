/**
 * @fileoverview
 * @enhanceable
 * @suppress {messageConventions} JS Compiler reports an error if a variable or
 *     field starts with 'MSG_' and isn't a translatable message.
 * @public
 */
// GENERATED CODE -- DO NOT EDIT!

goog.provide('proto.evntsrc.storer.ReplayEventRequest');

goog.require('jspb.BinaryReader');
goog.require('jspb.BinaryWriter');
goog.require('jspb.Message');


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
proto.evntsrc.storer.ReplayEventRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.evntsrc.storer.ReplayEventRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.evntsrc.storer.ReplayEventRequest.displayName = 'proto.evntsrc.storer.ReplayEventRequest';
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
proto.evntsrc.storer.ReplayEventRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.evntsrc.storer.ReplayEventRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.evntsrc.storer.ReplayEventRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.evntsrc.storer.ReplayEventRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    stream: jspb.Message.getFieldWithDefault(msg, 1, 0),
    eventid: jspb.Message.getFieldWithDefault(msg, 2, "")
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
 * @return {!proto.evntsrc.storer.ReplayEventRequest}
 */
proto.evntsrc.storer.ReplayEventRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.evntsrc.storer.ReplayEventRequest;
  return proto.evntsrc.storer.ReplayEventRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.evntsrc.storer.ReplayEventRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.evntsrc.storer.ReplayEventRequest}
 */
proto.evntsrc.storer.ReplayEventRequest.deserializeBinaryFromReader = function(msg, reader) {
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
      var value = /** @type {string} */ (reader.readString());
      msg.setEventid(value);
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
proto.evntsrc.storer.ReplayEventRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.evntsrc.storer.ReplayEventRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.evntsrc.storer.ReplayEventRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.evntsrc.storer.ReplayEventRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getStream();
  if (f !== 0) {
    writer.writeInt32(
      1,
      f
    );
  }
  f = message.getEventid();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional int32 Stream = 1;
 * @return {number}
 */
proto.evntsrc.storer.ReplayEventRequest.prototype.getStream = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/** @param {number} value */
proto.evntsrc.storer.ReplayEventRequest.prototype.setStream = function(value) {
  jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional string EventID = 2;
 * @return {string}
 */
proto.evntsrc.storer.ReplayEventRequest.prototype.getEventid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.evntsrc.storer.ReplayEventRequest.prototype.setEventid = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


