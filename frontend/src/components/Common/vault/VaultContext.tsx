import React, { createContext, useState } from "react";
import { VaultStatus } from "../../../VaultStatus";
import { VaultSealStatus } from "../../../VaultSealStatus";

import {
    Button,
  } from "@chakra-ui/react"
// export const VaultContext = createContext<VaultStatus>({

export interface VaultContextState {
  vaultStatus: VaultStatus
  setVaultStatus: React.Dispatch<React.SetStateAction<VaultStatus>>
  vaultSealStatus: VaultSealStatus
  setVaultSealStatus: React.Dispatch<React.SetStateAction<VaultSealStatus>>
}
export const VaultContext = React.createContext<VaultContextState>({
    vaultStatus: {
        cluster_id: "",
        cluster_name: "",
        initialized: false,
        last_wal: 0,
        sealed: false,
        server_time_utc: 0,
        standby: false,
        version: "0.0.0"
    },
    setVaultSealStatus: () => {
        throw new Error("[setVaultSealStatus] - Not Implemented")
    },
    setVaultStatus: () => {
        throw new Error("[setVaultStatus] - Not Implemented")
    },
    vaultSealStatus: {
        cluster_id: "",
        cluster_name: "",
        initialized: false,
        nonce: "",
        progress: 0,
        sealed: false,
        storage_type: "",
        type: "",
        version: "0.0.0",
        warnings: []
    }
});


export const VaultContextProvider: React.FC<{children: React.ReactNode}> = ({ children }) => {
    const [ vaultStatus, setVaultStatus] = useState<VaultStatus>({
        cluster_id: '',
        cluster_name: '',
        initialized: false,
        last_wal: -1,
        sealed: false,
        server_time_utc: 0,
        standby: false,
        version: "-1"
      })
    const [ vaultSealStatus, setVaultSealStatus] = useState<VaultSealStatus>({
        cluster_id: "",
        cluster_name: "",
        initialized: false,
        nonce: "",
        progress: 0,
        sealed: false,
        storage_type: "",
        type: "",
        version: "0.0.0",
        warnings: []
    })
    
    return (
        <VaultContext.Provider value={{vaultStatus, setVaultStatus, vaultSealStatus, setVaultSealStatus}}>{children}</VaultContext.Provider>
    )

}

