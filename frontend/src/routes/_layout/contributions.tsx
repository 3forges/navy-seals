import {
  Container,
  EmptyState,
  Flex,
  Heading,
  Table,
  VStack,
} from "@chakra-ui/react"
import { useQuery } from "@tanstack/react-query"
import { createFileRoute, useNavigate } from "@tanstack/react-router"
import { FiSearch } from "react-icons/fi"
import { z } from "zod"

import { ContributionsService } from "@/client"
import { ContributionActionsMenu } from "@/components/Common/ContributionActionsMenu"
// import AddContribution from "@/components/Contributions/AddContribution"
// import AddContribution from "@/components/Contributions2/AddContribution"
import AddContribution from "@/components/Contributions3/AddContribution"
import PendingContributions from "@/components/Pending/PendingContributions"
import {
  PaginationItems,
  PaginationNextTrigger,
  PaginationPrevTrigger,
  PaginationRoot,
} from "@/components/ui/pagination.tsx"
import { ContributionPublic } from "../../client"
import React from "react"

const contributionsSearchSchema = z.object({
  page: z.number().catch(1),
})

const PER_PAGE = 5

function getContributionsQueryOptions({ page }: { page: number }) {
  return {
    queryFn: () =>
      ContributionsService.readContributions({ skip: (page - 1) * PER_PAGE, limit: PER_PAGE }),
    queryKey: ["contributions", { page }],
  }
}

export const Route = createFileRoute("/_layout/contributions")({
  component: Contributions,
  validateSearch: (search) => contributionsSearchSchema.parse(search),
})

function ContributionsTable() {
  const navigate = useNavigate({ from: Route.fullPath })
  const { page } = Route.useSearch()

  const { data, isLoading, isPlaceholderData } = useQuery({
    ...getContributionsQueryOptions({ page }),
    placeholderData: (prevData) => prevData,
  })

  const setPage = (page: number) =>
    navigate({
      search: (prev: { [key: string]: string }) => ({ ...prev, page }),
    })

  const contributions = data?.data.slice(0, PER_PAGE) ?? []
  const count = data?.count ?? 0

  if (isLoading) {
    return <PendingContributions />
  }

  if (contributions.length === 0) {
    return (
      <EmptyState.Root>
        <EmptyState.Content>
          <EmptyState.Indicator>
            <FiSearch />
          </EmptyState.Indicator>
          <VStack textAlign="center">
            <EmptyState.Title>You don't have any contributions yet</EmptyState.Title>
            <EmptyState.Description>
              Add a new contribution to get started
            </EmptyState.Description>
          </VStack>
        </EmptyState.Content>
      </EmptyState.Root>
    )
  }

  return (
    <>
      <Table.Root size={{ base: "sm", md: "md" }}>
        <Table.Header>
          <Table.Row>
            <Table.ColumnHeader w="30%">ID</Table.ColumnHeader>
            <Table.ColumnHeader w="30%">Idea Text</Table.ColumnHeader>
            <Table.ColumnHeader w="10%">Actions</Table.ColumnHeader>
          </Table.Row>
        </Table.Header>
        <Table.Body>
          {contributions?.map((contribution: ContributionPublic) => (
            <Table.Row key={contribution.id} opacity={isPlaceholderData ? 0.5 : 1}>
              <Table.Cell truncate maxW="30%">
                {contribution.id}
              </Table.Cell>
              <Table.Cell
                color={!contribution.idea_text ? "gray" : "inherit"}
                truncate
                maxW="30%"
              >
                {contribution.idea_text || "N/A"}
              </Table.Cell>
              <Table.Cell width="10%">
                <ContributionActionsMenu contribution={contribution} />
              </Table.Cell>
            </Table.Row>
          ))}
        </Table.Body>
      </Table.Root>
      <Flex justifyContent="flex-end" mt={4}>
        <PaginationRoot
          count={count}
          pageSize={PER_PAGE}
          onPageChange={({ page }) => setPage(page)}
        >
          <Flex>
            <PaginationPrevTrigger />
            <PaginationItems />
            <PaginationNextTrigger />
          </Flex>
        </PaginationRoot>
      </Flex>
    </>
  )
}

function Contributions() {
  return (
    <Container maxW="full">
      <Heading size="lg" pt={12}>
        Contributions Management
      </Heading>
      <AddContribution />
      <ContributionsTable />
    </Container>
  )
}
