import React, { useContext, useEffect, useState } from "react"
import {
  Avatar,
  HStack,
  Status,
  Icon,
  Link,
  Portal,
  Stack,
  Text,
  Button,
} from "@chakra-ui/react"
import { Card } from "@chakra-ui/react"
import { VaultContext } from "./VaultContext"
import { LuRefreshCcw, LuChartLine, LuGitCommitHorizontal, LuLockOpen, LuLock } from "react-icons/lu"
import { BsPersonStanding, BsPersonDown, BsGit } from "react-icons/bs"
interface VaultStatusCardProps {
  initialVaultUnsealStatus?: VaultUnsealStatus
}

interface VaultUnsealStatus {
  initialized: boolean
  sealed: boolean
  standby: boolean
  server_time_utc: number
  version: string
  cluster_name: string
  cluster_id: string
  last_wal: number
}








/**
 * This one card should not be a hover card, but just  card instead
 * @param param0 
 * @returns 
 */
export const VaultStatusCard = (/* {  initialVaultUnsealStatus }: VaultStatusCardProps*/): React.JSX.Element => {
  const vaultContext = useContext(VaultContext);
  // const [vaultStatus, SetVaultStatus] = useState<VaultStatus>(initialVaultStatus)
  const fetchVaultUnsealStatus = async (signal: AbortSignal) => {
    console.info(` just calling fetchVaultUnsealStatus now`)
    const fetchVaultUnsealStatusResponsePromise = await fetch('/vault-status', {
      method: 'GET',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      },
      signal: signal
    }).catch((error) => {
      // alert(`fetchVaultUnsealStatus //   ERROR raw  fetching API raw Error: [${error}]`)
      // alert(`fetchVaultUnsealStatus //   ERROR  fetching API: [${JSON.stringify({ error: error }, null, 2)}]`)
      console.log(`fetchVaultUnsealStatus //   ERROR  fetching API: [${JSON.stringify({ error: error }, null, 2)}]`)
    });
    const fetchVaultUnsealStatusResponse = await fetchVaultUnsealStatusResponsePromise?.json();
    // --- 
    // 
    // console.log(fetchVaultUnsealStatusResponse);
    console.log(` fetchVaultUnsealStatus //  Here is the  fetchVaultUnsealStatusResponse = [${JSON.stringify({ fetchVaultUnsealStatusResponse: fetchVaultUnsealStatusResponse }, null, 2)}] `);
    const newStatus: VaultUnsealStatus = {
      cluster_id: fetchVaultUnsealStatusResponse.cluster_id,
      cluster_name: fetchVaultUnsealStatusResponse.cluster_name,
      initialized: fetchVaultUnsealStatusResponse.initialized,
      last_wal: fetchVaultUnsealStatusResponse.last_wal,
      sealed: fetchVaultUnsealStatusResponse.sealed,
      server_time_utc: fetchVaultUnsealStatusResponse.server_time_utc,
      standby: fetchVaultUnsealStatusResponse.standby,
      version: fetchVaultUnsealStatusResponse.version,
    }
    console.log(` fetchVaultUnsealStatus // JUST BEFORE CALLING SetVaultUnsealStatus(newStatus) Here is the  newStatus = [${JSON.stringify({ newStatus: newStatus }, null, 2)}] `);
    
    vaultContext.setVaultStatus(newStatus)
  }
  useEffect(() => {
    const abortController = new AbortController();
    const signal = abortController.signal; 
    fetchVaultUnsealStatus(signal)
    return () => {
      abortController.abort()
    }
  }, [/* vaultStatus.ClusterID, vaultStatus.ServerTimeUTC */ ])
  /*
  const [vaultStatus, SetVaultUnsealStatus] = useState<VaultUnsealStatus>({
    ClusterID: '',
    ClusterName: '',
    Initialized: false,
    LastWAL: -1,
    Sealed: false,
    ServerTimeUTC: 0,
    Standby: false,
    Version: "-1"
  })
    */
  return (
    <>

<Card.Root width="320px">
      <Card.Body gap="2">
        <Avatar.Root size="lg" shape="rounded">
          <Avatar.Image src="https://picsum.photos/200/300" />
          <Avatar.Fallback name="Nue Camp" />
        </Avatar.Root>
        <Card.Title mt="2">Vault Status Card</Card.Title>
        <Card.Description>
          This is the Vault Status Card body. Lorem ipsum dolor sit amet, consectetur
          adipiscing elit. Curabitur nec odio vel dui euismod fermentum.
          Curabitur nec odio vel dui euismod fermentum.
        </Card.Description>
        <Stack gap="4" direction="row">
                <Avatar.Root>
                  <Avatar.Image src="https://avatars.githubusercontent.com/u/152585220?s=48&v=4" />
                  <Avatar.Fallback name="Chakra UI" />
                </Avatar.Root>
                <Stack gap="3">
                  <Stack gap="1">
                    <Text textStyle="sm" fontWeight="semibold">
                    Vault Status Card
                    </Text>
                    <Text textStyle="sm" color="fg.muted">
                      The most powerful Vault Status Card.
                    </Text>
                  </Stack>
                  <HStack color="fg.subtle">
                    <BsGit className="sm" />
                    <Text textStyle="xs">{vaultContext.vaultStatus.version}</Text>
                  </HStack>
                  <HStack color="fg.subtle">
                    <Status.Root colorPalette={vaultContext.vaultStatus.initialized?"green":"orange"}>
                      <Status.Indicator />
                      Initialized : {vaultContext.vaultStatus.initialized?"Yes":"No"}
                    </Status.Root>
                    <Status.Root colorPalette={vaultContext.vaultStatus.sealed?"green":"red"}>
                      <Status.Indicator />
                      {vaultContext.vaultStatus.sealed?<LuLock/>:<LuLockOpen />}
                      Sealed: {vaultContext.vaultStatus.sealed?"Yes":"No"}
                    </Status.Root>
                  </HStack>
                  <HStack color="fg.subtle">
                    <Status.Root colorPalette={vaultContext.vaultStatus.standby?"blue":"gray"}>
                      <Status.Indicator />
                      {vaultContext.vaultStatus.standby?<BsPersonStanding />:<BsPersonDown />}
                      Standby
                    </Status.Root>
                    <Status.Root colorPalette={"blue"}>
                      <Status.Indicator />
                      {vaultContext.vaultStatus.server_time_utc?<BsPersonStanding />:<BsPersonDown />}
                      ServerTimeUTC: {vaultContext.vaultStatus.server_time_utc}
                    </Status.Root>
                  </HStack>
                  <HStack color="fg:subtle" >

                  <Button title="Refresh status" variant="solid" type="button" size="md" onClick={async () => {
        console.log(`Onclick jbl Trigger fetch vaultStatus.Version=[${vaultContext.vaultStatus.version}]`)
        // setShowScanner(true)
        const abortController = new AbortController()
        const signal = abortController.signal
        await fetchVaultUnsealStatus(signal)
      }}>
        <LuRefreshCcw /> Refresh status
      </Button>
                  </HStack>
                </Stack>
              </Stack>
      </Card.Body>
      <Card.Footer justifyContent="flex-end">
        <Button variant="outline">View</Button>
        <Button>Join</Button>
      </Card.Footer>
    </Card.Root>







      
    </>
  )
}












