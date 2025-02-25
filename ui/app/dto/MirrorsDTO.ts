import { FlowConnectionConfigs, QRepConfig } from '@/grpc_generated/flow';
import { Dispatch, SetStateAction } from 'react';

export type UCreateMirrorResponse = {
  created: boolean;
};

export type UValidateMirrorResponse = {
  ok: boolean;
  errorMessage: string;
};

export type UDropMirrorResponse = {
  dropped: boolean;
  errorMessage: string;
};

export type CDCConfig = FlowConnectionConfigs;
export type MirrorConfig = CDCConfig | QRepConfig;
export type MirrorSetter = Dispatch<SetStateAction<CDCConfig | QRepConfig>>;
export type TableMapRow = {
  schema: string;
  source: string;
  destination: string;
  partitionKey: string;
  exclude: Set<string>;
  selected: boolean;
  canMirror: boolean;
};

export type SyncStatusRow = {
  batchId: bigint;
  startTime: Date;
  endTime: Date | null;
  numRows: number;
};
