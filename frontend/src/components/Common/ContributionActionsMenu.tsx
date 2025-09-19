import { IconButton } from "@chakra-ui/react"
import { BsThreeDotsVertical } from "react-icons/bs"
import { MenuContent, MenuRoot, MenuTrigger } from "../ui/menu"

import type { ContributionPublic } from "@/client"
import DeleteContribution from "../Contributions/DeleteContribution"
import EditContribution from "../Contributions/EditContribution"
import React from "react"

interface ContributionActionsMenuProps {
  contribution: ContributionPublic
}

export const ContributionActionsMenu = ({ contribution }: ContributionActionsMenuProps): React.JSX.Element => {
  return (
    <MenuRoot>
      <MenuTrigger asChild>
        <IconButton variant="ghost" color="inherit">
          <BsThreeDotsVertical />
        </IconButton>
      </MenuTrigger>
      <MenuContent>
        <EditContribution contribution={contribution} />
        <DeleteContribution id={contribution.id} />
      </MenuContent>
    </MenuRoot>
  )
}
