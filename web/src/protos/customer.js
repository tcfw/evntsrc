/**
 * @fileoverview
 * @enhanceable
 * @suppress {messageConventions} JS Compiler reports an error if a variable or
 *     field starts with 'MSG_' and isn't a translatable message.
 * @public
 */
// GENERATED CODE -- DO NOT EDIT!

goog.provide('proto.evntsrc.billing.Customer');

goog.require('jspb.BinaryReader');
goog.require('jspb.BinaryWriter');
goog.require('jspb.Map');
goog.require('jspb.Message');
goog.require('proto.evntsrc.billing.Source');
goog.require('proto.evntsrc.billing.Subscription');


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
proto.evntsrc.billing.Customer = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.evntsrc.billing.Customer.repeatedFields_, null);
};
goog.inherits(proto.evntsrc.billing.Customer, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.evntsrc.billing.Customer.displayName = 'proto.evntsrc.billing.Customer';
}
/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.evntsrc.billing.Customer.repeatedFields_ = [15,16];



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
proto.evntsrc.billing.Customer.prototype.toObject = function(opt_includeInstance) {
  return proto.evntsrc.billing.Customer.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.evntsrc.billing.Customer} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.evntsrc.billing.Customer.toObject = function(includeInstance, msg) {
  var f, obj = {
    accountbalance: jspb.Message.getFieldWithDefault(msg, 1, 0),
    created: jspb.Message.getFieldWithDefault(msg, 2, 0),
    currency: jspb.Message.getFieldWithDefault(msg, 3, ""),
    deleted: jspb.Message.getFieldWithDefault(msg, 5, false),
    delinquent: jspb.Message.getFieldWithDefault(msg, 6, false),
    description: jspb.Message.getFieldWithDefault(msg, 7, ""),
    email: jspb.Message.getFieldWithDefault(msg, 9, ""),
    id: jspb.Message.getFieldWithDefault(msg, 10, ""),
    invoiceprefix: jspb.Message.getFieldWithDefault(msg, 11, ""),
    livemode: jspb.Message.getFieldWithDefault(msg, 12, false),
    metadataMap: (f = msg.getMetadataMap()) ? f.toObject(includeInstance, undefined) : [],
    sourcesList: jspb.Message.toObjectList(msg.getSourcesList(),
    proto.evntsrc.billing.Source.toObject, includeInstance),
    subscriptionsList: jspb.Message.toObjectList(msg.getSubscriptionsList(),
    proto.evntsrc.billing.Subscription.toObject, includeInstance),
    taxinfo: jspb.Message.getFieldWithDefault(msg, 17, ""),
    taxinfoverification: jspb.Message.getFieldWithDefault(msg, 18, "")
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
 * @return {!proto.evntsrc.billing.Customer}
 */
proto.evntsrc.billing.Customer.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.evntsrc.billing.Customer;
  return proto.evntsrc.billing.Customer.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.evntsrc.billing.Customer} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.evntsrc.billing.Customer}
 */
