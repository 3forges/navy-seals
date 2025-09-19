export interface VaultSealStatus {
  type: string,
  initialized: boolean,
  sealed: boolean,
  progress: number,
  nonce: string,
  version: string,
  cluster_name: string,
  cluster_id: string,
  storage_type: string,
  warnings: string[],
}

/*

In the Golang API:

type SealStatusResponse struct {
	Type         string   `json:"type"`
	Initialized  bool     `json:"initialized"`
	Sealed       bool     `json:"sealed"`
	T            int      `json:"t"`
	N            int      `json:"n"`
	Progress     int      `json:"progress"`
	Nonce        string   `json:"nonce"`
	Version      string   `json:"version"`
	BuildDate    string   `json:"build_date"`
	Migration    bool     `json:"migration"`
	ClusterName  string   `json:"cluster_name,omitempty"`
	ClusterID    string   `json:"cluster_id,omitempty"`
	RecoverySeal bool     `json:"recovery_seal"`
	StorageType  string   `json:"storage_type,omitempty"`
	Warnings     []string `json:"warnings,omitempty"`
}


*/
