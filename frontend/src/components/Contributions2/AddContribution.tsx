import axios from "axios";
import { useMutation, useQueryClient } from "@tanstack/react-query"
import { type SubmitHandler, useForm } from "react-hook-form"
import {
  Button,
  DialogActionTrigger,
  DialogTitle,
  Input,
  Text,
  VStack,
  Textarea,
  Alert,
  Spinner
} from "@chakra-ui/react"
import { useState } from "react"
import { FaPlus } from "react-icons/fa"

import { type ContributionCreate, ContributionsService } from "@/client"
import type { ApiError } from "@/client/core/ApiError"
import useCustomToast from "@/hooks/useCustomToast"
import { handleError } from "@/utils"
import {
  DialogBody,
  DialogCloseTrigger,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogRoot,
  DialogTrigger,
} from "../ui/dialog"

import { Field } from "../ui/field"
import React from "react"
import RecordMessage from "../Speechtotext/RecordMessage";

const AddContribution = (): React.JSX.Element => {
  const [isOpen, setIsOpen] = useState(false)
  const queryClient = useQueryClient()
  const { showSuccessToast } = useCustomToast()
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isValid, isSubmitting },
  } = useForm<ContributionCreate>({
    mode: "onBlur",
    criteriaMode: "all",
    defaultValues: {
      //title: "",
      //description: "",
      idea_text: '',
    },
  })

  const mutation = useMutation({
    mutationFn: (data: ContributionCreate) =>
      ContributionsService.createContribution({ requestBody: data }),
    onSuccess: () => {
      showSuccessToast("Contribution created successfully.")
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

  const onSubmit: SubmitHandler<ContributionCreate> = (data) => {
    mutation.mutate(data)
  }

  /**
   * State of the content of the texte area
   */
  const [ideaSentences, setIdeaSentences] = useState<string[]>(["Share Your Idea"])
  /**
   * A spinner icon will show that Whisper AI is working on transcribing the audio
   */
  const [isTranscribing, setIsTranscribing] = useState(false);
  /**
   * Microphone handler
   */
  const handleStopRecording = async (blobUrl: string) => {
    setIsTranscribing(true);

    // Append recorded message to ideaSentences
    const myIdeaSentence = `Just a test sentence added to your idea in the [handleStopRecording] method`;
    const ideaSentencesArr = [...ideaSentences, myIdeaSentence];

    // convert blob url to blob object
    fetch(blobUrl)
      .then((res) => res.blob())
      .then(async (blob) => {
        // Construct audio to send file
        const formData = new FormData();
        formData.append("audio", blob, "myVoiceFile.wav");

        // send form data to api endpoint
        await axios
          .post("http://localhost:8001/api/v1/transcribe/", formData, {
            // .post("http://mongo.pesto.io:8001/api/v1/transcribe/", formData, {
            headers: {
              // "Content-Type": "audio/mpeg",
              // "Content-Type": "multipart/form-data; boundary=------------------------7e1169bfdffcca45", // Simply don't set the Content-Type header manually and the browser will automatically set "multipart/form-data; boundary=..." value.
              "Content-Type": "multipart/form-data", // Simply don't set the Content-Type header manually and the browser will automatically set "multipart/form-data; boundary=..." value.
            },
            // responseType: "arrayBuffer", // Set the response type to handle binary data
            responseType: "json",
          })
          .then((res: any) => {
            const transcribedMessage = res.data;
            // const audio = new Audio();
            // audio.src = createBlobURL(blob);

            // Append to audio
            //ideaSentencesArr.push(transcribedMessage);
            setIdeaSentences([
              ...ideaSentences,
              ,transcribedMessage]);

            // Play audio
            setIsTranscribing(false);
            // audio.play();
          })
          .catch((err: any) => {
            console.error(err);
            setIsTranscribing(false);
          });
      });
  };

  return (
    <DialogRoot
      size={{ base: "xs", md: "md" }}
      placement="center"
      open={isOpen}
      onOpenChange={({ open }) => setIsOpen(open)}
    >
      <DialogTrigger asChild>
        <Button value="add-item" my={4}>
          <FaPlus fontSize="16px" />
          Add Contribution
        </Button>
      </DialogTrigger>
      <DialogContent>
        <form onSubmit={handleSubmit(onSubmit)}>
          <DialogHeader>
            <DialogTitle>Add Contribution</DialogTitle>
          </DialogHeader>
          <DialogBody>
            <Text mb={4}>Write your idea!</Text>
            <VStack gap={4}>
              {isTranscribing && (
                <Alert.Root status="info" title="Whisper AI">
                  <Alert.Indicator />
                  <Alert.Title> Whisper is transcribing, wait a few seconds...</Alert.Title>
                  <Spinner size="sm" />
                </Alert.Root>


              )}
              <Field
                invalid={!!errors.idea_text}
                errorText={errors.idea_text?.message}
                label="Your Idea"
              >
                <Textarea
                  id="idea_text"
                  value={ideaSentences.join('\r')}
                  placeholder={`Share your idea!`}
                  {...register("idea_text", { required: "Bio is required" })}
                />
              </Field>
              <RecordMessage handleStop={handleStopRecording} />
            </VStack>
          </DialogBody>

          <DialogFooter gap={2}>
            <DialogActionTrigger asChild>
              <Button
                variant="subtle"
                colorPalette="gray"
                disabled={isSubmitting}
              >
                Cancel
              </Button>
            </DialogActionTrigger>
            <Button
              variant="solid"
              type="submit"
              disabled={!isValid}
              loading={isSubmitting}
            >
              Save
            </Button>
          </DialogFooter>
        </form>
        <DialogCloseTrigger />
      </DialogContent>
    </DialogRoot>
  )
}

export default AddContribution
