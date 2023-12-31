// @generated by protoc-gen-es v1.6.0 with parameter "target=ts"
// @generated from file api/v2/collo.proto (package api.v2, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3, Timestamp } from "@bufbuild/protobuf";

/**
 * @generated from message api.v2.ColloRateWebStreamRequest
 */
export class ColloRateWebStreamRequest extends Message<ColloRateWebStreamRequest> {
  /**
   * @generated from field: string keyword = 1;
   */
  keyword = "";

  /**
   * @generated from field: google.protobuf.Timestamp from = 2;
   */
  from?: Timestamp;

  /**
   * @generated from field: google.protobuf.Timestamp until = 3;
   */
  until?: Timestamp;

  /**
   * @generated from field: uint32 forcus_node_id = 4;
   */
  forcusNodeId = 0;

  /**
   * @generated from field: repeated uint32 part_of_speech_types = 5;
   */
  partOfSpeechTypes: number[] = [];

  /**
   * @generated from field: repeated string stopwords = 6;
   */
  stopwords: string[] = [];

  constructor(data?: PartialMessage<ColloRateWebStreamRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "api.v2.ColloRateWebStreamRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "keyword", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "from", kind: "message", T: Timestamp },
    { no: 3, name: "until", kind: "message", T: Timestamp },
    { no: 4, name: "forcus_node_id", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
    { no: 5, name: "part_of_speech_types", kind: "scalar", T: 13 /* ScalarType.UINT32 */, repeated: true },
    { no: 6, name: "stopwords", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ColloRateWebStreamRequest {
    return new ColloRateWebStreamRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ColloRateWebStreamRequest {
    return new ColloRateWebStreamRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ColloRateWebStreamRequest {
    return new ColloRateWebStreamRequest().fromJsonString(jsonString, options);
  }

  static equals(a: ColloRateWebStreamRequest | PlainMessage<ColloRateWebStreamRequest> | undefined, b: ColloRateWebStreamRequest | PlainMessage<ColloRateWebStreamRequest> | undefined): boolean {
    return proto3.util.equals(ColloRateWebStreamRequest, a, b);
  }
}

/**
 * @generated from message api.v2.ColloRateWebStreamResponse
 */
export class ColloRateWebStreamResponse extends Message<ColloRateWebStreamResponse> {
  /**
   * @generated from field: repeated api.v2.RateNode nodes = 1;
   */
  nodes: RateNode[] = [];

  /**
   * @generated from field: repeated api.v2.RateEdge edges = 2;
   */
  edges: RateEdge[] = [];

  /**
   * @generated from field: uint32 dones = 3;
   */
  dones = 0;

  /**
   * @generated from field: uint32 needs = 4;
   */
  needs = 0;

  constructor(data?: PartialMessage<ColloRateWebStreamResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "api.v2.ColloRateWebStreamResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "nodes", kind: "message", T: RateNode, repeated: true },
    { no: 2, name: "edges", kind: "message", T: RateEdge, repeated: true },
    { no: 3, name: "dones", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
    { no: 4, name: "needs", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ColloRateWebStreamResponse {
    return new ColloRateWebStreamResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ColloRateWebStreamResponse {
    return new ColloRateWebStreamResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ColloRateWebStreamResponse {
    return new ColloRateWebStreamResponse().fromJsonString(jsonString, options);
  }

  static equals(a: ColloRateWebStreamResponse | PlainMessage<ColloRateWebStreamResponse> | undefined, b: ColloRateWebStreamResponse | PlainMessage<ColloRateWebStreamResponse> | undefined): boolean {
    return proto3.util.equals(ColloRateWebStreamResponse, a, b);
  }
}

/**
 * @generated from message api.v2.RateNode
 */
export class RateNode extends Message<RateNode> {
  /**
   * @generated from field: uint32 node_id = 1;
   */
  nodeId = 0;

  /**
   * @generated from field: string word = 2;
   */
  word = "";

  /**
   * @generated from field: float rate = 3;
   */
  rate = 0;

  constructor(data?: PartialMessage<RateNode>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "api.v2.RateNode";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "node_id", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
    { no: 2, name: "word", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "rate", kind: "scalar", T: 2 /* ScalarType.FLOAT */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): RateNode {
    return new RateNode().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): RateNode {
    return new RateNode().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): RateNode {
    return new RateNode().fromJsonString(jsonString, options);
  }

  static equals(a: RateNode | PlainMessage<RateNode> | undefined, b: RateNode | PlainMessage<RateNode> | undefined): boolean {
    return proto3.util.equals(RateNode, a, b);
  }
}

/**
 * @generated from message api.v2.RateEdge
 */
export class RateEdge extends Message<RateEdge> {
  /**
   * @generated from field: uint32 edge_id = 1;
   */
  edgeId = 0;

  /**
   * @generated from field: uint32 node_id1 = 2;
   */
  nodeId1 = 0;

  /**
   * @generated from field: uint32 node_id2 = 3;
   */
  nodeId2 = 0;

  /**
   * @generated from field: float rate = 4;
   */
  rate = 0;

  constructor(data?: PartialMessage<RateEdge>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "api.v2.RateEdge";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "edge_id", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
    { no: 2, name: "node_id1", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
    { no: 3, name: "node_id2", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
    { no: 4, name: "rate", kind: "scalar", T: 2 /* ScalarType.FLOAT */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): RateEdge {
    return new RateEdge().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): RateEdge {
    return new RateEdge().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): RateEdge {
    return new RateEdge().fromJsonString(jsonString, options);
  }

  static equals(a: RateEdge | PlainMessage<RateEdge> | undefined, b: RateEdge | PlainMessage<RateEdge> | undefined): boolean {
    return proto3.util.equals(RateEdge, a, b);
  }
}

