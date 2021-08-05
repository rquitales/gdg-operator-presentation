import {Slide, Box, Heading} from 'spectacle';
import Editor from '../common/code';

const host = window.location.host;
const client = new WebSocket(`ws://${host}/socket`);
// const client = new WebSocket("ws://localhost:8082/socket");

const startCode = `#Make project directory and move there
mkdir -p ~/projects/serverless
cd ~/projects/serverless

# Initialize kubebuilder project
kubebuilder init --domain rquitales.com --repo github.com/rquitales/faas

# Add a new CRD
yes | kubebuilder create api --group serverless --version v1alpha1 --kind Function
`

export default function Page() {
  return (
    <Slide>
      <Heading>
        Scaffold Project
      </Heading>
      <Box>
        <Editor client={client} value={startCode} language={'sh'} runner="shell" outputPane={true}/>
      </Box>
    </Slide>
  );
}
