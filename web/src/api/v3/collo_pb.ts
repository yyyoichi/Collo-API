// @generated by protoc-gen-es v1.6.0 with parameter "target=ts"
// @generated from file api/v3/collo.proto (package api.v3, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3, Timestamp } from "@bufbuild/protobuf";

/**
 * 共起ネットワーク情報を取得する。
 *
 * @generated from message api.v3.NetworkStreamRequest
 */
export class NetworkStreamRequest extends Message<NetworkStreamRequest> {
  /**
   * リクエスト設定
   *
   * @generated from field: api.v3.RequestConfig config = 1;
   */
  config?: RequestConfig;

  /**
   * 特定のノードに関連する共起情報を返す。Falsyの時もっとも重要な3つのノードとその関連する共起情報を返す。
   *
   * @generated from field: uint32 forcus_node_id = 2;
   */
  forcusNodeId = 0;

  constructor(data?: PartialMessage<NetworkStreamRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "api.v3.NetworkStreamRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "config", kind: "message", T: RequestConfig },
    { no: 2, name: "forcus_node_id", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): NetworkStreamRequest {
    return new NetworkStreamRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): NetworkStreamRequest {
    return new NetworkStreamRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): NetworkStreamRequest {
    return new NetworkStreamRequest().fromJsonString(jsonString, options);
  }

  static equals(a: NetworkStreamRequest | PlainMessage<NetworkStreamRequest> | undefined, b: NetworkStreamRequest | PlainMessage<NetworkStreamRequest> | undefined): boolean {
    return proto3.util.equals(NetworkStreamRequest, a, b);
  }
}

/**
 * 共起ネットワーク情報。とその進捗情報。
 *
 * @generated from message api.v3.NetworkStreamResponse
 */
export class NetworkStreamResponse extends Message<NetworkStreamResponse> {
  /**
   * @generated from field: repeated api.v3.Node nodes = 1;
   */
  nodes: Node[] = [];

  /**
   * @generated from field: repeated api.v3.Edge edges = 2;
   */
  edges: Edge[] = [];

  /**
   * @generated from field: api.v3.Meta meta = 3;
   */
  meta?: Meta;

  /**
   * 進捗
   *
   * @generated from field: float process = 4;
   */
  process = 0;

  constructor(data?: PartialMessage<NetworkStreamResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "api.v3.NetworkStreamResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "nodes", kind: "message", T: Node, repeated: true },
    { no: 2, name: "edges", kind: "message", T: Edge, repeated: true },
    { no: 3, name: "meta", kind: "message", T: Meta },
    { no: 4, name: "process", kind: "scalar", T: 2 /* ScalarType.FLOAT */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): NetworkStreamResponse {
    return new NetworkStreamResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): NetworkStreamResponse {
    return new NetworkStreamResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): NetworkStreamResponse {
    return new NetworkStreamResponse().fromJsonString(jsonString, options);
  }

  static equals(a: NetworkStreamResponse | PlainMessage<NetworkStreamResponse> | undefined, b: NetworkStreamResponse | PlainMessage<NetworkStreamResponse> | undefined): boolean {
    return proto3.util.equals(NetworkStreamResponse, a, b);
  }
}

/**
 * 共起行列中の任意の数のノードを重要度順に取得する。
 *
 * @generated from message api.v3.NodeRateStreamRequest
 */
export class NodeRateStreamRequest extends Message<NodeRateStreamRequest> {
  /**
   * リクエスト設定
   *
   * @generated from field: api.v3.RequestConfig config = 1;
   */
  config?: RequestConfig;

  /**
   * start with 0 デフォルト0。
   *
   * @generated from field: uint32 offset = 2;
   */
  offset = 0;

  /**
   * デフォルト100。 
   *
   * @generated from field: uint32 limit = 3;
   */
  limit = 0;

