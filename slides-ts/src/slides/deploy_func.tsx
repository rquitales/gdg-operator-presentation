import { Slide, Box, Heading } from 'spectacle';
import Editor from '../common/code';
import raw from 'raw.macro';
// import code from './code/functions_api.txt';

const host = window.location.host;
const client = new WebSocket(`ws://${host}/socket`);
// const client = new WebSocket('ws://localhost:8082/socket');

const startCode = raw('./code/deploy.txt');

export default function Page() {
  return (
    <Slide>
      <Heading>Deploy Function</Heading>
      <Box>
        <Editor
          client={client}
          value={startCode}
          language={'yaml'}
          runner="kubectlApply"
          outputPane={false}
        />
      </Box>
    </Slide>
  );
}
