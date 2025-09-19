// import { useState } from "react";
import Controller from "./../../components/Speechtotext/Controller";
// import React from "react";
import * as React from "react"
import { createFileRoute } from "@tanstack/react-router"

export const Route = createFileRoute("/_layout/speakout")({
  component: Speakout,
})

function Speakout(): React.JSX.Element {
  //const [count, setCount] = useState(0);

  return (
    <div>
      <>
        <Controller />
      </>
    </div>
  );
}

export default Speakout;
