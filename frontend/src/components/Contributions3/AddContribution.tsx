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
import { useEffect, useState } from "react"
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
  const [ideaSentences, setIdeaSentences] = useState<string[]>([])
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
        const controller = new AbortController();
        const run = async (controller:AbortController) => {
          // Construct audio to send file
          const formData = new FormData();
          // formData.append("audio", blob, "myVoiceFile.wav");
          formData.append("file", blob, "myVoiceFile.wav");
          // formData.append("language", "en");
          // formData.append("model", "small");
          formData.append("model", "medium");
          formData.append("language", "fr");
          formData.append("stream", 'true');
          const response = await fetch("http://localhost:8002/v1/audio/transcriptions", {
            method: "POST",
            body: formData,
          });
          if (response.body) {
            const reader = response.body.getReader();
            while (true) {
              const { done, value } = await reader.read();
              if (done) {
                // Do something with last chunk of data then exit reader
                const lastTranscribedMessage = new TextDecoder().decode(new Uint8Array(value?value:[]));
                console.log(`This is the last chiunk returned by the streaming response: [${lastTranscribedMessage}]`)
                console.log(`Speaches AI streaming response ended`)
                setIsTranscribing(false);
                return;
              } else {
                // Otherwise do something here to process current chunk
                console.log(`Speaches AI streaming response - This is NOT the last chiunk returned by the streaming response: [${value}]`)
                // const transcribedMessage = value;
                // const transcribedMessageDecoded = new TextDecoder().decode(new Uint8Array(value));
                const decodedValue = new TextDecoder().decode(new Uint8Array(value));
                const arrayOftranscribedMessageChunks = `${decodedValue}`.split('\n\n').map((item) => {
                  let toReturn = `{
                      "data": { "text": "" }
                    }`
                  if ( `${item}` == '' ) {
                    let toReturn = `{
                      "data": { "text": "" }
                    }`
                  } else {
                    toReturn = `{ ${item} }`
                  }
                  return toReturn.replace('data:', '"data":')
                }).map((item) => {
                  console.log(`Speaches AI streaming response - BEFORE JSON.PARSE - item is : [${JSON.stringify({item: item}, null, 2)}]`)
                  return JSON.parse(item)
                })
                console.log(`Speaches AI streaming response - arrayOftranscribedMessageChunks is : [${JSON.stringify({arrayOftranscribedMessageChunks: arrayOftranscribedMessageChunks}, null, 2)}]`)
                let transcribedMessage = ''
                arrayOftranscribedMessageChunks.forEach((item) => {
                  if (!item.data.text.includes(`Amara.org`)) {
                    transcribedMessage = `${transcribedMessage} \n ${item.data.text}`
                  }
                })
                
                console.log(`Speaches AI streaming response - Speaches AI returned: [${JSON.stringify({jbl_transcribedMessage: transcribedMessage}, null, 2)}]`)
                
                setIdeaSentences((ideaSentences) =>{ return [
                  ...ideaSentences,
                  `${transcribedMessage}`
                ]})
                
                // setIdeaSentences([
                //   ...ideaSentences,
                //   ,transcribedMessage]);
              }

              
            }

          } else {
            throw new Error(`response.body is undefined !`)
          }
        }
        run(controller)

        return () => controller.abort();
        
        })
        .catch((err: any) => {
            console.error(err);
            setIsTranscribing(false);
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