proto.evntsrc.billing.Customer.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setAccountbalance(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setCreated(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setCurrency(value);
      break;
    case 5:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setDeleted(value);
      break;
    case 6:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setDelinquent(value);
      break;
    case 7:
      var value = /** @type {string} */ (reader.readString());
      msg.setDescription(value);
      break;
    case 9:
      var value = /** @type {string} */ (reader.readString());
      msg.setEmail(value);
      break;
    case 10:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 11:
      var value = /** @type {string} */ (reader.readString());
      msg.setInvoiceprefix(value);
      break;
    case 12:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setLivemode(value);
      break;
    case 13:
      var value = msg.getMetadataMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readString, null, "");
         });
      break;
    case 15:
      var value = new proto.evntsrc.billing.Source;
      reader.readMessage(value,proto.evntsrc.billing.Source.deserializeBinaryFromReader);
      msg.addSources(value);
      break;
    case 16:
      var value = new proto.evntsrc.billing.Subscription;
      reader.readMessage(value,proto.evntsrc.billing.Subscription.deserializeBinaryFromReader);
      msg.addSubscriptions(value);
      break;
    case 17:
      var value = /** @type {string} */ (reader.readString());
      msg.setTaxinfo(value);
      break;
    case 18:
      var value = /** @type {string} */ (reader.readString());
      msg.setTaxinfoverification(value);
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
proto.evntsrc.billing.Customer.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.evntsrc.billing.Customer.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.evntsrc.billing.Customer} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.evntsrc.billing.Customer.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getAccountbalance();
  if (f !== 0) {
    writer.writeInt64(
      1,
      f
    );
  }
  f = message.getCreated();
  if (f !== 0) {
    writer.writeInt64(
      2,
      f
    );
  }
  f = message.getCurrency();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getDeleted();
  if (f) {
    writer.writeBool(
      5,
      f
    );
  }
  f = message.getDelinquent();
  if (f) {
    writer.writeBool(
      6,
      f
    );
  }
  f = message.getDescription();
  if (f.length > 0) {
    writer.writeString(
      7,
      f
    );
  }
  f = message.getEmail();
  if (f.length > 0) {
    writer.writeString(
      9,
      f
    );
  }
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      10,
      f
    );
  }
  f = message.getInvoiceprefix();
  if (f.length > 0) {
    writer.writeString(
      11,
      f
    );
  }
  f = message.getLivemode();
  if (f) {
    writer.writeBool(
      12,
      f
    );
  }
  f = message.getMetadataMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(13, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeString);
  }
  f = message.getSourcesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      15,
      f,
      proto.evntsrc.billing.Source.serializeBinaryToWriter
    );
  }
  f = message.getSubscriptionsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      16,
      f,
      proto.evntsrc.billing.Subscription.serializeBinaryToWriter
    );
  }
  f = message.getTaxinfo();
  if (f.length > 0) {
    writer.writeString(
      17,
      f
    );
  }
  f = message.getTaxinfoverification();
  if (f.length > 0) {
    writer.writeString(
      18,
      f
    );
  }
};


/**
 * optional int64 accountBalance = 1;
 * @return {number}
 */
proto.evntsrc.billing.Customer.prototype.getAccountbalance = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/** @param {number} value */
proto.evntsrc.billing.Customer.prototype.setAccountbalance = function(value) {
  jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional int64 created = 2;
 * @return {number}
 */
proto.evntsrc.billing.Customer.prototype.getCreated = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/** @param {number} value */
proto.evntsrc.billing.Customer.prototype.setCreated = function(value) {
  jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional string currency = 3;
 * @return {string}
 */
proto.evntsrc.billing.Customer.prototype.getCurrency = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.evntsrc.billing.Customer.prototype.setCurrency = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional bool deleted = 5;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.evntsrc.billing.Customer.prototype.getDeleted = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 5, false));
};


/** @param {boolean} value */
proto.evntsrc.billing.Customer.prototype.setDeleted = function(value) {
  jspb.Message.setProto3BooleanField(this, 5, value);
};


/**
 * optional bool delinquent = 6;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.evntsrc.billing.Customer.prototype.getDelinquent = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 6, false));
};


/** @param {boolean} value */
proto.evntsrc.billing.Customer.prototype.setDelinquent = function(value) {
  jspb.Message.setProto3BooleanField(this, 6, value);
};


/**
 * optional string description = 7;
 * @return {string}
 */
proto.evntsrc.billing.Customer.prototype.getDescription = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 7, ""));
};


/** @param {string} value */
proto.evntsrc.billing.Customer.prototype.setDescription = function(value) {
  jspb.Message.setProto3StringField(this, 7, value);
};


/**
 * optional string email = 9;
 * @return {string}
 */
proto.evntsrc.billing.Customer.prototype.getEmail = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 9, ""));
};


/** @param {string} value */
proto.evntsrc.billing.Customer.prototype.setEmail = function(value) {
  jspb.Message.setProto3StringField(this, 9, value);
};


