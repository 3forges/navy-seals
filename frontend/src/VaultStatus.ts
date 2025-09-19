export interface VaultStatus {
  initialized: boolean;
  sealed: boolean;
  standby: boolean;
  server_time_utc: number;
  version: string;
  cluster_name: string;
  cluster_id: string;
  last_wal: number;
}
