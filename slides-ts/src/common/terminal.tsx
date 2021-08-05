import Terminal from 'terminal-in-react';
import { Payload } from '../common/interfaces';

const host = window.location.host;
const client = new WebSocket(`ws://${host}/socket`);
// const client = new WebSocket('ws://localhost:8082/socket');

let resp = {
  done: false,
  responses: [] as Payload[],
};

client.onmessage = (message: any) => {
  const dataFromServer: Payload = JSON.parse(message.data);

  dataFromServer.Kind === 'end'
    ? (resp.done = true)
    : resp.responses.push(dataFromServer);
};

export default function Term() {
    return(
        <Terminal
        color="#333333"
        style={{height: "75%"}}
        prompt="#a917a8"
        outputColor="#333333"
        backgroundColor="#fdf6e4"
        barColor="#333333"
        commandPassThrough={(cmd, print: (msg: string | undefined) => void) => {
          resp.done = false;
          let data = JSON.stringify({
            Id: 'terminal',
            Kind: 'run',
            Body: '#!/bin/sh\n' + (cmd as unknown as string[]).join(' '),
          });
          console.log(data);
          client.send(data);

          const wait = async () => {
            if (resp.done) {
              resp.responses.forEach((ele) => print(ele.Body));
              resp.responses = [] as Payload[];
              return;
            }

            let last = resp.responses.shift();
            if (last) {
              print(last.Body);
            }
            await setTimeout(wait, 100);
            return;
          };

          wait();
        }}
        allowTabs={false}
      />
    )
}