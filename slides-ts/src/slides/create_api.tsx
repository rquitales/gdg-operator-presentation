import { Slide, Box, Heading } from 'spectacle';
import Editor from '../common/code';
import raw from 'raw.macro';
// import code from './code/functions_api.txt';

const host = window.location.host;
const client = new WebSocket(`ws://${host}/socket`);
// const client = new WebSocket('ws://localhost:8082/socket');

const startCode = raw('./code/functions_api.txt');
const path = '/home/nonroot/projects/serverless/api/v1alpha1/function_types.go';

export default function Page() {
  return (
    <Slide>
      <Heading>Define CRD</Heading>
      {path}
      <Box>
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
