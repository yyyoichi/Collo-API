// @generated by protoc-gen-connect-es v1.1.3 with parameter "target=ts"
// @generated from file api/v2/collo.proto (package api.v2, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import { ColloNetworkStreamRequest, ColloNetworkStreamResponse, ColloWebStreamRequest, ColloWebStreamResponse } from "./collo_pb.js";
import { MethodKind } from "@bufbuild/protobuf";

/**
 * @generated from service api.v2.ColloNetworkService
 */
export const ColloNetworkService = {
  typeName: "api.v2.ColloNetworkService",
  methods: {
    /**
     * @generated from rpc api.v2.ColloNetworkService.ColloNetworkStream
     */
    colloNetworkStream: {
      name: "ColloNetworkStream",
      I: ColloNetworkStreamRequest,
      O: ColloNetworkStreamResponse,
      kind: MethodKind.BiDiStreaming,
    },
  }
} as const;

/**
 * @generated from service api.v2.ColloWebService
 */
export const ColloWebService = {
  typeName: "api.v2.ColloWebService",
  methods: {
    /**
     * @generated from rpc api.v2.ColloWebService.ColloWebStream
     */
    colloWebStream: {
      name: "ColloWebStream",
      I: ColloWebStreamRequest,
      O: ColloWebStreamResponse,
      kind: MethodKind.ServerStreaming,
    },
  }
} as const;