/**
 * optional string id = 10;
 * @return {string}
 */
proto.evntsrc.billing.Customer.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 10, ""));
};


/** @param {string} value */
proto.evntsrc.billing.Customer.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 10, value);
};


/**
 * optional string invoicePrefix = 11;
 * @return {string}
 */
proto.evntsrc.billing.Customer.prototype.getInvoiceprefix = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 11, ""));
};


/** @param {string} value */
proto.evntsrc.billing.Customer.prototype.setInvoiceprefix = function(value) {
  jspb.Message.setProto3StringField(this, 11, value);
};


/**
 * optional bool liveMode = 12;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.evntsrc.billing.Customer.prototype.getLivemode = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 12, false));
};


/** @param {boolean} value */
proto.evntsrc.billing.Customer.prototype.setLivemode = function(value) {
  jspb.Message.setProto3BooleanField(this, 12, value);
};


/**
 * map<string, string> metadata = 13;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,string>}
 */
proto.evntsrc.billing.Customer.prototype.getMetadataMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,string>} */ (
      jspb.Message.getMapField(this, 13, opt_noLazyCreate,
      null));
};


proto.evntsrc.billing.Customer.prototype.clearMetadataMap = function() {
  this.getMetadataMap().clear();
};


/**
 * repeated Source sources = 15;
 * @return {!Array<!proto.evntsrc.billing.Source>}
 */
proto.evntsrc.billing.Customer.prototype.getSourcesList = function() {
  return /** @type{!Array<!proto.evntsrc.billing.Source>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.evntsrc.billing.Source, 15));
};


/** @param {!Array<!proto.evntsrc.billing.Source>} value */
proto.evntsrc.billing.Customer.prototype.setSourcesList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 15, value);
};


/**
 * @param {!proto.evntsrc.billing.Source=} opt_value
 * @param {number=} opt_index
 * @return {!proto.evntsrc.billing.Source}
 */
proto.evntsrc.billing.Customer.prototype.addSources = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 15, opt_value, proto.evntsrc.billing.Source, opt_index);
};


proto.evntsrc.billing.Customer.prototype.clearSourcesList = function() {
  this.setSourcesList([]);
};


/**
 * repeated Subscription subscriptions = 16;
 * @return {!Array<!proto.evntsrc.billing.Subscription>}
 */
proto.evntsrc.billing.Customer.prototype.getSubscriptionsList = function() {
  return /** @type{!Array<!proto.evntsrc.billing.Subscription>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.evntsrc.billing.Subscription, 16));
};


/** @param {!Array<!proto.evntsrc.billing.Subscription>} value */
proto.evntsrc.billing.Customer.prototype.setSubscriptionsList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 16, value);
};


/**
 * @param {!proto.evntsrc.billing.Subscription=} opt_value
 * @param {number=} opt_index
 * @return {!proto.evntsrc.billing.Subscription}
 */
proto.evntsrc.billing.Customer.prototype.addSubscriptions = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 16, opt_value, proto.evntsrc.billing.Subscription, opt_index);
};


proto.evntsrc.billing.Customer.prototype.clearSubscriptionsList = function() {
  this.setSubscriptionsList([]);
};


/**
 * optional string taxInfo = 17;
 * @return {string}
 */
proto.evntsrc.billing.Customer.prototype.getTaxinfo = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 17, ""));
};


/** @param {string} value */
proto.evntsrc.billing.Customer.prototype.setTaxinfo = function(value) {
  jspb.Message.setProto3StringField(this, 17, value);
};


/**
 * optional string taxInfoVerification = 18;
 * @return {string}
 */
proto.evntsrc.billing.Customer.prototype.getTaxinfoverification = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 18, ""));
};


/** @param {string} value */
proto.evntsrc.billing.Customer.prototype.setTaxinfoverification = function(value) {
  jspb.Message.setProto3StringField(this, 18, value);
};