  constructor(data?: PartialMessage<NodeRateStreamRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "api.v3.NodeRateStreamRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "config", kind: "message", T: RequestConfig },
    { no: 2, name: "offset", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
    { no: 3, name: "limit", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): NodeRateStreamRequest {
    return new NodeRateStreamRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): NodeRateStreamRequest {
    return new NodeRateStreamRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): NodeRateStreamRequest {
    return new NodeRateStreamRequest().fromJsonString(jsonString, options);
  }

  static equals(a: NodeRateStreamRequest | PlainMessage<NodeRateStreamRequest> | undefined, b: NodeRateStreamRequest | PlainMessage<NodeRateStreamRequest> | undefined): boolean {
    return proto3.util.equals(NodeRateStreamRequest, a, b);
  }
}

/**
 * 共起行列中の任意の数のノードを重要度順に返す。
 *
 * @generated from message api.v3.NodeRateStreamResponse
 */
export class NodeRateStreamResponse extends Message<NodeRateStreamResponse> {
  /**
   * @generated from field: repeated api.v3.Node nodes = 1;
   */
  nodes: Node[] = [];

  /**
   * ノード保持総数。
   *
   * @generated from field: uint32 num = 2;
   */
  num = 0;

  /**
   * 次回取得位置。0のときは最後まで取得済みとみなす。
   *
   * @generated from field: uint32 next = 3;
   */
  next = 0;

  /**
   * 返却個数。
   *
   * @generated from field: uint32 count = 4;
   */
  count = 0;

  /**
   * @generated from field: api.v3.Meta meta = 5;
   */
  meta?: Meta;

  /**
   * 進捗
   *
   * @generated from field: float process = 6;
   */
  process = 0;

  constructor(data?: PartialMessage<NodeRateStreamResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "api.v3.NodeRateStreamResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "nodes", kind: "message", T: Node, repeated: true },
    { no: 2, name: "num", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
    { no: 3, name: "next", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
    { no: 4, name: "count", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
    { no: 5, name: "meta", kind: "message", T: Meta },
    { no: 6, name: "process", kind: "scalar", T: 2 /* ScalarType.FLOAT */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): NodeRateStreamResponse {
    return new NodeRateStreamResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): NodeRateStreamResponse {
    return new NodeRateStreamResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): NodeRateStreamResponse {
    return new NodeRateStreamResponse().fromJsonString(jsonString, options);
  }

  static equals(a: NodeRateStreamResponse | PlainMessage<NodeRateStreamResponse> | undefined, b: NodeRateStreamResponse | PlainMessage<NodeRateStreamResponse> | undefined): boolean {
    return proto3.util.equals(NodeRateStreamResponse, a, b);
  }
}

/**
 * @generated from message api.v3.RequestConfig
 */
export class RequestConfig extends Message<RequestConfig> {
  /**
   * API検索キーワード。必須。
   *
   * @generated from field: string keyword = 1;
   */
  keyword = "";

  /**
   * API検索開始日。必須。
   *
   * @generated from field: google.protobuf.Timestamp from = 2;
   */
  from?: Timestamp;

  /**
   * API検索終了日。必須。
   *
   * @generated from field: google.protobuf.Timestamp until = 3;
   */
  until?: Timestamp;

  /**
   * 行列使用品詞。オプション。
   *
   * @generated from field: repeated uint32 part_of_speech_types = 4;
   */
  partOfSpeechTypes: number[] = [];

  /**
   * 解析に使用しない単語。オプション。
   *
   * @generated from field: repeated string stopwords = 5;
   */
  stopwords: string[] = [];

  /**
   * 使用する共起行列。任意の文字列あるいは"total"(文書全体)、指定しない場合、すべてのグループと文書全体の共起行列を生成する。
   *
   * @generated from field: string forcus_group_id = 6;
   */
  forcusGroupId = "";

  /**
   * グルーピング手法。デフォルトは会議録ごと（1）。 
   *
   * @generated from field: api.v3.RequestConfig.PickGroupType pick_group_type = 7;
   */
  pickGroupType = RequestConfig_PickGroupType.UNSPECIFIED;

  /**
   * 使用するAPIの種類。デフォルトは、発言単位（1）。
   *
   * @generated from field: api.v3.RequestConfig.NdlApiType ndl_api_type = 8;
   */
  ndlApiType = RequestConfig_NdlApiType.UNSPECIFIED;

  /**
   * apiフェッチ時cacheを利用するか。デフォルトは使用しない。
   *
   * @generated from field: bool use_ndl_cache = 9;
   */
  useNdlCache = false;

  /**
   * apiフェッチ後、cacheを作成するか
   *
   * @generated from field: bool create_ndl_cache = 10;
   */
  createNdlCache = false;

  /**
   * apiフェッチのキャッシュ利用場所
   *
   * @generated from field: string ndl_cache_dir = 11;
   */
  ndlCacheDir = "";

  constructor(data?: PartialMessage<RequestConfig>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "api.v3.RequestConfig";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "keyword", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "from", kind: "message", T: Timestamp },
    { no: 3, name: "until", kind: "message", T: Timestamp },
    { no: 4, name: "part_of_speech_types", kind: "scalar", T: 13 /* ScalarType.UINT32 */, repeated: true },
    { no: 5, name: "stopwords", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 6, name: "forcus_group_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 7, name: "pick_group_type", kind: "enum", T: proto3.getEnumType(RequestConfig_PickGroupType) },
    { no: 8, name: "ndl_api_type", kind: "enum", T: proto3.getEnumType(RequestConfig_NdlApiType) },
    { no: 9, name: "use_ndl_cache", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
    { no: 10, name: "create_ndl_cache", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
    { no: 11, name: "ndl_cache_dir", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): RequestConfig {
    return new RequestConfig().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): RequestConfig {
    return new RequestConfig().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): RequestConfig {
    return new RequestConfig().fromJsonString(jsonString, options);
  }

  static equals(a: RequestConfig | PlainMessage<RequestConfig> | undefined, b: RequestConfig | PlainMessage<RequestConfig> | undefined): boolean {
    return proto3.util.equals(RequestConfig, a, b);
  }
}

/**
 * @generated from enum api.v3.RequestConfig.PickGroupType
 */
export enum RequestConfig_PickGroupType {
  /**
   * @generated from enum value: PICK_GROUP_TYPE_UNSPECIFIED = 0;
   */
  UNSPECIFIED = 0,

  /**
   * 会議録ごと
   *
   * @generated from enum value: PICK_GROUP_TYPE_ISSUEID = 1;
   */
  ISSUEID = 1,

  /**
   * 月ごと(ex: 2024-01-01)
   *
   * @generated from enum value: PICK_GROUP_TYPE_MONTH = 2;
   */
  MONTH = 2,
}
// Retrieve enum metadata with: proto3.getEnumType(RequestConfig_PickGroupType)
proto3.util.setEnumType(RequestConfig_PickGroupType, "api.v3.RequestConfig.PickGroupType", [
  { no: 0, name: "PICK_GROUP_TYPE_UNSPECIFIED" },
  { no: 1, name: "PICK_GROUP_TYPE_ISSUEID" },
  { no: 2, name: "PICK_GROUP_TYPE_MONTH" },
]);

/**
 * @generated from enum api.v3.RequestConfig.NdlApiType
 */
export enum RequestConfig_NdlApiType {
  /**
   * @generated from enum value: NDL_API_TYPE_UNSPECIFIED = 0;
   */
  UNSPECIFIED = 0,

  /**
   * 発言単位
   *
   * @generated from enum value: NDL_API_TYPE_SPEECH_UNSPECIFIED = 1;
   */
  SPEECH_UNSPECIFIED = 1,

  /**
   * 会議単位
   *
   * @generated from enum value: NDL_API_TYPE_MEETING = 2;
   */
  MEETING = 2,
}
// Retrieve enum metadata with: proto3.getEnumType(RequestConfig_NdlApiType)
proto3.util.setEnumType(RequestConfig_NdlApiType, "api.v3.RequestConfig.NdlApiType", [
  { no: 0, name: "NDL_API_TYPE_UNSPECIFIED" },
  { no: 1, name: "NDL_API_TYPE_SPEECH_UNSPECIFIED" },
  { no: 2, name: "NDL_API_TYPE_MEETING" },
]);

/**
 * @generated from message api.v3.Node
 */
export class Node extends Message<Node> {
  /**
   * @generated from field: uint32 node_id = 1;
   */
  nodeId = 0;

  /**
   * @generated from field: string word = 2;
   */
  word = "";

  /**
   * 中心度
   *
   * @generated from field: float rate = 3;
   */
  rate = 0;

  /**
   * 共起ワード数
   *
   * @generated from field: uint32 num_edges = 4;
   */
  numEdges = 0;

  constructor(data?: PartialMessage<Node>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "api.v3.Node";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "node_id", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
    { no: 2, name: "word", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "rate", kind: "scalar", T: 2 /* ScalarType.FLOAT */ },
    { no: 4, name: "num_edges", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Node {
    return new Node().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Node {
    return new Node().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Node {
    return new Node().fromJsonString(jsonString, options);
  }

  static equals(a: Node | PlainMessage<Node> | undefined, b: Node | PlainMessage<Node> | undefined): boolean {
    return proto3.util.equals(Node, a, b);
  }
}

/**
 * @generated from message api.v3.Edge
 */
export class Edge extends Message<Edge> {
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

  constructor(data?: PartialMessage<Edge>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "api.v3.Edge";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "edge_id", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
    { no: 2, name: "node_id1", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
    { no: 3, name: "node_id2", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
    { no: 4, name: "rate", kind: "scalar", T: 2 /* ScalarType.FLOAT */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Edge {
    return new Edge().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Edge {
    return new Edge().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Edge {
    return new Edge().fromJsonString(jsonString, options);
  }

  static equals(a: Edge | PlainMessage<Edge> | undefined, b: Edge | PlainMessage<Edge> | undefined): boolean {
    return proto3.util.equals(Edge, a, b);
  }
}

/**
 * @generated from message api.v3.Meta
 */
export class Meta extends Message<Meta> {
  /**
   * 共起行列が所属するグループ
   *
   * @generated from field: string group_id = 1;
   */
  groupId = "";

  /**
   * @generated from field: google.protobuf.Timestamp from = 2;
   */
  from?: Timestamp;

  /**
   * @generated from field: google.protobuf.Timestamp until = 3;
   */
  until?: Timestamp;

  /**
   * @generated from field: repeated api.v3.DocMeta metas = 4;
   */
  metas: DocMeta[] = [];

  constructor(data?: PartialMessage<Meta>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "api.v3.Meta";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "group_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "from", kind: "message", T: Timestamp },
    { no: 3, name: "until", kind: "message", T: Timestamp },
    { no: 4, name: "metas", kind: "message", T: DocMeta, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Meta {
    return new Meta().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Meta {
    return new Meta().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Meta {
    return new Meta().fromJsonString(jsonString, options);
  }

  static equals(a: Meta | PlainMessage<Meta> | undefined, b: Meta | PlainMessage<Meta> | undefined): boolean {
    return proto3.util.equals(Meta, a, b);
  }
}

/**
 * @generated from message api.v3.DocMeta
 */
export class DocMeta extends Message<DocMeta> {
  /**
   * @generated from field: string key = 1;
   */
  key = "";

  /**
   * @generated from field: string group_id = 2;
   */
  groupId = "";

  /**
   * @generated from field: string name = 3;
   */
  name = "";

  /**
   * @generated from field: google.protobuf.Timestamp at = 4;
   */
  at?: Timestamp;

  /**
   * @generated from field: string description = 5;
   */
  description = "";

  constructor(data?: PartialMessage<DocMeta>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "api.v3.DocMeta";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "key", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "group_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 4, name: "at", kind: "message", T: Timestamp },
    { no: 5, name: "description", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): DocMeta {
    return new DocMeta().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): DocMeta {
    return new DocMeta().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): DocMeta {
    return new DocMeta().fromJsonString(jsonString, options);
  }

  static equals(a: DocMeta | PlainMessage<DocMeta> | undefined, b: DocMeta | PlainMessage<DocMeta> | undefined): boolean {
    return proto3.util.equals(DocMeta, a, b);
  }
}
