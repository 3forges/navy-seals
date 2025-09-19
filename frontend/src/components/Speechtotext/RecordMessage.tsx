import { ReactMediaRecorder } from "react-media-recorder";

// import React from "react";
import * as React from "react"
import { FiMic as Microphone } from "react-icons/fi"
import { Button, HStack } from "@chakra-ui/react"

type Props = {
  handleStop: any;
};

const RecordMessage = ({ handleStop }: Props): React.JSX.Element => {
  return (
    <ReactMediaRecorder
      audio
      onStop={handleStop}
      render={({ status, startRecording, stopRecording }) => (
        <div className="mt-2">


          <HStack>
            <Button
              colorPalette="teal"
              variant="solid"
              onMouseDown={startRecording}
              onMouseUp={stopRecording}
            >
              {
                /**
                 * https://github.com/DeltaCircuit/react-media-recorder?tab=readme-ov-file#status
                 */
              }
              <Microphone color={status==`recording`?`#EF0000`:`current`} /> {status}
            </Button>
          </HStack>
        </div>
      )}
    />
  );
};

export default RecordMessage;
