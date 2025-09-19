import { useState } from "react";
// import Title from "./Title";
import axios from "redaxios";
import RecordMessage from "./RecordMessage";
// import React from "react";
import * as React from "react"


// CHAKRA SURGERY
import {
  // Button,
  Flex,
  Heading,
  HStack,
  Input,
  Stack,
  Text,
} from '@chakra-ui/react';
import { Alert } from "@chakra-ui/react"
import { Spinner } from "@chakra-ui/react"

type MessageProps = {
  text: string;
  actor: 'user' | 'bot';
};
const Message = ({ text, actor }: MessageProps) => {
  return (
    <Flex
      p={4}
      bg={actor === 'user' ? 'blue.500' : 'gray.100'}
      color={actor === 'user' ? 'white' : 'gray.600'}
      borderRadius="lg"
      w="fit-content"
      alignSelf={actor === 'user' ? 'flex-end' : 'flex-start'}
    >
      <Text>{text}</Text>
    </Flex>
  );
};









const Controller = (): React.JSX.Element => {
  const [isLoading, setIsLoading] = useState(false);
  const [messages, setMessages] = useState<any[]>([]);
  /*
  function createBlobURL(data: any) {
    const blob = new Blob([data], { type: "audio/mpeg" });
    const url = window.URL.createObjectURL(blob);
    return url;
  }
  */

  const handleStop = async (blobUrl: string) => {
    setIsLoading(true);

    // Append recorded message to messages
    const myMessage = { sender: "me", blobUrl };
    const messagesArr = [...messages, myMessage];

    // convert blob url to blob object
    fetch(blobUrl)
      .then((res) => res.blob())
      .then(async (blob) => {
        // Construct audio to send file
        const formData = new FormData();
        formData.append("file", blob, "myFile.wav");

        // send form data to api endpoint
        await axios
          .post("http://localhost:8001/api/v1/transcribe/", formData, {
          // .post("http://mongo.pesto.io:8001/api/v1/transcribe/", formData, {
            headers: {
              // "Content-Type": "audio/mpeg",
              "Content-Type": "multipart/form-data",
            },
            // responseType: "arrayBuffer", // Set the response type to handle binary data
            responseType: "json",
          })
          .then((res: any) => {
            const transcribedMessage = res.data;
            // const audio = new Audio();
            // audio.src = createBlobURL(blob);

            // Append to audio
            const rachelMessage = { sender: "rachel", transcribedMessage: transcribedMessage };
            messagesArr.push(rachelMessage);
            setMessages(messagesArr);

            // Play audio
            setIsLoading(false);
            // audio.play();
          })
          .catch((err: any) => {
            console.error(err);
            setIsLoading(false);
          });
      });
  };

  return (
    <Flex h="100vh" py={12}>
      <Flex
        flexDirection="column"
        w="2xl"
        m="auto"
        h="full"
        borderWidth="1px"
        roundedTop="lg"
      >
        <HStack p={4} bg="blue.500">
          <Heading size="lg" color="white">
            Chat App
          </Heading>
        </HStack>

        <Stack
          px={4}
          py={8}
          overflow="auto"
          flex={1}
          css={{
            '&::-webkit-scrollbar': {
              width: '4px',
            },
            '&::-webkit-scrollbar-track': {
              width: '6px',
            },
            '&::-webkit-scrollbar-thumb': {
              background: '#d5e3f7',
              borderRadius: '24px',
            },
          }}
        >
          {messages?.map((message, index) => {
            return (
              <Message key={index + message.sender} text={message.transcribedMessage} actor={message.sender} />
            );
          })}

          {messages.length == 0 && !isLoading && (
            <Alert.Root status="info" title="This is the alert title">
              <Alert.Indicator />
              <Alert.Title>Send Rachel a message...</Alert.Title>
            </Alert.Root>
          )}

          {isLoading && (
          <Alert.Root status="info" title="This is the alert title">
            <Alert.Indicator />
            <Alert.Title> Gimme a few seconds...</Alert.Title>
            <Spinner size="sm" />
          </Alert.Root>
            

          )}
          <Message text="Hi" actor="user" />
          <Message text="How may I help you?" actor="bot" />
          <Message text="Hi" actor="user" />
          <Message text="How may I help you?" actor="bot" />
          <Message text="Hi" actor="user" />
          <Message text="How may I help you?" actor="bot" />
          <Message text="Hi" actor="user" />
          <Message text="How may I help you?" actor="bot" />
        </Stack>

        <HStack p={4} bg="gray.100">
          <Input bg="white" placeholder="Enter your text" />
          {/* Recorder */}
          <RecordMessage handleStop={handleStop} />
        </HStack>
      </Flex>
    </Flex>
  );
}
// END OF CHAKRA SURGERY

export default Controller;
