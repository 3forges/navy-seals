import {
  Button,
  ButtonGroup,
  DialogActionTrigger,
  Input,
  Text,
  VStack,
} from "@chakra-ui/react"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import { useState } from "react"
import { type SubmitHandler, useForm } from "react-hook-form"
import { FaExchangeAlt } from "react-icons/fa"

import { type ApiError, type ContributionPublic, ContributionsService } from "@/client"
import useCustomToast from "@/hooks/useCustomToast"
import { handleError } from "@/utils"
import {
  DialogBody,
  DialogCloseTrigger,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogRoot,
  DialogTitle,
  DialogTrigger,
} from "../ui/dialog"
import { Field } from "../ui/field"
import React from "react"

interface EditContributionProps {
  contribution: ContributionPublic
}

interface ContributionUpdateForm {
  // title: string
  // description?: string
  idea_text: string
}

const EditContribution = ({ contribution }: EditContributionProps): React.JSX.Element => {
  const [isOpen, setIsOpen] = useState(false)
  const queryClient = useQueryClient()
  const { showSuccessToast } = useCustomToast()
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<ContributionUpdateForm>({
    mode: "onBlur",
    criteriaMode: "all",
    defaultValues: {
      ...contribution,
      idea_text: contribution.idea_text ?? undefined,
    },
  })

  const mutation = useMutation({
    mutationFn: (data: ContributionUpdateForm) =>
      ContributionsService.updateContribution({ id: contribution.id, requestBody: data }),
    onSuccess: () => {
      showSuccessToast("Contribution updated successfully.")
      reset()
      setIsOpen(false)
    },
    onError: (err: ApiError) => {
      handleError(err)
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ["contributions"] })
    },
  })

  const onSubmit: SubmitHandler<ContributionUpdateForm> = async (data) => {
    mutation.mutate(data)
  }

  return (
    <DialogRoot
      size={{ base: "xs", md: "md" }}
      placement="center"
      open={isOpen}
      onOpenChange={({ open }) => setIsOpen(open)}
    >
      <DialogTrigger asChild>
        <Button variant="ghost">
          <FaExchangeAlt fontSize="16px" />
          Edit Contribution
        </Button>
      </DialogTrigger>
      <DialogContent>
        <form onSubmit={handleSubmit(onSubmit)}>
          <DialogHeader>
            <DialogTitle>Edit Contribution</DialogTitle>
          </DialogHeader>
          <DialogBody>
            <Text mb={4}>Update the contribution details below.</Text>
            <VStack gap={4}>
              <Field
                required
                invalid={!!errors.idea_text}
                errorText={errors.idea_text?.message}
                label="Your Idea"
              >
                <Input
                  id="idea_text"
                  {...register("idea_text", {
                    required: "Title is required",
                  })}
                  placeholder="Your Idea"
                  type="text"
                />
              </Field>
            </VStack>
          </DialogBody>

          <DialogFooter gap={2}>
            <ButtonGroup>
              <DialogActionTrigger asChild>
                <Button
                  variant="subtle"
                  colorPalette="gray"
                  disabled={isSubmitting}
                >
                  Cancel
                </Button>
              </DialogActionTrigger>
              <Button variant="solid" type="submit" loading={isSubmitting}>
                Save
              </Button>
            </ButtonGroup>
          </DialogFooter>
        </form>
        <DialogCloseTrigger />
      </DialogContent>
    </DialogRoot>
  )
}

export default EditContribution
