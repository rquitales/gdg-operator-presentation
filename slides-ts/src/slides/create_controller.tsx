import { Slide, Box, Heading } from 'spectacle';
import Editor from '../common/code';
import raw from 'raw.macro';
// import code from './code/functions_api.txt';

const host = window.location.host;
const client = new WebSocket(`ws://${host}/socket`);
// const client = new WebSocket('ws://localhost:8082/socket');

const startCode = raw('./code/functions_controller.txt');
const path = '/home/nonroot/projects/serverless/controllers/function_controller.go';

export default function Page() {
  return (
    <Slide>
      <Heading>Create Controller Logic</Heading>
      <Box>
      {path}
        <Editor
          client={client}
          value={startCode}
          language={'golang'}
          runner="saveFile"
          outputPane={false}
          path={path}
        />
      </Box>
    </Slide>
  );
}
